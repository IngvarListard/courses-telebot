package store

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserRepository interface {
	// GetOrCreate creating user method
	GetOrCreate(*tgbotapi.User, int64) (*models.User, error)
}

type LearningNodeRepository interface {
	// GetCourses get all learning nodes that don't have a parent node
	GetNodesByParentID(parentID int) ([]*models.LearningNode, error)
	GetNodeByID(ID int) (*models.LearningNode, error)
}

type DocumentRepository interface {
	// GetDocumentByID
	GetDocumentByID(int) (*models.Document, error)
	GetDocumentsByParentID(parentID int) ([]*models.Document, error)
}

type ChatRepository interface {
	// GetOrCreate creating chat method
	GetOrCreate(*tgbotapi.Chat) (*models.Chat, error)
}
