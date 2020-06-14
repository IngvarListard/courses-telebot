package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/models"
	"github.com/IngvarListard/courses-telebot/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func start(message *tgbotapi.Message) (err error) {

	msg := tgbotapi.NewMessage(message.Chat.ID, "START")
	msg.ReplyToMessageID = message.MessageID
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return err
}

func hello(message *tgbotapi.Message) (err error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "HELLO WORLD")
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return err
}

func courses(message *tgbotapi.Message) (err error) {
	courses := store.GetCourses()
	keyboard, _ := genCoursesKeyboard(courses, []*models.Document{})

	msg := tgbotapi.NewMessage(message.Chat.ID, "Доступные курсы")
	msg.ReplyMarkup = keyboard
	_, err = Bot.Send(msg)
	return err
}
