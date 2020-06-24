package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

var icons = map[string]string{
	"dir":   "\xF0\x9F\x93\x81 ",
	"audio": "\xF0\x9F\x94\x8A ",
	"doc":   "\xF0\x9F\x93\x84",
	"up":    "\xE2\xAC\x86",
}

func genCoursesKeyboard(nodes []*models.LearningNode, documents []*models.Document) (*tgbotapi.InlineKeyboardMarkup, error) {

	var rows [][]tgbotapi.InlineKeyboardButton
	var parentID int

	for _, v := range nodes {
		callback := fmt.Sprintf("sendNodeList:%v", v.ID)
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(icons["dir"]+v.Name, callback))
		parentID = v.ParentID
		rows = append(rows, newRow)
	}
	for _, v := range documents {
		icon := icons["doc"]
		if strings.HasSuffix(v.FileName, ".mp3") {
			icon = icons["audio"]
		}
		callback := fmt.Sprintf("sendDocument:%v", v.ID)
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(icon+v.Name, callback))
		parentID = v.NodeID
		rows = append(rows, newRow)
	}
	if parentID != 0 {
		callback := fmt.Sprintf("upDir:%v", parentID)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(icons["up"], callback)))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
