package db

import (
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	return DB, err
}

// MigrateSchema database schema synchronization
func MigrateSchema() {
	DB.AutoMigrate(&models.Chat{}, &models.User{}, &models.LearningNode{}, &models.Document{})
}

func AddConversation(user *tgbotapi.User, chat *tgbotapi.Chat) {
	newChat := models.Chat{Type: chat.Type, Title: chat.Title}
	DB.Where(models.Chat{ID: chat.ID}).Attrs(newChat).FirstOrCreate(&newChat)

	newUser := models.User{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		LanguageCode: user.LanguageCode,
		IsBot:        user.IsBot,
	}
	if newChat.ID != 0 {
		newUser.ChatID = newChat.ID
	}
	DB.Where(models.User{ID: user.ID}).Attrs(newUser).FirstOrCreate(&newUser)
}

func GetCourses() []models.LearningNode {
	var nodes []models.LearningNode
	DB.Where("parent_id IS NULL").Find(&nodes)
	return nodes
}
