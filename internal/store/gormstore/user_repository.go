package gormstore

import (
	"errors"
	"github.com/IngvarListard/courses-telebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) GetOrCreate(u *tgbotapi.User, chatID int64) (*models.User, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}

	newUser := models.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		UserName:     u.UserName,
		LanguageCode: u.LanguageCode,
		IsBot:        u.IsBot,
	}
	if chatID != 0 {
		newUser.ChatID = chatID
	}

	if err := ur.store.db.Where(models.User{ID: u.ID}).Attrs(newUser).FirstOrCreate(&newUser).Error; err != nil {
		log.Printf("error creating user in database: %v", err)
		return nil, err
	}
	return &newUser, nil
}
