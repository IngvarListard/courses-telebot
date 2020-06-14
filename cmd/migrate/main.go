package main

import (
	"github.com/IngvarListard/courses-telebot/internal/store"
	"log"
)

func main() {
	_db, err := store.NewDB()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}
	defer func() {
		if err := _db.Close(); err != nil {
			log.Printf("error during closing the DB connection: %v", err)
		}
	}()
	store.MigrateSchema()
}
