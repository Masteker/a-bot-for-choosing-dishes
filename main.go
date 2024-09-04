package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Syfaro/telegram-bot-api"
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
		if update.Message == nil { // игнорируем не Message обновления
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
	buttons := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton("Искать Рецепт")},
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	msg := tgbotapi.NewMessage(chatID, "Привет, что ты хочешь приготовить?")
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func handleSearch(bot *tgbotapi.BotAPI, chatID int64, searchQuery string) {
	if searchQuery == "Искать Рецепт" {
		bot.Send(tgbotapi.NewMessage(chatID, "Введите название рецепта"))
		return
	}

	url := fmt.Sprintf("https://gotovim-doma.ru/wp-json/wp/v2/posts?search=%s", searchQuery)
	resp, err := http.Get(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при запросе к сайту: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка: не удалось получить данные от сайта (статус: %d)", resp.StatusCode)))
		return
	}

	var results []struct {
		Title struct {
			Rendered string `json:"rendered"`
		} `json:"title"`
		Link string `json:"link"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при обработке ответа: "+err.Error()))
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
