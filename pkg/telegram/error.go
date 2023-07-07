package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	unknownError      = errors.New("unknown error")
	errorAuthorized   = errors.New("user is not Authorized")
	errorGetBalance   = errors.New("error get balance")
	authTokenNotFound = errors.New("auth token not found")
)

func (b *Bot) handleError(chatID int64, err error) error {

	msg := tgbotapi.NewMessage(chatID, "")

	switch err {
	case unknownError:
		msg.Text = b.messages.Errors.UnknownError
	case authTokenNotFound:
		msg.Text = b.messages.Errors.AuthTokenNotFound
	case errorAuthorized:
		msg.Text = b.messages.Errors.Unauthorized
	case errorGetBalance:
		msg.Text = b.messages.Errors.GetBalance
	default:
		msg.Text = b.messages.Errors.Default
	}

	b.bot.Send(msg)

	return err
}
