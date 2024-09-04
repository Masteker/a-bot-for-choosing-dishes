package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var botToken = "7317495569:AAEGfPna-0UwVwMAB2rgs8zLPASqt8jLO7g"

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Ошибка при создании бота:", err)
		return
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // ignore non-Message Updates
			continue
		}

		switch update.Message.Text {
		case "/start":
			sendStartMessage(bot, update.Message.Chat.ID)
		default:
			if update.Message.Chat.IsPrivate() {
				handleSearch(bot, update.Message.Chat.ID, update.Message.Text)
			}
		}
	}
}

func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) {
	receptSearchButton := tgbotapi.NewKeyboardButton("Искать рецепт")
	keyboard := tgbotapi.NewReplyKeyboard(receptSearchButton)

	msg := tgbotapi.NewMessage(chatID, "ВОТ ТУТ ПРИВЕТСТВЕННОЕ СООБЩЕНИЕ БОТА ПОСЛЕ КНОПКИ СТАРТ")
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func handleSearch(bot *tgbotapi.BotAPI, chatID int64, searchQuery string) {
	if searchQuery == "Искать рецепт" {
		bot.Send(tgbotapi.NewMessage(chatID, "Введи, что ты хочешь приготовить."))
		return
	}

	url := fmt.Sprintf("Вhttps://gotovim-doma.ru/wp-json/wp/v2/posts?search=%s", searchQuery)
	resp, err := http.Get(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при запросе к сайту"))
		return
	}
	defer resp.Body.Close()

	var results []struct {
		Title struct {
			Rendered string `json:"rendered"`
		} `json:"title"`
		Link string `json:"link"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при обработке ответа"))
		return
	}

	if len(results) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Я ничего не нашел по вашему запросу"))
		return
	}

	for _, result := range results {
		completeMessage := fmt.Sprintf("%s %s", result.Title.Rendered, result.Link)
		bot.Send(tgbotapi.NewMessage(chatID, completeMessage))
	}
}
