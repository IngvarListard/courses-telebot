package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var (
	ApiKey string
	Debug  bool
	DB     *gorm.DB
)

// Setup bot configuration
func init() {
	db, err := gorm.Open("sqlite3", "courses_bot.db")
	if err != nil {
		log.Panic("failed to connect database")
	}

	DB = db
	ApiKey = os.Getenv("API_KEY")
	Debug = os.Getenv("DEBUG") == "1"
	if ApiKey == "" {
		log.Fatal("API_KEY is missing")
	}
}
