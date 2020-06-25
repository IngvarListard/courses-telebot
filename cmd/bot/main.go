package main

import (
	"context"
	"github.com/IngvarListard/courses-telebot/internal/coursesbot"
	"github.com/IngvarListard/courses-telebot/internal/store"
	"github.com/IngvarListard/courses-telebot/internal/store/gormstore"
	"log"
	"os"
	"os/signal"
)

func main() {

	config, err := coursesbot.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := store.NewDB()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error during closing the DB connection: %v", err)
		}
	}()

	gormStore := gormstore.New(db)

	bot, err := coursesbot.New(config, gormStore)
	if err != nil {
		log.Fatalf("error when setup the bot: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	if err := bot.Start(ctx); err != nil {
		log.Fatalf("program runtime error: %v", err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)

	for {
		sig := <-sigCh
		switch sig {
		case os.Interrupt:
			cancel()
			return
		}
	}

}
