package coursesbot

import (
	"context"
	"errors"
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	fnName int = iota
	args
	PageLimit = 10
)

type (
	callbackHandlers map[string]func(*Bot, *tgbotapi.CallbackQuery, string) error
	commandHandlers  map[string]func(*Bot, *tgbotapi.Message) error
)

type Bot struct {
	TgAPI  *tgbotapi.BotAPI
	Store  store.Store
	Config *Config
	callbackHandlers
	commandHandlers
}

func New(config *Config, store store.Store) (bot *Bot, err error) {
	var botApi *tgbotapi.BotAPI
	httpClient, err := NewHTTPProxyClient(config.Proxy)

	if config.Proxy != nil {
		botApi, err = tgbotapi.NewBotAPIWithClient(config.APIKey, httpClient)
	} else {
		botApi, err = tgbotapi.NewBotAPI(config.APIKey)
	}
	if err != nil {
		return
	}

	bot = &Bot{
		botApi,
		store,
		config,
		callbackHandlers{},
		commandHandlers{},
	}

	bot.registerCallbacks()
	bot.registerCommands()

	bot.TgAPI.Debug = config.Debug
	log.Printf("Authorized on account %s", bot.TgAPI.Self.UserName)
	return bot, err
}

func (b *Bot) Start(ctx context.Context) (err error) {

	if b.TgAPI == nil {
		log.Fatalln("improperly configured: no bot instance")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.TgAPI.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("error while getting updates: %v", err)
	}

	var wg sync.WaitGroup
	for {
		select {
		case update := <-updates:
			switch {
			case update.Message != nil:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := b.handleCommand(update.Message); err != nil {
						b.errorHandler(err, update.Message.Chat.ID, "message", update.Message.Text)
					}
				}()
			case update.CallbackQuery != nil:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := b.handleCallback(update.CallbackQuery); err != nil {
						b.errorHandler(err, update.CallbackQuery.Message.Chat.ID, "callback", update.CallbackQuery.Data)
					}
				}()
			}
		case <-ctx.Done():
			wg.Wait()
			log.Println("bot stopped")
			return
		}
	}

	return err
}

func (b *Bot) errorHandler(err error, chatID int64, uType, content string) {
	log.Printf("error during processing of the message: '%v'; update type: '%v'; content: '%v'", err, uType, content)

	text := "Во время обработки запроса произошла ошибка"
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := b.TgAPI.Send(msg); err != nil {
		log.Printf("error sending message: %v: %v", text, err)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) (err error) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	command := message.Command()
	if f, ok := b.commandHandlers[command]; ok {
		err = f(b, message)
	}

	return err
}

func (b *Bot) handleCallback(c *tgbotapi.CallbackQuery) (err error) {
	if _, e := b.TgAPI.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: c.ID}); e != nil {
		log.Printf("error during sending callback answer: %v", e)
	}

	data := strings.Split(c.Data, ":")
	callbackHandler, ok := b.callbackHandlers[data[fnName]]
	if !ok {
		return fmt.Errorf("unknown callback function: %v", data[fnName])
	}
	if err := callbackHandler(b, c, data[args]); err != nil {
		return fmt.Errorf("handleCallback: %v", err)
	}
	return
}

func (b *Bot) registerCommands() {
	b.commandHandlers["start"] = start
	b.commandHandlers["hello"] = hello
	b.commandHandlers["courses"] = courses
}

func (b *Bot) registerCallbacks() {
	b.callbackHandlers["sendNodeList"] = sendNodeList
	b.callbackHandlers["sendDocument"] = sendDocument
	b.callbackHandlers["upDir"] = upDirectory
	b.callbackHandlers["sendAllDocs"] = sendAllDocuments
	b.callbackHandlers["nextPage"] = nextPage
}

func NewHTTPProxyClient(config *ProxyConfig) (client *http.Client, err error) {

	if config == nil {
		err = errors.New("config is nil")
		return
	}

	if config.Type == "socks5" {
		dialer, err := proxy.SOCKS5(
			"tcp",
			os.Getenv("SOCKS5_URL"),
			&proxy.Auth{User: os.Getenv("SOCKS5_USER"), Password: os.Getenv("SOCKS5_PASSWORD")},
			proxy.Direct,
		)
		if err != nil {
			log.Panicf("Error in proxy %s", err)
		}

		client = &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}}}
	} else {
		msg := fmt.Sprintf("unsupported proxy type: %v\n", config.Type)
		log.Printf(msg)
		err = errors.New(msg)
		return
	}

	return
}
