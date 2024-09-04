package main

import (
	"encoding/json" // Не забудьте импортировать json
	"log"
	"net/http"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var botToken = "YOUR_BOT_TOKEN"

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // игнорируем не сообщения
			continue
		}

		switch update.Message.Text {
		case "/start":
			sendStartMessage(bot, update.Message.Chat.ID)
		default:
			if update.Message.Chat.IsPrivate() {
				findRecipe(bot, update.Message.Chat.ID, update.Message.Text)
			}
		}
	}
}

func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) {
	// Создаем кнопку
	receptSearchButton := tgbotapi.NewKeyboardButton("Искать рецепт")
	// Создаем срез с одной кнопкой
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{receptSearchButton}, // Запускаем срез с кнопками
		},
		ResizeKeyboard: true,
	}

	msg := tgbotapi.NewMessage(chatID, "Привет! Чем могу помочь?")
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func findRecipe(bot *tgbotapi.BotAPI, chatID int64, query string) {
	// Здесь будет логика поиска рецептов
	// Например, если у вас есть API для поиска рецептов
	response, err := http.Get("https://example.com/api/search?query=" + query)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при поиске рецепта. Пожалуйста, попробуйте позже."))
		return
	}
	defer response.Body.Close()

	var recipes []string // Предположим, ответ — это массив строк (названий рецептов)
	if err := json.NewDecoder(response.Body).Decode(&recipes); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка декодирования ответа."))
		return
	}

	// Отправляем найденные рецепты (например, первые 3)
	if len(recipes) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Рецепты не найдены."))
		return
	}

	messageText := "Найденные рецепты:\n"
	for _, recipe := range recipes[:3] { // Показываем только первые 3 рецепта
		messageText += "- " + recipe + "\n"
	}

	bot.Send(tgbotapi.NewMessage(chatID, messageText))
}
