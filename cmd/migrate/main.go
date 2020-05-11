package main

import (
	"github.com/IngvarListard/courses-telebot/internal/db"
	"log"
)

func main() {
	_db, err := db.Setup()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}
	defer func() {
		if err := _db.Close(); err != nil {
			log.Printf("error during closing the DB connection: %v", err)
		}
	}()
	db.MigrateSchema()
}
