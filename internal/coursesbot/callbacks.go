package coursesbot

import (
	"fmt"
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

func sendNodeList() callbackHandler {
	return func(b *Bot, c *tgbotapi.CallbackQuery, nodeID string) (err error) {
		nodeIDint, err := strconv.ParseInt(nodeID, 10, 64)
		if err != nil {
			return fmt.Errorf("incorrect node ID in callback: %v", err)
		}

		nodes, err := b.Store.LearningNode().GetNodesByParentID(int(nodeIDint))
		if err != nil {
			err = fmt.Errorf("error querying learning nodes: %v", err)
		}

		documents, err := b.Store.Document().GetDocumentsByParentID(int(nodeIDint))
		if err != nil {
			err = fmt.Errorf("error querying documents: %v", err)
		}

		keyboard, _ := genCoursesKeyboard(nodes, documents)

		msg := tgbotapi.NewEditMessageReplyMarkup(c.Message.Chat.ID, c.Message.MessageID, *keyboard)

		_, err = b.TgAPI.Send(msg)
		return
	}
}

func sendDocument() callbackHandler {
	return func(b *Bot, c *tgbotapi.CallbackQuery, documentID string) error {
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
}

func sendPage() callbackHandler {
	return func(b *Bot, c *tgbotapi.CallbackQuery, argsStr string) error {
		args := strings.Split(argsStr, argsSep)
		if args[pType] == "node" {

		} else if args[pType] == "document" {

		}
		return nil
	}
}
