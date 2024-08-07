package main

import (
	"log"
	"main/database"
	"main/handlers"
	"os"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	database.Init()
	db := database.DB

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				handlers.HandleStart(update, bot, db)
			case "Меню":
				handlers.HandleMenu(update, bot)
			default:
			}
		} else if update.CallbackQuery != nil {
			handlers.HandleCallbackQuery(update, bot)
		}
	}
}
