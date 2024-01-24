package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-bot/pkg/repository"
)

func (b *Bot) initAuthorization(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	_, err = b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.config.Messages.Start, authLink)))
	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessToken)
}

func (b *Bot) generateAuthorizationLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectURL(chatID)
	token, err := b.pocket.GetRequestToken(context.TODO(), redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, token, repository.RequestToken); err != nil {
		return "", err
	}

	return b.pocket.GetAuthorizationURL(token, redirectURL)
}

func (b *Bot) generateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.config.AuthServerUrl, chatID)
}
