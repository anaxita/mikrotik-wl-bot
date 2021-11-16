package robot

import (
	"github.com/anaxita/mikrotik-wl-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Robot) helpCommandHandler(_ *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "/start /help /add_ip /remove_ip")
}

func (b *Robot) removeIPCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusRemoveIP

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter ip in 192.168.1.25 format")
}

func (b *Robot) addIPCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusAddIP

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter ip in 192.168.1.25 format")
}

func (b *Robot) handleStartCommand(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusStart

	numericKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/add_ip"),
			tgbotapi.NewKeyboardButton("/remove_ip"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select a command")
	msg.ReplyMarkup = numericKeyboard

	return msg
}
