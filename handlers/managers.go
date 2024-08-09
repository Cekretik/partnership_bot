package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var ActiveDialogs = make(map[int]ActiveDialog)

type ActiveDialog struct {
	ManagerChatID int64
	LastActivity  int64
}

func HandleMessages(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	if dialog, isInDialogue := ActiveDialogs[int(userID)]; isInDialogue {
		// Пересылаем сообщение от пользователя к менеджеру
		msg := tgbotapi.NewMessage(dialog.ManagerChatID, update.Message.Text)
		bot.Send(msg)

		// Обновляем время последней активности
		ActiveDialogs[int(userID)] = ActiveDialog{
			ManagerChatID: dialog.ManagerChatID,
			LastActivity:  time.Now().Unix(),
		}
	} else if userID, isInDialogue := findUserByManagerChatID(chatID); isInDialogue {
		// Пересылаем сообщение от менеджера к пользователю
		msg := tgbotapi.NewMessage(int64(userID), update.Message.Text)
		bot.Send(msg)

		// Обновляем время последней активности
		ActiveDialogs[userID] = ActiveDialog{
			ManagerChatID: chatID,
			LastActivity:  time.Now().Unix(),
		}
	}
}

func HandleManagerRequest(bot *tgbotapi.BotAPI, update tgbotapi.Update, chatID int64, userID int) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botLink := os.Getenv("BOT_LINK")
	managerGroupChatID, err := strconv.ParseInt(os.Getenv("MANAGER_CHAT"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing MANAGER_CHAT environment variable:", err)
	}

	// Проверка, не находится ли уже менеджер в диалоге с этим пользователем
	if isUserInDialogue(userID) {
		msg := tgbotapi.NewMessage(chatID, "Вы уже находитесь в диалоге с менеджером.")
		bot.Send(msg)
		return
	}

	// Формирование ссылки для начала диалога
	dialogLink := fmt.Sprintf(botLink+"?start=dialog_%d", userID)

	// Сообщение для группы менеджеров
	notificationMsg := fmt.Sprintf("Пользователь @%s запросил общение с менеджером. Нажмите на ссылку для начала диалога: %s", update.CallbackQuery.From.UserName, dialogLink)
	msg := tgbotapi.NewMessage(managerGroupChatID, notificationMsg)
	bot.Send(msg)

	// Уведомление пользователя
	userMsg := tgbotapi.NewMessage(chatID, "Запрос на общение с менеджером отправлен.")
	bot.Send(userMsg)
}

func StartDialog(bot *tgbotapi.BotAPI, managerChatID int64, userID int) {
	// Проверка, не состоит ли уже другой менеджер в диалоге с этим пользователем
	if _, exists := ActiveDialogs[userID]; exists {
		msg := tgbotapi.NewMessage(managerChatID, "Другой менеджер уже общается с этим пользователем.")
		bot.Send(msg)
		return
	}

	// Устанавливаем диалог
	ActiveDialogs[userID] = ActiveDialog{
		ManagerChatID: managerChatID,
		LastActivity:  time.Now().Unix(),
	}

	// Уведомляем менеджера
	msg := tgbotapi.NewMessage(managerChatID, fmt.Sprintf("Вы начали диалог с пользователем %d. Напишите сообщение для общения.", userID))
	bot.Send(msg)

	// Уведомляем пользователя
	userMsg := tgbotapi.NewMessage(int64(userID), "Менеджер присоединился к чату. Вы можете начать общение.")
	bot.Send(userMsg)

	// Добавляем кнопку завершения диалога для обеих сторон
	endDialogButton := tgbotapi.NewInlineKeyboardButtonData("Завершить диалог", "end_dialog")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(endDialogButton),
	)

	// Отправляем сообщение менеджеру
	managerMessage := tgbotapi.NewMessage(managerChatID, "Нажмите кнопку, чтобы завершить диалог.")
	managerMessage.ReplyMarkup = keyboard
	bot.Send(managerMessage)

	// Отправляем сообщение пользователю
	userMessage := tgbotapi.NewMessage(int64(userID), "Нажмите кнопку, чтобы завершить диалог.")
	userMessage.ReplyMarkup = keyboard
	bot.Send(userMessage)
}

func HandleEndButton(bot *tgbotapi.BotAPI, chatID int64) {
	// Определяем, кто завершил диалог
	if userID, isManager := findUserByManagerChatID(chatID); isManager {
		// Менеджер завершает диалог
		endDialog(bot, userID)
	} else {
		// Пользователь завершает диалог
		for userID, dialog := range ActiveDialogs {
			if dialog.ManagerChatID == chatID {
				endDialog(bot, userID)
				break
			}
		}
	}
}

func isUserInDialogue(userID int) bool {
	_, exists := ActiveDialogs[userID]
	return exists
}

func endDialog(bot *tgbotapi.BotAPI, userID int) {
	if dialog, exists := ActiveDialogs[userID]; exists {
		delete(ActiveDialogs, userID)

		// Уведомляем менеджера
		msg := tgbotapi.NewMessage(dialog.ManagerChatID, "Диалог с пользователем завершен.")
		bot.Send(msg)

		// Уведомляем пользователя
		userMsg := tgbotapi.NewMessage(int64(userID), "Диалог с менеджером завершен.")
		bot.Send(userMsg)
	}
}

func MonitorDialogs(bot *tgbotapi.BotAPI) {
	for {
		time.Sleep(1 * time.Minute)

		currentTime := time.Now().Unix()
		for userID, dialog := range ActiveDialogs {
			if currentTime-dialog.LastActivity > 600 {
				endDialog(bot, userID)
			}
		}
	}
}

func HandleEndCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID, err := strconv.Atoi(update.Message.CommandArguments())
	if err == nil {
		endDialog(bot, userID)
	}
}

func findUserByManagerChatID(managerChatID int64) (int, bool) {
	for userID, chatID := range ActiveDialogs {
		if chatID.ManagerChatID == managerChatID {
			return userID, true
		}
	}
	return 0, false
}
