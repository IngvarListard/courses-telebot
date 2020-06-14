package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
	"github.com/jinzhu/gorm"
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

func (s *Store) User() *UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}

func (s *Store) Chat() *ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = &ChatRepository{
			store: s,
		}
	}
	return s.chatRepository
}

func (s *Store) LearningNode() *LearningNodeRepository {
	if s.learningNodeRepository == nil {
		s.learningNodeRepository = &LearningNodeRepository{
			store: s,
		}
	}
	return s.learningNodeRepository
}

func (s *Store) Document() *DocumentRepository {
	if s.documentRepository == nil {
		s.documentRepository = &DocumentRepository{
			store: s,
		}
	}
	return s.documentRepository
}
