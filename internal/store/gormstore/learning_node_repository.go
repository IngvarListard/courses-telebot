package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
)

type LearningNodeRepository struct {
	store *Store
}

// GetCourses returns list of all available courses
func (lnr *LearningNodeRepository) GetCourses() ([]*models.LearningNode, error) {
	var nodes []*models.LearningNode
	err := lnr.store.db.Where("parent_id IS NULL").Find(&nodes).Error
	return nodes, err
}

func (lnr *LearningNodeRepository) GetNodesByParentID(parentID int) ([]*models.LearningNode, error) {
	var nodes []*models.LearningNode
	err := lnr.store.db.Where(models.LearningNode{ParentID: parentID}).Find(&nodes).Error
	return nodes, err
}
