package store

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserRepository interface {
	GetOrCreate(*tgbotapi.User, int64) (*models.User, error)
}

type LearningNodeRepository interface {
	GetNodeByID(ID int) (*models.LearningNode, error)
	GetNodesByParentID(parentID int) ([]*models.LearningNode, error)
}

type DocumentRepository interface {
	GetDocumentByID(int) (*models.Document, error)
	GetDocumentsByParentID(parentID int) ([]*models.Document, error)
}

type ChatRepository interface {
	GetOrCreate(*tgbotapi.Chat) (*models.Chat, error)
}
