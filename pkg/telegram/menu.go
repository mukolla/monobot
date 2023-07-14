package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mukolla/monobot/pkg/repository"
)

func (b *Bot) botAuthMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/"+commandList),
			tgbotapi.NewKeyboardButton("/"+commandBalance),
		),
	)
	return keyboard
}

func (b *Bot) botAccountMenu(userInfo repository.UserInfo) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard()

	for _, account := range userInfo.Accounts {
		buttonText := "/list " + account.MaskedPan[0]
		button := tgbotapi.NewKeyboardButton(buttonText)
		keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(button))
	}

	return keyboard
}

func (b *Bot) sendMenuMessage(replyMarkup tgbotapi.ReplyKeyboardMarkup, messageText string, message *tgbotapi.Message) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
	msg.ReplyMarkup = replyMarkup
	return b.bot.Send(msg)
}
