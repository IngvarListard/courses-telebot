package main

import (
	"github.com/IngvarListard/courses-telebot/internal/store"
	"github.com/IngvarListard/courses-telebot/internal/store/gormstore"
	"log"
)

func main() {
	db, err := store.NewDB()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error during closing the DB connection: %v", err)
		}
	}()
	gormstore := gormstore.New(db)
	gormstore.MigrateSchema()
}
