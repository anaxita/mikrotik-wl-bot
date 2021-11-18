package main

import (
	"github.com/anaxita/mikrotik-wl-bot/robot"
	router2 "github.com/anaxita/mikrotik-wl-bot/router"
	"github.com/anaxita/mikrotik-wl-bot/storage"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	config, err := godotenv.Read()
	if err != nil {
		log.Fatalln(err)
	}

	router := router2.NewRouter(config["ROUTER_ADDR"], config["ROUTER_USER"], config["ROUTER_PASSWORD"])

	store := storage.NewStorage()

	bot, err := robot.NewBot(config["BOT_TOKEN"], config["DYNAMIC_WL"], store, router)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Start()
}
