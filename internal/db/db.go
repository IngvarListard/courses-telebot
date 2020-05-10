package db

import (
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	DB *gorm.DB
)

// Database connection setup
func Setup() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open("sqlite3", "courses_bot.db")
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return nil, err
	}
	MigrateSchema()
	return DB, err
}

// MigrateSchema database schema synchronization
func MigrateSchema() {
	DB.AutoMigrate(&models.Chat{}, &models.User{}, &models.LearningNode{}, &models.Document{})
}

func AddConversation(user *tgbotapi.User, chat *tgbotapi.Chat) {
	newChat := DB.FirstOrCreate(&models.Chat{ID: chat.ID}, models.Chat{
		ID:    chat.ID,
		Type:  chat.Type,
		Title: chat.Title,
	})

	chatVal := newChat.Value.(*models.Chat)
	newUser := models.User{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		LanguageCode: user.LanguageCode,
		IsBot:        user.IsBot,
	}
	if chatVal != nil {
		newUser.ChatID = chatVal.ID
	}
	DB.FirstOrCreate(&models.User{ID: user.ID}, newUser)
}
