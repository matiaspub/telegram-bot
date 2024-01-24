package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/repository"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocket          *pocket.Client
	tokenRepository repository.TokenRepository
	config          *config.Config
}

func NewBot(bot *tgbotapi.BotAPI, pocket *pocket.Client, tr repository.TokenRepository, config *config.Config) *Bot {
	return &Bot{bot: bot, pocket: pocket, tokenRepository: tr, config: config}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	err = b.handleUpdates(updates)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	return updates, err
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		err := b.handleMessage(update.Message)
		if err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
	return nil
}
