package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

const argsSep = ","

func sendNodeList(c *tgbotapi.CallbackQuery, nodeID string) (err error) {
	nodeIDint, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect node ID in callback: %v", err)
	}

	var nodes []*models.LearningNode
	var documents []*models.Document
	if e := db.DB.Where(models.LearningNode{ParentID: int(nodeIDint)}).Find(&nodes).Error; e != nil {
		err = fmt.Errorf("error querying learning nodes: %v", err)
	}
	if e := db.DB.Where(models.Document{NodeID: int(nodeIDint)}).Find(&documents).Error; e != nil {
		err = fmt.Errorf("error querying learning nodes: %v", err)
	}
	keyboard, _ := genCoursesKeyboard(nodes, documents)

	msg := tgbotapi.NewMessage(c.Message.Chat.ID, "Доступные курсы")
	msg.ReplyMarkup = keyboard
	_, err = Bot.Send(msg)
	return
}

func sendDocument(c *tgbotapi.CallbackQuery, documentID string) error {
	document := &models.Document{}
	documentIDint, err := strconv.ParseInt(documentID, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect node ID in callback: %v", err)
	}
	db.DB.First(document, models.Document{ID: int(documentIDint)})
	d := tgbotapi.NewDocumentShare(c.Message.Chat.ID, document.FileID)
	_, err = Bot.Send(d)

	return err
}

func sendPage(c *tgbotapi.CallbackQuery, argsStr string) error {
	const (
		nodeID = iota
		page
		pType
	)
	args := strings.Split(argsStr, argsSep)
	if args[pType] == "node" {

	} else if args[pType] == "document" {

	}
	return nil
}
