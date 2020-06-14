package gormstore

import "github.com/IngvarListard/courses-telebot/internal/models"

type DocumentRepository struct {
	store *Store
}

func (dr *DocumentRepository) GetDocumentByID(ID int) (*models.Document, error) {
	d := new(*models.Document)
	return *d, nil
}
