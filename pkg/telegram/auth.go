package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mukolla/monobot/pkg/repository"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Response.Start)
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	return b.bot.Send(msg)
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessToken)
}

func (b *Bot) saveAccessToken(chatID int64, token string) (string, error) {
	if err := b.tokenRepository.Save(chatID, token, repository.AccessToken); err != nil {
		return "", err
	}
	return token, nil
}

func (b *Bot) getAuthToken(message *tgbotapi.Message) (string, error) {
	token, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		_, err = b.initAuthorizationProcess(message)
		if err != nil {
			return "", err
		}
	}
	return token, nil
}
