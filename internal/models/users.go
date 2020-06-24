package models

type User struct {
	ID           int    `json:"id";gorm:"primary_key"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`
	ChatID       int64  `json:"chat_id"`
	Chat         *Chat  `json:"chat"`
	Model
}

type Chat struct {
	ID    int64  `json:"id";gorm:"primary_key"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Model
}
