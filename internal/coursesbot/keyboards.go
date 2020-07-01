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
	"send":  "\xF0\x9F\x93\xA9",
	"back":  "\xE2\x86\xA9",
	"next":  "\xe2\x9e\xa1",
	"prev":  "\xe2\xac\x85",
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

	bottomRow := tgbotapi.NewInlineKeyboardRow()

	if parentID != 0 {
		callback := fmt.Sprintf("upDir:%v", parentID)
		bottomRow = append(bottomRow, tgbotapi.NewInlineKeyboardButtonData(icons["up"]+" Вверх", callback))
	}

	if len(documents) > 0 {
		sendAllCb := fmt.Sprintf("sendAllDocs:%v", parentID)
		bottomRow = append(bottomRow, tgbotapi.NewInlineKeyboardButtonData(icons["send"]+" Отправить файлы", sendAllCb))
	}

	if len(bottomRow) > 0 {
		rows = append(rows, bottomRow)
	}

	nextPage := tgbotapi.NewInlineKeyboardButtonData(icons["next"], fmt.Sprintf("nextPage:%v,%v,%v", parentID, 2, PageLimit))
	prevPage := tgbotapi.NewInlineKeyboardButtonData(icons["prev"], "333")
	pagination := tgbotapi.NewInlineKeyboardRow(nextPage, prevPage)
	rows = append(rows, pagination)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}

func genCoursesKeyboard2(items []*PageItem) (*tgbotapi.InlineKeyboardMarkup, error) {

	var rows [][]tgbotapi.InlineKeyboardButton
	var parentID int
	hasDocuments := false

	for _, v := range items {
		button, err := v.GetInlineButton()
		if err != nil {
			return nil, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*button))
		parentID = v.ParentID
		if v.Type == "document" {
			hasDocuments = true
		}
	}

	bottomRow := tgbotapi.NewInlineKeyboardRow()

	if parentID != 0 {
		callback := fmt.Sprintf("upDir:%v", parentID)
		bottomRow = append(bottomRow, tgbotapi.NewInlineKeyboardButtonData(icons["up"]+" Вверх", callback))
	}

	if hasDocuments {
		sendAllCb := fmt.Sprintf("sendAllDocs:%v", parentID)
		bottomRow = append(bottomRow, tgbotapi.NewInlineKeyboardButtonData(icons["send"]+" Отправить файлы", sendAllCb))
	}

	if len(bottomRow) > 0 {
		rows = append(rows, bottomRow)
	}

	nextPage := tgbotapi.NewInlineKeyboardButtonData(icons["next"], fmt.Sprintf("nextPage:%v,%v,%v", parentID, 2, PageLimit))
	prevPage := tgbotapi.NewInlineKeyboardButtonData(icons["prev"], "333")
	pagination := tgbotapi.NewInlineKeyboardRow(nextPage, prevPage)
	rows = append(rows, pagination)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
