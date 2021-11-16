package robot

import (
	"github.com/anaxita/mikrotik-wl-bot/router"
	"github.com/anaxita/mikrotik-wl-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net"
	"sync"
)

const (
	statusDefault = iota
	statusStart
	statusAddIP
	statusRemoveIP
)

const (
	commandStart    = "start"
	commandHelp     = "help"
	commandAddIP    = "add_ip"
	commandRemoveIP = "remove_ip"
)

type Robot struct {
	mux    sync.Mutex
	api    *tgbotapi.BotAPI
	store  *storage.Storage
	router *router.Router
}

func NewBot(token string, store *storage.Storage, router *router.Router) (*Robot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Robot{
		mux:    sync.Mutex{},
		api:    bot,
		store:  store,
		router: router,
	}, nil

}

func (b *Robot) handleCommands(update tgbotapi.Update) {
	var chatID = update.Message.Chat.ID
	var msg tgbotapi.MessageConfig

	user, ok := b.store.Users[chatID]
	if !ok {
		user = b.addUser(update)
	}

	user.LastMessageID = update.Message.MessageID

	switch update.Message.Command() {
	case commandStart:
		msg = b.handleStartCommand(user, &update)
	case commandHelp:
		msg = b.helpCommandHandler(user, &update)
	case commandRemoveIP:
		msg = b.removeIPCommandHandler(user, &update)
	case commandAddIP:
		msg = b.addIPCommandHandler(user, &update)
	default:
		msg = tgbotapi.NewMessage(chatID, "Unknown command. Send /start to begin")
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Println("[ERROR] Can't send a message: ", err)
	}
}

func (b *Robot) handleMessages(update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig
	var chatID = update.Message.Chat.ID
	var text = update.Message.Text

	user, ok := b.store.Users[chatID]
	if !ok {
		user = b.addUser(update)
	}

	switch user.Status {
	case statusStart:
		msg = tgbotapi.NewMessage(chatID, "Please select a command and click on it")
	case statusRemoveIP:
		ip := net.ParseIP(text)
		if ip == nil {
			msg = tgbotapi.NewMessage(chatID, "Incorrect IP addresb. It should be XXX.XXX.XXX.XXX. Try again.")
		} else {
			err := b.router.RemoveIP(ip)
			if err != nil {
				msg = tgbotapi.NewMessage(chatID, err.Error())
			} else {
				msg = tgbotapi.NewMessage(chatID, "Success.")
			}
			user.Status = statusStart
		}
	case statusAddIP:
		ip := net.ParseIP(text)
		if ip == nil {
			msg = tgbotapi.NewMessage(chatID, "Incorrect IP addresb. It should be XXX.XXX.XXX.XXX. Try again.")
		} else {
			err := b.router.AddIP(ip)
			if err != nil {
				msg = tgbotapi.NewMessage(chatID, err.Error())
			} else {
				msg = tgbotapi.NewMessage(chatID, "Success.")
			}

			user.Status = statusStart
		}
	default:
		msg = tgbotapi.NewMessage(chatID, "Send /start to begin")
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Println("send a message:", err)
	}
}

func (b *Robot) Start() {
	b.api.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.api.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommands(update)
			continue
		}

		b.handleMessages(update)
	}
}

func (b *Robot) addUser(update tgbotapi.Update) *storage.User {
	user := &storage.User{
		ID:            update.Message.From.ID,
		Username:      update.Message.From.UserName,
		LastMessageID: update.Message.MessageID,
		Status:        statusDefault,
	}

	b.store.Users[update.Message.Chat.ID] = user

	return user
}