package main

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/db/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// migrateSchema синхронизация сземы БД
func migrateSchema() {
	conf.Db.AutoMigrate(&models.User{})
}

func main() {
	defer func() {
		if err := conf.Db.Close(); err != nil {
			fmt.Printf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()
	migrateSchema()
	Start()
}
