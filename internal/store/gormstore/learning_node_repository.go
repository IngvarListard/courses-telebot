package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
)

type LearningNodeRepository struct {
	store *Store
}

func (lnr *LearningNodeRepository) GetCourses() ([]*models.LearningNode, error) {
	r := new([]*models.LearningNode)
	return *r, nil
}
