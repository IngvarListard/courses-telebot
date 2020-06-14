package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type ChatRepository struct {
	store *Store
}

func (ur *ChatRepository) GetOrCreate(c *tgbotapi.Chat) (*models.Chat, error) {
	newChat := models.Chat{Type: c.Type, Title: c.Title}

	if err := ur.store.db.Where(models.Chat{ID: c.ID}).Attrs(newChat).FirstOrCreate(&newChat).Error; err != nil {
		log.Printf("error creating chat in database: %v", err)
		return nil, err
	}
	return &newChat, nil
}
