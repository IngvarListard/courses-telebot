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

type PageArgs struct{ ParentID, PerPage, LastNodeID, LastDocumentID int }

type PagedCourse struct {
	perPage   int
	nodes     []*models.LearningNode
	documents []*models.Document
	items     []*PageItem
}

func NewPagedCourse(perPage int, nodes []*models.LearningNode, documents []*models.Document) *PagedCourse {

	course := &PagedCourse{
		perPage:   perPage,
		nodes:     nodes,
		documents: documents,
	}

	for _, v := range nodes {
		course.items = append(course.items, &PageItem{
			ID:       v.ID,
			ParentID: v.ParentID,
			Name:     v.Name,
			Type:     "node",
			GetInlineButton: func() (*tgbotapi.InlineKeyboardButton, error) {
				callback := fmt.Sprintf("sendNodeList:%v", v.ID)
				newButton := tgbotapi.NewInlineKeyboardButtonData(icons["dir"]+v.Name, callback)
				return &newButton, nil
			},
		})
	}

	for _, v := range documents {
		course.items = append(course.items, &PageItem{
			ID:       v.ID,
			ParentID: v.NodeID,
			Name:     v.Name,
			Type:     "document",
			GetInlineButton: func() (*tgbotapi.InlineKeyboardButton, error) {
				icon := icons["doc"]
				if strings.HasSuffix(v.Name, ".mp3") {
					icon = icons["audio"]
				}
				callback := fmt.Sprintf("sendDocument:%v", v.ID)
				newButton := tgbotapi.NewInlineKeyboardButtonData(icon+v.Name, callback)

				return &newButton, nil
			},
		})
	}

	return course
}

func (pc *PagedCourse) GetPage(n int) *Page {
	last := n * pc.perPage
	first := last - pc.perPage
	items := pc.items[first:last]
	page := &Page{
		pageNumber:  1,
		PagedCourse: pc,
		Items:       items,
	}
	return page
}

type PageItem struct {
	ID              int
	ParentID        int
	Name            string
	Type            string
	GetInlineButton func() (*tgbotapi.InlineKeyboardButton, error)
}

type Page struct {
	pageNumber int
	Items      []*PageItem
	*PagedCourse
}

func (pc *Page) HasNext() bool { return false }

func nextPage(b *Bot, c *tgbotapi.CallbackQuery, argsStr string) (err error) {
	var nodes []*models.LearningNode
	var documents []*models.Document

	var iArgs [4]int

	sArgs := strings.Split(argsStr, argsSep)

	for i, sArg := range sArgs {
		iArg, err := strconv.Atoi(sArg)
		if err != nil {
			return fmt.Errorf("nextPage: error parsing arg %v: %v", sArg, err)
		}
		iArgs[i] = iArg
	}

	nodes, err = b.Store.LearningNode().GetNodesByParentID(iArgs[0])

	if len(nodes) != PageLimit {
		docsLimit := PageLimit - len(nodes)

		documents, err = b.Store.Document().GetDocumentsByParentID(iArgs[0])
		if len(documents) > docsLimit {
			documents = documents[:docsLimit]
		}
	}

	course := NewPagedCourse(iArgs[2], nodes, documents)

	page := course.GetPage(iArgs[0])

	keyboard, _ := genCoursesKeyboard2(page.Items)
	msg := tgbotapi.NewEditMessageReplyMarkup(c.Message.Chat.ID, c.Message.MessageID, *keyboard)
	_, err = b.TgAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("nextPage: error sending response: %v", err)
	}

	return nil

}
