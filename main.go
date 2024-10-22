package main

import (
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var botToken = "7317495569:AAEGfPna-0UwVwMAB2rgs8zLPASqt8jLO7g"

// Это карта для хранения рецептов
var recipes = map[string]string{
	"Борщ":            "Ингредиенты: свекла, капуста, морковь, картофель, мясо, специи. Приготовление: Нарезать овощи, варить в бульоне до готовности.",
	"Салат Оливье":    "Ингредиенты: картофель, морковь, зеленый горошек, яйца, колбаса, майонез. Приготовление: Отварить всё, нарезать и смешать.",
	"Панкейки":        "Ингредиенты: мука, яйца, молоко, сахар, разрыхлитель. Приготовление: Смешать все ингредиенты и жарить на сковороде.",
	"Паста карбонара": "Ингредиенты: макароны, яйца, сыр Пармезан, бекон, черный перец. Приготовление: Сварить пасту, смешать с другими ингредиентами.",
}

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
		if update.Message == nil { // Игнорируем не сообщения
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
	receptSearchButton := tgbotapi.NewKeyboardButton("Искать рецепт")
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       [][]tgbotapi.KeyboardButton{{receptSearchButton}},
		ResizeKeyboard: true,
	}

	msg := tgbotapi.NewMessage(chatID, "Привет! Чем могу помочь?")
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func findRecipe(bot *tgbotapi.BotAPI, chatID int64, query string) {
	// Поиск рецептов в мапе
	var foundRecipes []string
	for name, recipe := range recipes {
		if contains(name, query) {
			foundRecipes = append(foundRecipes, name+": "+recipe)
		}
	}

	// Отправляем найденные рецепты
	if len(foundRecipes) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Рецепты не найдены."))
		return
	}

	messageText := "Найденные рецепты:\n"
	for _, recipe := range foundRecipes {
		messageText += "- " + recipe + "\n"
	}

	bot.Send(tgbotapi.NewMessage(chatID, messageText))
}

// Функция для проверки наличия строки в названии блюда
func contains(name, query string) bool {
	return len(query) > 0 && (name == query || contains(name, query))
}
