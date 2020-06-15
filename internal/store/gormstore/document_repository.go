package gormstore

import (
	"github.com/IngvarListard/courses-telebot/internal/models"
)

type DocumentRepository struct {
	store *Store
}

func (dr *DocumentRepository) GetDocumentsByParentID(parentID int) ([]*models.Document, error) {
	var documents []*models.Document
	err := dr.store.db.Where(models.Document{NodeID: parentID}).Find(&documents).Error
	return documents, err
}

func (dr *DocumentRepository) GetDocumentByID(ID int) (*models.Document, error) {
	document := &models.Document{}
	err := dr.store.db.First(document, models.Document{ID: ID}).Error
	return document, err
}
