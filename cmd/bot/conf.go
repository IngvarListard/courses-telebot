package main

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type Settings struct {
	ApiKey string
	Debug  bool
	Db     *gorm.DB
}

func (s *Settings) init() {
	db, err := gorm.Open("sqlite3", "courses_bot.db")
	if err != nil {
		log.Panic("failed to connect database")
	}

	s.Db = db
	s.ApiKey = os.Getenv("API_KEY")
	s.Debug = os.Getenv("DEBUG") == "1"
}

func (s *Settings) Validate() {
	if s.ApiKey == "" {
		log.Fatal("API_KEY is missing")
	}
}

var conf Settings

// Setup startup bot configuration
func init() {
	conf = Settings{}
	conf.init()
	conf.Validate()
}

func Start() {
	bot, err := tgbotapi.NewBotAPI(conf.ApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.Text == "/start" {
			user := conf.Db.FirstOrCreate(&models.User{TgId: update.Message.From.ID}, models.User{
				Name: update.Message.From.FirstName,
				TgId: update.Message.From.ID,
			})
			fmt.Println(user)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println("Ошибка отправки:", err)
		}
	}
}
