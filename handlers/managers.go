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

var (
	ActiveDialogs = make(map[int]ActiveDialog)
	NextDialogID  int64
)

type ActiveDialog struct {
	ManagerChatID int64
	LastActivity  int64
	DialogID      int64
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
			DialogID:      dialog.DialogID,
		}
	} else if userID, isInDialogue := findUserByManagerChatID(chatID); isInDialogue {
		// Пересылаем сообщение от менеджера к пользователю
		msg := tgbotapi.NewMessage(int64(userID), update.Message.Text)
		bot.Send(msg)

		// Обновляем время последней активности
		ActiveDialogs[userID] = ActiveDialog{
			ManagerChatID: chatID,
			LastActivity:  time.Now().Unix(),
			DialogID:      dialog.DialogID,
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

	// Генерация нового уникального ID диалога
	dialogID := GenerateNextDialogID()

	// Обновление ActiveDialogs с новым DialogID
	ActiveDialogs[userID] = ActiveDialog{
		ManagerChatID: 0, // Установится при соединении с менеджером
		LastActivity:  time.Now().Unix(),
		DialogID:      dialogID,
	}

	// Формирование ссылки для начала диалога
	dialogLink := fmt.Sprintf("%s?start=dialog_%d_%d", botLink, userID, dialogID)

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
	if dialog, exists := ActiveDialogs[userID]; exists {
		if dialog.ManagerChatID != 0 && time.Now().Unix() < dialog.LastActivity+10 { // 1 час в секундах
			msg := tgbotapi.NewMessage(managerChatID, "Другой менеджер уже общается с этим пользователем.")
			bot.Send(msg)
			return
		}
		delete(ActiveDialogs, userID) // Удаляем старый диалог
	}

	// Устанавливаем диалог
	dialogID := GenerateNextDialogID()
	ActiveDialogs[userID] = ActiveDialog{
		ManagerChatID: managerChatID,
		LastActivity:  time.Now().Unix(),
		DialogID:      dialogID,
	}

	// Уведомляем менеджера
	msg := tgbotapi.NewMessage(managerChatID, fmt.Sprintf("Вы начали диалог с пользователем %d. Напишите сообщение для общения.", userID))
	bot.Send(msg)

	// Уведомляем пользователя
	userMsg := tgbotapi.NewMessage(int64(userID), "Менеджер присоединился к чату. Вы можете начать общение.")
	bot.Send(userMsg)

	// Добавляем кнопку завершения диалога для менеджера и пользователя
	managerEndDialogButton := tgbotapi.NewInlineKeyboardButtonData("Завершить диалог", "end_dialog")
	userEndDialogButton := tgbotapi.NewInlineKeyboardButtonData("Завершить диалог", "end_dialog")

	managerKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(managerEndDialogButton),
	)
	userKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(userEndDialogButton),
	)

	// Отправляем сообщение менеджеру
	managerMessage := tgbotapi.NewMessage(managerChatID, "Нажмите кнопку, чтобы завершить диалог.")
	managerMessage.ReplyMarkup = managerKeyboard
	bot.Send(managerMessage)

	// Отправляем сообщение пользователю
	userMessage := tgbotapi.NewMessage(int64(userID), "Нажмите кнопку, чтобы завершить диалог.")
	userMessage.ReplyMarkup = userKeyboard
	bot.Send(userMessage)
}

func HandleEndButton(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID

	// Проверка, является ли этот chatID менеджером или пользователем
	if userID, isManager := findUserByManagerChatID(chatID); isManager {
		// Менеджер завершает диалог
		endDialog(bot, userID)
	} else {
		// Пользователь завершает диалог
		if _, exists := ActiveDialogs[int(chatID)]; exists {
			endDialog(bot, int(chatID))
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

		// Уведомляем менеджера, если он еще в чате
		if dialog.ManagerChatID != 0 {
			managerMsg := tgbotapi.NewMessage(dialog.ManagerChatID, "Диалог с пользователем завершен.")
			bot.Send(managerMsg)
		}

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
			if currentTime-dialog.LastActivity > 30 {
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

func GenerateNextDialogID() int64 {
	NextDialogID++
	return NextDialogID
}
