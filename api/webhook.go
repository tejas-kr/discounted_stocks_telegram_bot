package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func Handler(w http.ResponseWriter, r *http.Request) {
	bot, err := tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if update.Message == nil || !update.Message.IsCommand() {
		w.WriteHeader(http.StatusOK)
		return
	}

	handleCommand(bot, update.Message)
	w.WriteHeader(http.StatusOK)
}

func getAllDiscountedStocks(chatId int64) bool {
	url := fmt.Sprintf("%s/discounted_stocks?telegram_chat_id=%d", WorkerBaseURL, chatId)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func getAllStocks(chatId int64) bool {
	url := fmt.Sprintf("%s/all_stocks_status?telegram_chat_id=%d", WorkerBaseURL, chatId)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatId := message.Chat.ID
	msg := tgbotapi.NewMessage(chatId, "")

	switch message.Command() {
	case "discounted_stocks_all":
		success := getAllDiscountedStocks(chatId)
		if success {
			msg.Text = "🚀 Task started! I'll ping you here as soon as it gets completed."
		} else {
			msg.Text = "❌ Failed to reach the worker server."
		}
	case "stocks_all":
		success := getAllStocks(chatId)
		if success {
			msg.Text = "🚀 Task started! I'll ping you here as soon as it gets completed."
		} else {
			msg.Text = "❌ Failed to reach the worker server."
		}
	default:
		msg.Text = "I only know the following commands - /discounted_stocks_all, /stocks_all"
	}

	bot.Send(msg)
}
