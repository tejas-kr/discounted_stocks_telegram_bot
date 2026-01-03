package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var WorkerBaseURL string
var TelegramToken string

func init() {
	godotenv.Load("dev.env")
	WorkerBaseURL = os.Getenv("BASE_URL")
	if WorkerBaseURL == "" {
		log.Fatal("BASE_URL environment variable is required but not set")
	}
	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN environment variable is required but not set")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		handleCommand(bot, update.Message)
	}
}
