package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
)

type LearningNodeRepository struct {
	store *Store
}

func (lnr *LearningNodeRepository) GetNodesByParentID(parentID int) ([]*models.LearningNode, error) {
	var nodes []*models.LearningNode
	var expression interface{}
	if parentID == 0 {
		expression = "parent_id IS NULL"
	} else {
		expression = models.LearningNode{ParentID: parentID}
	}
	err := lnr.store.db.Where(expression).Order("priority").Find(&nodes).Error
	return nodes, err
}

func (lnr *LearningNodeRepository) GetNodeByID(nodeID int) (*models.LearningNode, error) {
	node := new(models.LearningNode)
	err := lnr.store.db.Where(models.LearningNode{ID: nodeID}).Find(node).Error
	return node, err
}
