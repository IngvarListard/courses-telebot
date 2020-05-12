package coursesbot

import (
	"context"
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	Bot      *tgbotapi.BotAPI
	commands = map[string]func(tgbotapi.Update) error{
		"/start": startHandler,
		"/hello": helloHandler,
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

func Start() (err error) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("error while getting updates: %v", err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if err := Handle(update); err != nil {
			text := update.Message.Text
			log.Printf("error during processing of the message: %v: %v", text, err)
			text = "Во время обработки зарпоса произошла ошибка"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			if _, err := Bot.Send(msg); err != nil {
				log.Printf("error sending message: %v: %v", text, err)
			}
		}
	}
	return err
}

func Handle(u tgbotapi.Update) (err error) {
	if Bot == nil {
		return fmt.Errorf("improperly configured")
	}
	log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)
	text := u.Message.Text
	// TODO: bot commands are specified as separate type MessageEntity
	if f, ok := commands[text]; ok {
		err = f(u)
	}

	return err
}

func startHandler(u tgbotapi.Update) (err error) {
	db.AddConversation(u.Message.From, u.Message.Chat)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "START")
	msg.ReplyToMessageID = u.Message.MessageID
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return err
}

func helloHandler(u tgbotapi.Update) (err error) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "HELLO WORLD")
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return err
}
