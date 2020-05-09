package models

import "time"

type Model struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index";json:"deleted_at"`
}

type LearningNode struct {
	Model
	ID       int    `json:"id";gorm:"primary_key"`
	Name     string `json:"name"`
	ParentID uint   `json:"parent_id"`
	Parent   *LearningNode
}

type Document struct {
	Model
	ID       int           `json:"id";gorm:"primary_key"`
	Name     string        `json:"name"`
	FileID   string        `json:"file_id"`
	Type     string        `json:"type"`
	Priority float32       `json:"priority"`
	NodeID   string        `json:"node_id"`
	Node     *LearningNode `json:"node"`
	URL      string        `json:"url"`
}
