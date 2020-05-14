package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

func coursesCommand(message *tgbotapi.Message) (err error) {
	courses := db.GetCourses()
	keyboard, _ := genCoursesKeyboard(courses, []*models.Document{})

	msg := tgbotapi.NewMessage(message.Chat.ID, "Доступные курсы")
	msg.ReplyMarkup = keyboard
	_, err = Bot.Send(msg)
	return err
}
