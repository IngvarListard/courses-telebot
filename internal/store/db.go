package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// Database connection setup
func NewDB() (*gorm.DB, error) {
	var err error
	DB, err := gorm.Open("sqlite3", "courses_bot.db")
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return nil, err
	}
	return DB, err
}
