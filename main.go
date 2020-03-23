package main

import (
	"./bot"
	"./conf"
	"./users"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// migrateSchema синхронизация сземы БД
func migrateSchema() {
	conf.Db.AutoMigrate(&users.User{})
}

func main() {
	conf.Setup()
	defer conf.Db.Close()
	migrateSchema()
	bot.Start()
}
