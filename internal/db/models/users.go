package models

type User struct {
	Model
	ID           int    `json:"id";gorm:"primary_key"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`
}

type Chat struct {
	Model
	ID    int    `json:"id";gorm:"primary_key"`
	Type  string `json:"type"`
	Title string `json:"title"`
}
