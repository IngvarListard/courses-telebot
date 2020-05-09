package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	DB       *gorm.DB
	Bot      *tgbotapi.BotAPI
	commands = map[string]func(tgbotapi.Update) error{
		"/start": startHandler,
		"/hello": helloHandler,
	}
)

func Setup(db *gorm.DB, bot *tgbotapi.BotAPI) {
	DB = db
	Bot = bot
}

func Handle(u tgbotapi.Update) (err error) {
	if DB == nil || Bot == nil {
		return fmt.Errorf("no DB instance in coursesbot handler")
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
	ID := u.Message.From.ID
	DB.FirstOrCreate(&models.User{ID: ID}, models.User{
		ID:           ID,
		FirstName:    u.Message.From.FirstName,
		LastName:     u.Message.From.LastName,
		UserName:     u.Message.From.UserName,
		LanguageCode: u.Message.From.LanguageCode,
		IsBot:        u.Message.From.IsBot,
	})

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
