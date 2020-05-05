package main

import (
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// migrateSchema синхронизация сземы БД
func migrateSchema() {
	Db.AutoMigrate(&models.User{})
}

func main() {
	Setup()
	defer Db.Close()
	migrateSchema()
	Start()
}
