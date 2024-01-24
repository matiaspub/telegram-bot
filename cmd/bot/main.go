package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/repository/boltdb"
	"telegram-bot/pkg/server"
	"telegram-bot/pkg/telegram"
)

func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatal("config " + err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Panic("tg " + err.Error())
	}
	bot.Debug = true

	pock, err := pocket.NewClient(conf.PocketConsumerKey)
	if err != nil {
		log.Panic(err)
	}

	db, err := boltdb.InitDB(conf)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pock, tokenRepository, conf)

	authorizationServer := server.NewAuthorizationServer(pock, tokenRepository, conf.TelegramBotUrl)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}
