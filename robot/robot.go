package robot

import (
	"fmt"
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
	statusAddAdmin
	statusRemoveAdmin
)
const (
	roleUser = iota
	roleAdmin
)

const (
	commandStart           = "start"
	commandHelp            = "help"
	commandAddIP           = "add_ip"
	commandRemoveIP        = "remove_ip"
	commandAddAdmin        = "add_admin"
	commandRemoveAdmin     = "remove_admin"
	commandShowAdmins      = "show_admins"
	commandShowDynamicLink = "show_dynamic_link"
)

const (
	answerSuccess = "success"
)

type Robot struct {
	mux       sync.Mutex
	api       *tgbotapi.BotAPI
	store     *storage.Storage
	router    *router.Router
	dynamicWL string
}

func NewBot(token, dynamicWL string, store *storage.Storage, router *router.Router) (*Robot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Robot{
		mux:       sync.Mutex{},
		api:       bot,
		store:     store,
		router:    router,
		dynamicWL: dynamicWL,
	}, nil

}

func (b *Robot) handleCommands(update tgbotapi.Update) {
	var chatID = update.Message.Chat.ID
	var msg tgbotapi.MessageConfig

	user := b.store.Users[chatID]

	user.LastMessageID = update.Message.MessageID

	switch update.Message.Command() {
	case commandStart:
		msg = b.startCommandHandler(user, &update)
	case commandHelp:
		msg = b.helpCommandHandler(user, &update)
	case commandRemoveIP:
		msg = b.removeIPCommandHandler(user, &update)
	case commandAddIP:
		msg = b.addIPCommandHandler(user, &update)
	case commandShowAdmins:
		msg = b.ShowAdmins(user, &update)
	case commandAddAdmin:
		msg = b.addAdminCommandHandler(user, &update)
	case commandRemoveAdmin:
		msg = b.removeAdminCommandHandler(user, &update)
	case commandShowDynamicLink:
		msg = b.showDynamicLinkCommandHandler(user, &update)
	default:
		msg = tgbotapi.NewMessage(chatID, "Unknown command. Send /help to see all allow commands")
	}

	msg.DisableNotification = true

	if _, err := b.api.Send(msg); err != nil {
		log.Println("[ERROR] Can't send a message: ", err)
	}
}

func (b *Robot) handleMessages(update tgbotapi.Update) {
	var chatID = update.Message.Chat.ID
	var text = update.Message.Text
	var msgText = answerSuccess // default

	user := b.store.Users[chatID]

	switch user.Status {
	case statusStart:
		msgText = "Пожалуйста, выберите нужную команду на клавиатуре."
	case statusRemoveIP:
		ip := net.ParseIP(text)
		if ip == nil {
			msgText = "Некорректный IP. Введите адрес в формате 127.0.0.1."

			break
		}

		err := b.router.RemoveIP(ip)
		if err != nil {
			msgText = err.Error()
		}
		user.Status = statusStart
	case statusAddIP:
		ip := net.ParseIP(text)
		if ip == nil {
			msgText = "Некорректный IP. Введите адрес в формате 127.0.0.1."
		} else {
			chatTitle := update.Message.Chat.Title
			firstName := update.Message.From.FirstName
			lastName := update.Message.From.LastName

			comment := fmt.Sprintf("BOT %s | %s %s", chatTitle, firstName, lastName)
			err := b.router.AddIP(ip, comment)
			if err != nil {
				msgText = err.Error()
			}

			user.Status = statusStart
		}
	case statusAddAdmin:
		b.store.AddAdmin(text)
		user.Status = statusStart
	case statusRemoveAdmin:
		b.store.RemoveAdmin(text)
		user.Status = statusStart
	default:
		msgText = "Отправьте /start для начала работы."
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.DisableNotification = true

	if _, err := b.api.Send(msg); err != nil {
		log.Println("send a message:", err)
	}
}

func (b *Robot) addUser(update tgbotapi.Update) *storage.User {
	user := &storage.User{
		ID:            update.Message.From.ID,
		Username:      update.Message.From.UserName,
		LastMessageID: update.Message.MessageID,
		Status:        statusDefault,
		Role:          b.roleByUsername(update.Message.From.UserName),
	}

	b.mux.Lock()
	defer b.mux.Unlock()

	b.store.Users[update.Message.Chat.ID] = user

	return user
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

		chatID := update.Message.Chat.ID

		if !b.isChatAllow(chatID) && !b.isAdmin(update.Message.From.UserName) {
			msg := tgbotapi.NewMessage(chatID, "Доступ запрещен!\nДля получения прав, напишите https://t.me/Mishagl")
			_, err := b.api.Send(msg)
			if err != nil {
				log.Println("Can't send a message: ", err)
			}
			continue
		}

		_, ok := b.store.Users[chatID]
		if !ok {
			b.addUser(update)
		}

		if update.Message.IsCommand() {
			b.handleCommands(update)
			continue
		}

		b.handleMessages(update)

	}
}
