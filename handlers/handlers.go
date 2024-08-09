package handlers

import (
	"fmt"
	"log"
	"main/database"
	"main/keyboards"
	"main/models"
	"main/utils"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func HandleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
	arg := update.Message.CommandArguments()
	chatID := update.Message.Chat.ID

	if strings.HasPrefix(arg, "dialog_") {
		userIDStr := strings.TrimPrefix(arg, "dialog_")
		userID, err := strconv.Atoi(userIDStr)
		if err == nil {
			StartDialog(bot, chatID, userID)
			return
		}
	}

	args := strings.Split(update.Message.Text, " ")
	referrerID := int64(0)

	if len(args) > 1 {
		referrerID, _ = strconv.ParseInt(args[1], 10, 64)
	}

	user := models.User{
		Username:   update.Message.From.UserName,
		UserID:     update.Message.From.ID,
		IncomeRate: 1,
	}

	var existingUser models.User
	db.First(&existingUser, "user_id = ?", user.UserID)

	if existingUser.ID == 0 {
		if referrerID != 0 && referrerID != user.UserID {
			referral := models.Referral{
				UserID:     user.UserID,
				UserName:   user.Username,
				ReferredBy: referrerID,
			}
			db.Create(&referral)

			var referrerUser models.User
			db.First(&referrerUser, "user_id = ?", referrerID)
			if referrerUser.ID != 0 {
				referrerUser.ReferralCount++
				db.Save(&referrerUser)
			}
		}

		db.Create(&user)
	}

	menuText := fmt.Sprintf("%s, мы на связи и готовы помочь☺️", update.Message.From.FirstName)
	menuMsg := tgbotapi.NewMessage(update.Message.Chat.ID, menuText)
	menuMsg.ReplyMarkup = keyboards.MenuButtonKeyboard()

	if _, err := bot.Send(menuMsg); err != nil {
		log.Println("Error sending menu message:", err)
	}

	mainMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите интересующую вас тему:")
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

func HandleWithdraw(callback *tgbotapi.CallbackQuery) (string, tgbotapi.InlineKeyboardMarkup) {
	if callback.Message == nil || callback.From == nil {
		log.Println("Error: callback.Message or callback.From is nil")
		return "Ошибка: не удалось получить данные пользователя", keyboards.BackButton()
	}
	userID := callback.From.ID
	var user models.User
	db := database.DB

	db.Where("user_id = ?", userID).First(&user)
	if user.ID == 0 {
		log.Println("Error: user not found")
		return "Ошибка: пользователь не найден", keyboards.BackButton()
	}
	if user.BonusToWithdraw == 0 {
		return "У вас нет бонусов для вывода", keyboards.BackButton()
	} else {
		msgText := fmt.Sprintf("💳Ваш бонус для вывода: %.2f", user.BonusToWithdraw)

		withdrawButton := tgbotapi.NewInlineKeyboardButtonData("Вывести", "withdraw_confirm")
		backButton := tgbotapi.NewInlineKeyboardButtonData("⬅️Назад", "back_to_partner")

		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(withdrawButton, backButton),
		)

		return msgText, replyMarkup
	}
}

func HandlePartnerProgram(callback *tgbotapi.CallbackQuery) (string, tgbotapi.InlineKeyboardMarkup) {
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
		"🔍Ваш ID: %d\n\n🤵‍♂️Количество рефералов: %d\n\n♻️Заработано всего: %.2f\n\n🔗Ваша партнерская ссылка: %s",
		user.UserID, user.ReferralCount, user.TotalBonus, referralLink,
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
		msgText := "Ваши рефералы:\n"
		for _, r := range referrals {
			msgText += fmt.Sprintf("\n🆔: %d, \nИмя: %s, \nСумма обмена: %v\n", r.UserID, r.UserName, r.TradeAmount)
		}

		replyMarkup := keyboards.BackButton()
		return msgText, replyMarkup
	}
}
