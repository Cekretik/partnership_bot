package main

import (
	"log"
	"os"
	"strings"

	"main/database"
	"main/handlers"
	"main/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
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
	go handlers.MonitorDialogs(bot)
	go utils.StartUpdateRoutine(
		db,
		"1GcJ6b_j1ZPx-33mlVV9sIEPXraH82nt7MdsH6PpyAfk", "Партнерская программа1!A6:G",
		"1GcJ6b_j1ZPx-33mlVV9sIEPXraH82nt7MdsH6PpyAfk", "Партнерская программа!B6:B",
		"1GcJ6b_j1ZPx-33mlVV9sIEPXraH82nt7MdsH6PpyAfk", "CRM_партнерка!B3:D",
	)

	for update := range updates {
		if update.Message != nil {
			if strings.HasPrefix(update.Message.Text, "/start") {
				handlers.HandleStart(update, bot, db)
			} else {
				switch update.Message.Text {
				case "Меню":
					handlers.HandleMenu(update, bot)
				case "end":
					handlers.HandleEndCommand(bot, update)
				default:
					handlers.HandleMessages(bot, update)
					//handlers.HandleStart(update, bot, db)
				}
			}
		} else if update.CallbackQuery != nil {
			handlers.HandleCallbackQuery(update, bot)
		}
	}
}
