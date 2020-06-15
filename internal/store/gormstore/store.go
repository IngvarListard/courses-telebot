package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
	"github.com/IngvarListard/courses-telebot/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"log"
)

type Store struct {
	db                     *gorm.DB
	userRepository         *UserRepository
	chatRepository         *ChatRepository
	learningNodeRepository *LearningNodeRepository
	documentRepository     *DocumentRepository
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

// MigrateSchema database schema synchronization
func (s *Store) MigrateSchema() {
	s.db.AutoMigrate(&models.Chat{}, &models.User{}, &models.LearningNode{}, &models.Document{})
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}

func (s *Store) Chat() store.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = &ChatRepository{
			store: s,
		}
	}
	return s.chatRepository
}

func (s *Store) LearningNode() store.LearningNodeRepository {
	if s.learningNodeRepository == nil {
		s.learningNodeRepository = &LearningNodeRepository{
			store: s,
		}
	}
	return s.learningNodeRepository
}

func (s *Store) Document() store.DocumentRepository {
	if s.documentRepository == nil {
		s.documentRepository = &DocumentRepository{
			store: s,
		}
	}
	return s.documentRepository
}

// AddConversation add user and related chat to db
func (s *Store) AddConversation(user *tgbotapi.User, chat *tgbotapi.Chat) {
	newChat, err := s.Chat().GetOrCreate(chat)
	if err != nil {
		log.Printf("error creating chat in database: %v", err)
	}
	_, err = s.User().GetOrCreate(user, newChat.ID)
	if err != nil {
		log.Printf("error creating user in database: %v", err)
	}
}
