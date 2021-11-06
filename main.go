package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net"
	"sync"
)

const (
	routerLogin    = "admin"
	routerPassword = "{tkgVb1"
	routerAddr     = "192.168.88.1:8728"
)

const (
	botToken = "2102715403:AAE13UIEQfDr7ZwUFubz2oVYo2P-knBK-sE"

	commandStart    = "start"
	commandAddIP    = "add_ip"
	commandRemoveIP = "remove_ip"
)

const (
	statusDefault = iota
	statusStart
	statusAddIP
	statusRemoveIP
)

// Sender provides method for handling messages from the bot
type Sender struct {
	mux    sync.Mutex
	bot    *tgbotapi.BotAPI
	users  map[int64]*User
	router *RouterController
}

// User is
type User struct {
	LastMessageID int
	Status        int
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/add_ip"),
		tgbotapi.NewKeyboardButton("/remove_ip"),
	),
)

func main() {
	routerController := NewRouterController()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	sender := &Sender{
		mux:    sync.Mutex{},
		bot:    bot,
		users:  make(map[int64]*User),
		router: routerController,
	}

	sender.start(updates)
}

func (s *Sender) handleCommands(update tgbotapi.Update) {
	var chatID = update.Message.Chat.ID
	var msg tgbotapi.MessageConfig
	var status int

	switch update.Message.Command() {
	case commandStart:
		msg = s.handleStartCommand(update)
		status = statusStart
	case commandRemoveIP:
		msg = tgbotapi.NewMessage(chatID, "Enter ip in XXX.XXX.XXX.XXX format")
		status = statusRemoveIP
	case commandAddIP:
		msg = tgbotapi.NewMessage(chatID, "Enter ip in XXX.XXX.XXX.XXX format")
		status = statusAddIP
	default:
		msg = tgbotapi.NewMessage(chatID, "Send /start to begin")
	}

	user, ok := s.users[chatID]

	if !ok {
		msg = tgbotapi.NewMessage(chatID, "Send /start to begin")
	} else {
		user.Status = status
		user.LastMessageID = update.Message.MessageID
	}

	if _, err := s.bot.Send(msg); err != nil {
		log.Println("send a message:", err)
	}
}

func (s *Sender) handleStartCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Select a command")
	msg.ReplyMarkup = numericKeyboard

	s.users[chatID] = &User{
		LastMessageID: update.Message.MessageID,
		Status:        statusStart,
	}

	return msg
}

func (s *Sender) handleMessages(update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig
	var chatID = update.Message.Chat.ID
	var text = update.Message.Text

	user, ok := s.users[chatID]
	if !ok {
		msg = tgbotapi.NewMessage(chatID, "Send /start to begin")
	} else {
		switch user.Status {
		case statusStart:
			msg = tgbotapi.NewMessage(chatID, "Please select a command and click on it")
		case statusRemoveIP:
			ip := net.ParseIP(text)
			if ip == nil {
				msg = tgbotapi.NewMessage(chatID, "Incorrect IP address. It should be XXX.XXX.XXX.XXX. Try again.")
			} else {
				err := s.router.RemoveIP(ip)
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
				msg = tgbotapi.NewMessage(chatID, "Incorrect IP address. It should be XXX.XXX.XXX.XXX. Try again.")
			} else {
				err := s.router.AddIP(ip)
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

	}

	if _, err := s.bot.Send(msg); err != nil {
		log.Println("send a message:", err)
	}
}

func (s *Sender) start(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			s.handleCommands(update)
			continue
		}

		s.handleMessages(update)
	}
}
