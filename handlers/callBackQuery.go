package handlers

import (
	"fmt"
	"log"
	"main/database"
	"main/keyboards"
	"main/models"
	"main/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCallbackQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	callback := update.CallbackQuery
	callbackData := callback.Data
	chatID := callback.Message.Chat.ID

	var msgText string
	var replyMarkup tgbotapi.InlineKeyboardMarkup

	switch callbackData {
	case "exchange":
		msgText = "Выберите удобный способ получения:"
		replyMarkup = keyboards.ExchangeOptionsKeyboard()
	case "partner":
		msgText, replyMarkup = handlePartnerProgram(callback)
	case "back":
		msgText = "Выберите интересующую вас тему:"
		replyMarkup = keyboards.MainInlineKeyboard()
	case "backToOptions":
		msgText = "Выберите удобный способ получения:"
		replyMarkup = keyboards.ExchangeOptionsKeyboard()
	case "atm":
		msgText = atmMsg
		replyMarkup = keyboards.ATMOptionsKeyboard()
	case "office":
		msgText = officeMsg
		replyMarkup = keyboards.OptionsKeyboard()
	case "courier":
		msgText = courierMsg
		replyMarkup = keyboards.OptionsKeyboard()
	case "account":
		msgText = accountMsg
		replyMarkup = keyboards.OptionsKeyboard()
	default:
		return
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = replyMarkup

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending callback response message:", err)
	}

	// Удаляем сообщение с кнопками, чтобы не было дубликатов
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Println("Error deleting callback message:", err)
	}
}

func handlePartnerProgram(callback *tgbotapi.CallbackQuery) (string, tgbotapi.InlineKeyboardMarkup) {
	userID := callback.From.ID
	var user models.User
	database.DB.First(&user, "user_id = ?", userID)

	if user.ID == 0 {
		user = models.User{
			Username: callback.From.UserName,
			UserID:   userID,
		}
		database.DB.Create(&user)
	}

	referralLink := utils.GenerateReferralLink(user.UserID)
	msgText := fmt.Sprintf(
		"Ваш ID: %d\nКоличество рефералов: %d\nВаш баланс: %.2f\nВаша партнерская ссылка: %s",
		user.UserID, user.ReferralCount, user.BonusToWithdraw, referralLink,
	)
	replyMarkup := keyboards.PartnerProgramKeyboard()

	return msgText, replyMarkup
}
