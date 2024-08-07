package handlers

import (
	"fmt"
	"log"
	"main/keyboards"
	"main/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func HandleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	user := models.User{
		Username: update.Message.From.UserName,
		UserID:   update.Message.From.ID,
	}

	db.FirstOrCreate(&user, models.User{UserID: user.UserID})

	// Сообщение с кнопкой "Меню"
	menuText := fmt.Sprintf("%s, мы на связи и готовы помочь☺️", update.Message.From.FirstName)
	menuMsg := tgbotapi.NewMessage(update.Message.Chat.ID, menuText)
	menuMsg.ReplyMarkup = keyboards.MenuButtonKeyboard()

	if _, err := bot.Send(menuMsg); err != nil {
		log.Println("Error sending menu message:", err)
	}

	// Сообщение с основными кнопками
	mainMsgText := "Выберите интересующую вас тему:"
	mainMsg := tgbotapi.NewMessage(update.Message.Chat.ID, mainMsgText)
	mainMsg.ReplyMarkup = keyboards.MainInlineKeyboard()

	if _, err := bot.Send(mainMsg); err != nil {
		log.Println("Error sending main message:", err)
	}
}

func HandleMenu(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msgText := fmt.Sprintf("%s, мы на связи и готовы помочь☺️", update.Message.From.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = keyboards.MainInlineKeyboard()

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending menu message:", err)
	}

	mainMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	mainMsg.ReplyMarkup = keyboards.MainInlineKeyboard()
}
