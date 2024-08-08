package handlers

import (
	"fmt"
	"log"
	"main/keyboards"
	"main/models"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func HandleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *gorm.DB) {
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
				UserID:       user.UserID,
				ReferralID:   referrerID,
				ReferralName: user.Username,
				ReferredBy:   referrerID,
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
	} else {
		if referrerID != 0 && referrerID != user.UserID {
			var referral models.Referral
			db.First(&referral, "user_id = ?", user.UserID)
			if referral.ID == 0 {
				referral = models.Referral{
					UserID:       user.UserID,
					ReferralID:   referrerID,
					ReferralName: user.Username,
					ReferredBy:   referrerID,
				}
				db.Create(&referral)

				var referrerUser models.User
				db.First(&referrerUser, "user_id = ?", referrerID)
				if referrerUser.ID != 0 {
					referrerUser.ReferralCount++
					db.Save(&referrerUser)
				}
			}
		}
	}

	// Сообщение с кнопкой "Меню"
	menuText := fmt.Sprintf("%s, мы на связи и готовы помочь☺️", update.Message.From.FirstName)
	menuMsg := tgbotapi.NewMessage(update.Message.Chat.ID, menuText)
	menuMsg.ReplyMarkup = keyboards.MenuButtonKeyboard()

	if _, err := bot.Send(menuMsg); err != nil {
		log.Println("Error sending menu message:", err)
	}

	// Сообщение с основными кнопками
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
