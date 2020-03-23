package bot

import (
	"../conf"
	"../users"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// Start ...
func Start() {

	bot, err := tgbotapi.NewBotAPI(conf.ApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.Text == "/start" {
			user := conf.Db.FirstOrCreate(&users.User{TgId: update.Message.From.ID}, users.User{
				Name: update.Message.From.FirstName,
				TgId: update.Message.From.ID,
			})
			fmt.Println(user)
		}

		bot.Send(msg)
	}
}
