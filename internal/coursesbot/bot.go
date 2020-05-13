package coursesbot

import (
	"context"
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	fnName int = iota
	args
)

var (
	Bot      *tgbotapi.BotAPI
	commands = map[string]func(*tgbotapi.Message) error{
		"start":   startCommand,
		"hello":   helloCommand,
		"courses": coursesCommand,
	}
	cbHandlers = map[string]func(*tgbotapi.CallbackQuery, string) error{
		"getNodeList": getNodeList,
		"getDocument": getDocument,
	}
)

func Setup(APIKey string, Debug bool) (err error) {
	dialer, proxyErr := proxy.SOCKS5(
		"tcp",
		os.Getenv("SOCKS5_URL"),
		&proxy.Auth{User: os.Getenv("SOCKS5_USER"), Password: os.Getenv("SOCKS5_PASSWORD")},
		proxy.Direct,
	)
	if proxyErr != nil {
		log.Panicf("Error in proxy %s", proxyErr)
	}
	client := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}}}
	if Bot, err = tgbotapi.NewBotAPIWithClient(APIKey, client); err != nil {
		return err
	}

	Bot.Debug = Debug
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	return err
}

func errorHandler(err error, chatID int64, uType, content string) {
	log.Printf("error during processing of the message: '%v'; update type: '%v'; content: '%v'", err, uType, content)

	text := "Во время обработки запроса произошла ошибка"
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("error sending message: %v: %v", text, err)
	}
}

func Start() (err error) {
	if Bot == nil {
		log.Fatalln("improperly configured: no bot instance")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("error while getting updates: %v", err)
	}

	for update := range updates {
		switch {
		case update.Message != nil:
			if err := HandleMessage(update.Message); err != nil {
				errorHandler(err, update.Message.Chat.ID, "message", update.Message.Text)
			}
		case update.CallbackQuery != nil:
			if err := HandleCallback(update.CallbackQuery); err != nil {
				errorHandler(err, update.CallbackQuery.Message.Chat.ID, "callback", update.CallbackQuery.Data)
			}
		}
	}
	return err
}

func HandleMessage(message *tgbotapi.Message) (err error) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	//text := u.Message.Text
	command := message.Command()
	if f, ok := commands[command]; ok {
		err = f(message)
	}

	return err
}

func HandleCallback(c *tgbotapi.CallbackQuery) (err error) {
	fmt.Printf("%v\n", c.Data)
	if _, e := Bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: c.ID}); e != nil {
		log.Printf("error during sending callback answer: %v", e)
	}

	data := strings.Split(c.Data, ":")
	callbackHandler, ok := cbHandlers[data[fnName]]
	if !ok {
		return fmt.Errorf("unknown callback function: %v", data[fnName])
	}
	if err := callbackHandler(c, data[args]); err != nil {
		return fmt.Errorf("")
	}
	return
}

func getNodeList(c *tgbotapi.CallbackQuery, nodeID string) (err error) {
	nodeIDint, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect node ID in callback: %v", err)
	}

	var nodes []*models.LearningNode
	var documents []*models.Document
	if e := db.DB.Where(models.LearningNode{ParentID: int(nodeIDint)}).Find(&nodes).Error; e != nil {
		err = fmt.Errorf("error querying learning nodes: %v", err)
	}
	if e := db.DB.Where(models.Document{NodeID: int(nodeIDint)}).Find(&documents).Error; e != nil {
		err = fmt.Errorf("error querying learning nodes: %v", err)
	}
	keyboard, _ := genCoursesKeyboard(nodes, documents)

	msg := tgbotapi.NewMessage(c.Message.Chat.ID, "Доступные курсы")
	msg.ReplyMarkup = keyboard
	_, err = Bot.Send(msg)
	return
}

func getDocument(c *tgbotapi.CallbackQuery, documentID string) error {
	document := &models.Document{}
	documentIDint, err := strconv.ParseInt(documentID, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect node ID in callback: %v", err)
	}
	db.DB.First(document, models.Document{ID: int(documentIDint)})
	d := tgbotapi.NewDocumentShare(c.Message.Chat.ID, document.FileID)
	_, err = Bot.Send(d)

	return err
}

func startCommand(message *tgbotapi.Message) (err error) {
	db.AddConversation(message.From, message.Chat)

	msg := tgbotapi.NewMessage(message.Chat.ID, "START")
	msg.ReplyToMessageID = message.MessageID
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return err
}

func helloCommand(message *tgbotapi.Message) (err error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "HELLO WORLD")
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return err
}

func genCoursesKeyboard(nodes []*models.LearningNode, documents []*models.Document) (*tgbotapi.InlineKeyboardMarkup, error) {

	var rows [][]tgbotapi.InlineKeyboardButton

	// TODO: pagination, navigation arrows, реализовать функцию обхода через интерфейс вместо двух циклов
	for _, v := range nodes {
		callback := fmt.Sprintf("getNodeList:%v", v.ID)
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(v.Name, callback))
		rows = append(rows, newRow)
	}
	for _, v := range documents {
		callback := fmt.Sprintf("getDocument:%v", v.ID)
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(v.Name, callback))
		rows = append(rows, newRow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}

func coursesCommand(message *tgbotapi.Message) (err error) {
	courses := db.GetCourses()
	keyboard, _ := genCoursesKeyboard(courses, []*models.Document{})

	msg := tgbotapi.NewMessage(message.Chat.ID, "Доступные курсы")
	msg.ReplyMarkup = keyboard
	_, err = Bot.Send(msg)
	return err
}
