package db

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	DB *gorm.DB
)

// Database connection setup
func Setup() (DB *gorm.DB, err error) {
	DB, err = gorm.Open("sqlite3", "courses_bot.db")
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return nil, err
	}
	MigrateSchema(DB)
	//DB.AutoMigrate(&models.Chat{}, &models.User{}, &models.LearningNode{}, &models.Document{})
	return DB, err
}

// MigrateSchema database schema synchronization
func MigrateSchema(DB *gorm.DB) {
	DB.AutoMigrate(&models.Chat{}, &models.User{}, &models.LearningNode{}, &models.Document{})
}

func AddConversation(user *tgbotapi.User, chat *tgbotapi.Chat) {
	c := DB.FirstOrCreate(&models.Chat{ID: chat.ID}, models.Chat{
		ID:    chat.ID,
		Type:  chat.Type,
		Title: chat.Title,
	})
	fmt.Printf("AAAAA: %v", c.Value)

	DB.FirstOrCreate(&models.User{ID: user.ID}, models.User{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		LanguageCode: user.LanguageCode,
		IsBot:        user.IsBot,
	})
}
