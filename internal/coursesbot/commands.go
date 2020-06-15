package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func start() commandHandler {
	return func(b *Bot, message *tgbotapi.Message) (err error) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "START")
		msg.ReplyToMessageID = message.MessageID
		if _, err := b.TgAPI.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
		return err
	}
}

func hello() commandHandler {
	return func(b *Bot, message *tgbotapi.Message) (err error) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "HELLO WORLD")
		if _, err := b.TgAPI.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
		return err
	}
}

func courses() commandHandler {
	return func(b *Bot, message *tgbotapi.Message) error {
		courses, err := b.Store.LearningNode().GetCourses()
		if err != nil {
			return err
		}

		keyboard, _ := genCoursesKeyboard(courses, []*models.Document{})

		msg := tgbotapi.NewMessage(message.Chat.ID, "Доступные курсы")
		msg.ReplyMarkup = keyboard
		_, err = b.TgAPI.Send(msg)
		return err
	}
}
