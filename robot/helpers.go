package robot

import (
	"github.com/anaxita/mikrotik-wl-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Robot) roleByUsername(userName string) int {
	for name := range b.store.Admins {
		if name == userName {
			return roleAdmin
		}
	}

	return roleUser
}

func (b *Robot) isChatAllow(chatID int64) bool {
	for _, id := range b.store.AllowChatIDs {
		if id == chatID {
			return true
		}
	}

	return false
}

func (b *Robot) isAdmin(username string) bool {
	for admin := range b.store.Admins {
		if username == admin {
			return true
		}
	}

	return false
}

func (b *Robot) sendNotification(text string) error {
	msg := tgbotapi.NewMessage(storage.KmsMailChatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := b.api.Send(msg)

	return err
}
