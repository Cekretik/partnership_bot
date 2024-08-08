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
	case "back_to_partner":
		msgText, replyMarkup = handlePartnerProgram(callback)
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
	case "qr_code":
		HandleQRCodeCallback(bot, update)
		return
	case "referral_list":
		msgText, replyMarkup = HandleReferals(callback, bot)
	case "how_it_works":
		msgText = howItWorksMsg
		replyMarkup = keyboards.BackButton()
	default:
		return
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = replyMarkup

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending callback message:", err)
	}

	// Удаляем сообщение с кнопками, чтобы не было дубликатов
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Println("Error deleting callback message:", err)
	}
}

func HandleQRCodeCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Chat.ID

	link := utils.GenerateReferralLink(int64(userID))
	qrCode, err := utils.GenerateQRCode(link)
	if err != nil {
		log.Println("Error generating QR code:", err)
		return
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  "qrcode.png",
		Bytes: qrCode,
	}
	photoMsg := tgbotapi.NewPhoto(chatID, fileBytes)
	photoMsg.Caption = "Ваш QR-код партнерской ссылки:"

	// Отправляем фото без изменения сообщения
	if _, err := bot.Send(photoMsg); err != nil {
		log.Println("Error sending QR code photo:", err)
		return
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

func HandleReferals(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) (string, tgbotapi.InlineKeyboardMarkup) {
	if callback.Message == nil || callback.From == nil {
		log.Println("Error: callback.Message or callback.From is nil")
		return "Ошибка: не удалось получить данные пользователя", keyboards.BackButton()
	}

	userID := callback.From.ID
	var referrals []models.Referral
	db := database.DB

	if db == nil {
		log.Println("Error: database connection is nil")
		return "Ошибка: не удалось подключиться к базе данных", keyboards.BackButton()
	}

	db.Where("referred_by = ?", userID).Find(&referrals)

	if len(referrals) == 0 {
		return "У вас нет рефералов", keyboards.BackButton()
	} else {
		msgText := "Ваши рефералы:"
		for _, r := range referrals {
			msgText += fmt.Sprintf("\n🆔: %d, \nИмя: %s, \nСумма обмена: %v", r.UserID, r.UserName, r.TradeAmount)
		}

		replyMarkup := keyboards.BackButton()
		return msgText, replyMarkup
	}
}
