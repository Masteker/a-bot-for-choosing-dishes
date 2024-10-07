package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    tgbotapi "github.com/Syfaro/telegram-bot-api"
    _ "github.com/mattn/go-sqlite3" // Импорт SQLite драйвера
)

var botToken = "7317495569:AAEGfPna-0UwVwMAB2rgs8zLPASqt8jLO7g"
var dbFile = "recipes.db" // Имя файла базы данных SQLite

func main() {
	// Создаем базу данных, если она не существует
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		initDB()
	}

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

func initDB() {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
		return
	}
	defer db.Close()

	// Создаем таблицу рецептов
	createTableSQL := `CREATE TABLE recipes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
		return
	}

	// Можно добавить несколько тестовых данных, если хотите
	insertSQL := `INSERT INTO recipes (name) VALUES
        ('Борщ'),
        ('Салат Оливье'),
        ('Панкейки'),
        ('Паста карбонара');`
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Fatal("Ошибка при вставке данных:", err)
	}
}

func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) {
	receptSearchButton := tgbotapi.NewKeyboardButton("Искать рецепт")
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{receptSearchButton},
		},
		ResizeKeyboard: true,
	}

	msg := tgbotapi.NewMessage(chatID, "Привет! Чем могу помочь?")
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func findRecipe(bot *tgbotapi.BotAPI, chatID int64, query string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка подключения к базе данных. Пожалуйста, попробуйте позже."))
		return
	}
	defer db.Close()

	// Поиск рецептов
	querySQL := "SELECT name FROM recipes WHERE name LIKE ? LIMIT 3"
	rows, err := db.Query(querySQL, "%"+query+"%")
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при выполнении запроса к базе данных. Пожалуйста, попробуйте позже."))
		return
	}
	defer rows.Close()

	var recipes []string
	for rows.Next() {
		var recipe string
		if err := rows.Scan(&recipe); err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении данных из базы."))
			return
		}
		recipes = append(recipes, recipe)
	}

	// Отправляем найденные рецепты
	if len(recipes) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Рецепты не найдены."))
		return
	}

	messageText := "Найденные рецепты:\n"
	for _, recipe := range recipes {
		messageText += "- " + recipe + "\n"
	}

	bot.Send(tgbotapi.NewMessage(chatID, messageText))
}
