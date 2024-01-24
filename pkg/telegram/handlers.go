package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	commandStart = "start"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		_, err := b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды 😟"))
		return err
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorization(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Messages.Start)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Messages.SavedSuccessfully)

	if _, err := url.ParseRequestURI(message.Text); err != nil {
		return errInvalidURL
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}

	err = b.pocket.Add(context.TODO(), pocket.AddInput{AccessToken: accessToken, URL: message.Text})
	if err != nil {
		return errUnableToSave
	}

	_, err = b.bot.Send(msg)
	return err
}
