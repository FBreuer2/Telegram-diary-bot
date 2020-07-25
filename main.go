package main

import (
	"log"
	"os"

	BotDatabase "github.com/FBreuer2/telegram-diary/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const DB_NAME = "test.sqlite"

var BOT_TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
var USER_NAME = os.Getenv("TELEGRAM_USER_NAME")

func main() {
	bot, err := tgbotapi.NewBotAPI(BOT_TOKEN)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	db, err := BotDatabase.New(DB_NAME)

	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	log.Printf("Connected to db %s\n", DB_NAME)

	if !db.Exists(USER_NAME) {
		db.Create(USER_NAME)
		log.Printf("Created user %s\n", USER_NAME)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.From.UserName != USER_NAME {
			continue
		}

		if update.Message.Photo != nil {
			log.Println("New photo received")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Photo saved")
			bot.Send(msg)
			continue
		}

		if update.Message.Text != "" { // ignore any non-Message Updates
			log.Println("New text received")
			db.AddText(update.Message.Text, USER_NAME)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Text saved")
			bot.Send(msg)
			continue
		}

		if update.Message.Location.Latitude != 0 && update.Message.Location.Longitude != 0 {
			log.Println("New location received")
			db.AddLocation(update.Message.Location, USER_NAME)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Location saved")
			bot.Send(msg)
			continue
		}
	}
}
