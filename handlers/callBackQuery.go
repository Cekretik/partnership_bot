package handlers

import (
	"log"
	"main/keyboards"
	"main/utils"
	"path/filepath"

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
		msgText, replyMarkup = HandlePartnerProgram(callback)
	case "back":
		msgText = "Выберите интересующую вас тему:"
		replyMarkup = keyboards.MainInlineKeyboard()
	case "back_to_partner":
		msgText, replyMarkup = HandlePartnerProgram(callback)
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
	case "instruction":
		sendInstructionPhotos(bot, update.CallbackQuery.Message.Chat.ID)
		return
	case "withdraw_bonus":
		msgText, replyMarkup = HandleWithdraw(callback)
	case "manager":
		HandleManagerRequest(bot, update, chatID, int(callback.From.ID))
		return
	case "end_dialog":
		HandleEndButton(bot, update)
		return
	default:
		return
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = replyMarkup
	msg.ParseMode = "HTML"
	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending callback message:", err)
	}

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

	if _, err := bot.Send(photoMsg); err != nil {
		log.Println("Error sending QR code photo:", err)
		return
	}
}

func sendInstructionPhotos(bot *tgbotapi.BotAPI, chatID int64) {
	imgDir := "./imgs"

	files, err := filepath.Glob(filepath.Join(imgDir, "*.PNG"))
	if err != nil {
		log.Println("Error reading image files:", err)
		return
	}

	if len(files) != 4 {
		log.Println("Expected 4 images, found:", len(files))
		return
	}

	var mediaGroup []interface{}

	for _, file := range files {
		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(file))
		mediaGroup = append(mediaGroup, photo)
	}

	msg := tgbotapi.NewMediaGroup(chatID, mediaGroup)
	_, err = bot.Request(msg)
	if err != nil {
		log.Println("Error sending media group:", err)
		return
	}
}
