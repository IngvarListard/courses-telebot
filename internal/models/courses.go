package models

import "time"

type Model struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index";json:"deleted_at"`
}

type LearningNode struct {
	ID       int     `json:"id";gorm:"primary_key"`
	Name     string  `json:"name"`
	DirName  string  `json:"dir_name"`
	Path     string  `json:"path"`
	Priority float32 `json:"priority"`
	ParentID int     `json:"parent_id"`
	Parent   *LearningNode
	Model
}

type Document struct {
	ID       int           `json:"id";gorm:"primary_key"`
	Name     string        `json:"name"`
	FileName string        `json:"file_name"`
	FileID   string        `json:"file_id"`
	Path     string        `json:"path"`
	Type     string        `json:"type"`
	Priority float32       `json:"priority"`
	NodeID   int           `json:"node_id"`
	Node     *LearningNode `json:"node"`
	URL      string        `json:"url"`
	Model
}
