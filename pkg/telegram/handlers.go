package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/mono"
	"log"
	"time"
)

const commandStart = "start"
const commandBalance = "balance"
const commandList = "list"

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	log.Printf("[%s] %s", message.From.UserName, message.Text)

	if message.ReplyToMessage != nil && message.ReplyToMessage.Text == b.messages.Response.Start {
		token := message.Text
		_, err := b.saveAccessToken(message.Chat.ID, token)

		if err != nil {
			return err
		}

		return b.sendMessage(message, b.messages.Response.SavedSuccessfully)
	}

	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	reply := "Привіт! Це бот для перевірки балансу. Виберіть одну з опцій:"
	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	msg.ReplyMarkup = createMainMenuKeyboard()

	_, err = b.bot.Send(msg)
	return err

	//return b.sendMessage(message, b.messages.Response.AlreadyUsed)
}

func createMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/"+commandList),
			tgbotapi.NewKeyboardButton("/"+commandBalance),
		),
	)
	return keyboard
}

func (b *Bot) sendMessage(message *tgbotapi.Message, messageText string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandList:
		return b.handleBalanceCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	return errorAuthorized
}

func (b *Bot) handleBalanceCommand(message *tgbotapi.Message) error {
	token, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	personal := mono.NewPersonal(token)

	user, err := personal.User(context.Background())
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	from := time.Now().Add(-730 * time.Hour)
	to := time.Now()

	var account mono.Account

	for _, acc := range user.Accounts {
		ccy, _ := mono.CurrencyFromISO4217(acc.CurrencyCode)
		if ccy.Code == "UAH" {
			account = acc
		}
	}

	transactions, err := personal.Transactions(context.Background(), account.ID, from, to)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	var messageText string
	messageText += fmt.Sprintf("Account: %s\n", account.ID)
	messageText += fmt.Sprintf("Balance: %s\n", account.Balance)
	messageText += fmt.Sprintf("MaskedPan: %s\n", account.MaskedPan)
	messageText += "Transactions:\n"
	for _, transaction := range transactions {
		messageText += fmt.Sprintf("%d\t%s\n", transaction.Amount/100, transaction.Description)
	}
	return b.sendMessage(message, messageText)
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Response.UnknownCommand+" ["+message.Command()+"]")
	_, error := b.bot.Send(msg)
	return error
}
