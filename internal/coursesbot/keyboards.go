package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func genCoursesKeyboard(nodes []*models.LearningNode, documents []*models.Document) (*tgbotapi.InlineKeyboardMarkup, error) {

	var rows [][]tgbotapi.InlineKeyboardButton

	// TODO: pagination, navigation arrows, реализовать функцию обхода через интерфейс вместо двух циклов
	for _, v := range nodes {
		callback := fmt.Sprintf("getNodeList:%v", v.ID)
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x93\x81 "+v.Name, callback))
		rows = append(rows, newRow)
	}
	for _, v := range documents {
		callback := fmt.Sprintf("getDocument:%v", v.ID)
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\x8A"+v.Name, callback))
		rows = append(rows, newRow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
