package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BackButton() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_partner"),
		),
	)
}

func MainInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Партнерская программа", "partner"),
			tgbotapi.NewInlineKeyboardButtonData("Способы обмена", "exchange"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Связаться с менеджером", "manager"),
			tgbotapi.NewInlineKeyboardButtonData("Запросить курс", "request"),
		),
	)
}

func PartnerProgramKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список рефералов", "referral_list"),
			tgbotapi.NewInlineKeyboardButtonData("Список начислений", "payment_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Как это работает", "how_it_works"),
			tgbotapi.NewInlineKeyboardButtonData("Вывод бонуса", "withdraw_bonus"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("QR-код партнерской ссылки", "qr_code"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️Назад", "back"),
		),
	)
}

func MenuButtonKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Меню"),
		),
	)
}

func ExchangeOptionsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Банкомат", "atm"),
			tgbotapi.NewInlineKeyboardButtonData("Офис", "office"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Курьер", "courier"),
			tgbotapi.NewInlineKeyboardButtonData("На счет", "account"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️Назад", "back"),
		),
	)
}

func ATMOptionsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обменять", "exchange_now"),
			tgbotapi.NewInlineKeyboardButtonData("Инструкция", "instruction"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Карта банкоматов", "atm_map"),
			tgbotapi.NewInlineKeyboardButtonData("⬅️Назад", "backToOptions"),
		),
	)
}

func OptionsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обменять", "exchange_now"),
			tgbotapi.NewInlineKeyboardButtonData("⬅️Назад", "backToOptions"),
		),
	)
}
