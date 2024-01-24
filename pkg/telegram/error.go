package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "🥴")
	switch {
	case errors.Is(err, errInvalidURL):
		msg.Text = b.config.Messages.InvalidUrl
	case errors.Is(err, errUnauthorized):
		msg.Text = b.config.Messages.Unauthorized
	case errors.Is(err, errUnableToSave):
		msg.Text = b.config.Messages.UnableToSave
	}
	_, _ = b.bot.Send(msg)
}
