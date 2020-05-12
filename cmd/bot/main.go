package main

import (
	"fmt"
	"github.com/IngvarListard/courses-telebot/internal/coursesbot"
	"github.com/IngvarListard/courses-telebot/internal/db"
	"log"
	"os"
)

var (
	APIKey string
	Debug  bool
)

// parseEnv bot configuration
func parseEnv() error {
	APIKey = os.Getenv("API_KEY")
	Debug = os.Getenv("DEBUG") == "1"
	if APIKey == "" {
		return fmt.Errorf("API_KEY is missiong")
	}
	return nil
}

func main() {
	if err := parseEnv(); err != nil {
		log.Fatal(err)
	}

	_db, err := db.Setup()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}

	defer func() {
		if err := _db.Close(); err != nil {
			log.Printf("error during closing the DB connection: %v", err)
		}
	}()

	if err := coursesbot.Setup(APIKey, Debug); err != nil {
		log.Fatalf("error when setup the bot: %v", err)
	}
	if err := coursesbot.Start(); err != nil {
		log.Fatalf("program runtime error: %v", err)
	}

}
