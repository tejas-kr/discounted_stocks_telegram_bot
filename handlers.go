package main

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


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

