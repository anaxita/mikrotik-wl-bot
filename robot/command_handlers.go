package robot

import (
	"fmt"
	"github.com/anaxita/mikrotik-wl-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Robot) helpCommandHandler(_ *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "/start\n/help\n/add_ip\n/remove_ip")
}

func (b *Robot) removeIPCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusRemoveIP

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter ip in 192.168.1.25 format")
}

func (b *Robot) addIPCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusAddIP

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter ip in 192.168.1.25 format")
}

func (b *Robot) startCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusStart

	numericKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/"+commandAddIP),
			tgbotapi.NewKeyboardButton("/"+commandRemoveIP),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/"+commandAddAdmin),
			tgbotapi.NewKeyboardButton("/"+commandRemoveAdmin),
			tgbotapi.NewKeyboardButton("/"+commandShowAdmins),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select a command")
	msg.ReplyMarkup = numericKeyboard

	return msg
}

func (b *Robot) addAdminCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusAddAdmin

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter admin username without @ to add")
}

func (b *Robot) removeAdminCommandHandler(user *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	user.Status = statusRemoveAdmin

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter admin username without @ to delete")
}

func (b *Robot) ShowAdmins(_ *storage.User, update *tgbotapi.Update) tgbotapi.MessageConfig {
	var text string

	admins := b.store.Admins

	for admin := range admins {
		text = fmt.Sprintf("%s@%s\n", text, admin)
	}
	return tgbotapi.NewMessage(update.Message.Chat.ID, text)
}
