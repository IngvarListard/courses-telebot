package coursesbot

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

const (
	argsSep = ","
	nodeID  = iota
	page
	pType
)

func sendNodeList(b *Bot, c *tgbotapi.CallbackQuery, nodeID string) (err error) {
	nodeIDint, err := strconv.Atoi(nodeID)
	if err != nil {
		return fmt.Errorf("incorrect node ID in callback: %v", err)
	}

	nodes, err := b.Store.LearningNode().GetNodesByParentID(nodeIDint)
	if err != nil {
		err = fmt.Errorf("error querying learning nodes: %v", err)
	}

	documents := *new([]*models.Document)
	if nodeIDint != 0 {
		documents, err = b.Store.Document().GetDocumentsByParentID(nodeIDint)
		if err != nil {
			err = fmt.Errorf("error querying documents: %v", err)
		}
	}

	keyboard, _ := genCoursesKeyboard(nodes, documents)

	msg := tgbotapi.NewEditMessageReplyMarkup(c.Message.Chat.ID, c.Message.MessageID, *keyboard)

	_, err = b.TgAPI.Send(msg)
	return
}

func sendDocument(b *Bot, c *tgbotapi.CallbackQuery, documentID string) error {
	documentIDint, err := strconv.ParseInt(documentID, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect node ID in callback: %v", err)
	}

	document, err := b.Store.Document().GetDocumentByID(int(documentIDint))
	if err != nil {
		return fmt.Errorf("error getting document by id: %v", err)
	}

	d := tgbotapi.NewDocumentShare(c.Message.Chat.ID, document.FileID)
	_, err = b.TgAPI.Send(d)

	return err
}

func sendPage(b *Bot, c *tgbotapi.CallbackQuery, argsStr string) error {
	args := strings.Split(argsStr, argsSep)
	if args[pType] == "node" {

	} else if args[pType] == "document" {

	}
	return nil
}

func upDirectory(b *Bot, c *tgbotapi.CallbackQuery, nodeID string) error {
	parentID, err := strconv.Atoi(nodeID)
	parent, err := b.Store.LearningNode().GetNodeByID(parentID)
	if err != nil {
		return err
	}

	ID := strconv.Itoa(parent.ParentID)
	return sendNodeList(b, c, ID)
}

func sendAllDocuments(b *Bot, c *tgbotapi.CallbackQuery, nodeID string) error {
	parentID, err := strconv.Atoi(nodeID)
	if err != nil {
		return fmt.Errorf("sendAllDocuments: error parsing parentID: %v", nodeID)
	}
	documents, err := b.Store.Document().GetDocumentsByParentID(parentID)
	if err != nil {
		return fmt.Errorf("sendAllDocuments: error getting documents: %v", err)
	}

	for _, doc := range documents {
		d := tgbotapi.NewDocumentShare(c.Message.Chat.ID, doc.FileID)
		_, err = b.TgAPI.Send(d)
	}
	return err
}
