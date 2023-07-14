package telegram

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mukolla/monobot/pkg/repository"
	"github.com/mukolla/monobot/pkg/repository/bank"
	"log"
	"strconv"
	"time"
)

const commandStart = "start"
const commandBalance = "balance"
const commandList = "list"

func (b *Bot) handleMessage(message *tgbotapi.Message) (tgbotapi.Message, error) {
	if message.ReplyToMessage != nil && message.ReplyToMessage.Text == b.messages.Response.Start {
		token := message.Text
		_, err := b.saveAccessToken(message.Chat.ID, token)

		if err != nil {
			return tgbotapi.Message{}, err
		}

		return b.sendMenuMessage(b.botAuthMenu(), b.messages.Response.SavedSuccessfully, message)
	}

	_, err := b.getAccessToken(message.Chat.ID)

	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	return b.sendMenuMessage(b.botAuthMenu(), b.messages.Response.AlreadyUsed, message)
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
	case commandBalance:
		return b.handleBalanceCommand(message)
	case commandList:
		return b.handleTransactionListCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	return b.handleBalanceCommand(message)
}

func (b *Bot) handleBalanceCommand(message *tgbotapi.Message) error {

	token, _ := b.getAuthToken(message)
	bank := bank.NewBankRepository(string(repository.BankMonoBank), token)
	userInfo, err := bank.Account()

	if err != nil {
		return err
	}

	menuText := ""
	for _, account := range userInfo.Accounts {
		menuText += account.MaskedPan[0] + " - Balance: " + strconv.Itoa(account.Balance) + "\n"
	}

	_, err = b.sendMenuMessage(b.botAccountMenu(userInfo), menuText, message)
	if err != nil {
		return err
	}

	return err
}

func (b *Bot) handleTransactionListCommand(message *tgbotapi.Message) error {

	token, _ := b.getAuthToken(message)
	bank := bank.NewBankRepository(string(repository.BankMonoBank), token)

	now := time.Now()
	to := roundTo15Minutes(now)
	from := to.Add(-time.Hour * 730)

	userInfoKey := token + from.String() + to.String()
	userInfo, err := b.cash.UserInfoGet(userInfoKey)

	if err != nil || userInfo.Name == "" {
		userInfo, err = bank.Account()
		b.cash.UserInfoSave(userInfoKey, userInfo)

		log.Println("\n\n SAVE UserInfo TO cash \n\n")

	} else {
		log.Println("\n\n LOAD UserInfo IN cash \n\n")
	}

	if err != nil {
		return err
	}

	maskedPan := message.CommandArguments()

	for _, account := range userInfo.Accounts {
		if maskedPan == account.MaskedPan[0] {

			cacheKey := account.ID + from.String() + to.String()

			log.Println("\n Transaction, cacheKey: [" + calculateMD5(cacheKey) + "]")

			emptyExists, err := b.cash.TransactionEmptyExists(cacheKey)
			if err != nil {
				return err
			}

			var messageText string

			if emptyExists == true {
				log.Println("\n Transaction cash exist and empty")
				messageText += fmt.Sprintf("MaskedPan: %s\n", maskedPan)
				messageText += fmt.Sprintf("Transaction: %s\n", "Transaction cash exist and empty")
			} else {

				transactions, err := b.cash.TransactionGet(cacheKey)

				if err != nil {
					return err
				}

				if len(transactions) == 0 {
					transactions, _ = bank.TransactionList(account.ID, from, to)
					err := b.cash.TransactionSave(cacheKey, transactions)

					if err != nil {
						return err
					}

					log.Println("\n Transaction SAVED cash")

				} else {
					log.Println("\n Transaction LOAD cash")
				}

				messageText += fmt.Sprintf("Name: %s\n", userInfo.Name)
				messageText += fmt.Sprintf("Account: %s\n", maskedPan)
				messageText += fmt.Sprintf("Balance: %s\n", account.Balance)
				messageText += fmt.Sprintf("MaskedPan: %s\n", account.MaskedPan)
				messageText += "Transactions:\n"

				for _, transaction := range transactions {
					messageText += fmt.Sprintf("%d\t%s\n", transaction.Amount/100, transaction.Description)
				}
			}

			return b.sendMessage(message, messageText)
		}
	}

	return errorGetTransactionList
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Response.UnknownCommand+" ["+message.Command()+"]")
	_, error := b.bot.Send(msg)
	return error
}

func calculateMD5(cacheKey string) string {
	hashes := md5.New()
	hashes.Write([]byte(cacheKey))
	hash := hex.EncodeToString(hashes.Sum(nil))
	return hash
}

func roundTo15Minutes(t time.Time) time.Time {
	roundedMinutes := (t.Minute() / 15) * 15
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMinutes, 0, 0, t.Location())
}
