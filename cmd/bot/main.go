package main

import (
	"github.com/IngvarListard/courses-telebot/internal/coursesbot"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// migrateSchema синхронизация сземы БД
func migrateSchema() {
	DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.LearningNode{}, &models.Document{})
}

func Start() {
	bot, err := tgbotapi.NewBotAPI(ApiKey)
	coursesbot.Setup(DB, bot)
	if err != nil {
		panic(err)
	}
	bot.Debug = Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if err := coursesbot.Handle(update); err != nil {
			text := update.Message.Text
			log.Printf("error during processing of the message: %v: %v", text, err)
			text = "Во время обработки зарпоса произошла ошибка"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("error sending message: %v: %v", text, err)
			}
		}
	}
}

func main() {
	defer func() {
		if err := DB.Close(); err != nil {
			log.Printf("error during closing the DB connection: %v", err)
		}
	}()
	migrateSchema()
	Start()
}
