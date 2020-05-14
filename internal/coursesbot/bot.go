package coursesbot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"os"
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
