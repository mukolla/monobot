package bank

import (
	"context"
	"errors"
	"fmt"
	"github.com/mukolla/monobot/pkg/repository"
	"github.com/shal/mono"
	"log"
	"os"
	"time"
)

type BankingRepository struct {
	bankName string
	token    string
}

func NewBankRepository(bankName string, token string) *BankingRepository {
	return &BankingRepository{bankName: bankName, token: token}
}

func (b *BankingRepository) Account() (repository.UserInfo, error) {

	personal := mono.NewPersonal(b.token)
	user, err := personal.User(context.Background())

	if err != nil {
		log.Fatal(err.Error())
		return repository.UserInfo{}, err
	}

	var userInfo repository.UserInfo

	userInfo.Name = user.Name

	for _, acc := range user.Accounts {
		rAccount := repository.Account{
			ID:           acc.ID,
			Balance:      acc.Balance,
			CurrencyCode: acc.CurrencyCode,
			IBAN:         acc.IBAN,
			MaskedPan:    acc.MaskedPan,
		}

		userInfo.Accounts = append(userInfo.Accounts, rAccount)
	}

	return userInfo, nil
}

func (b *BankingRepository) Balance() (string, error) {
	return "", nil
}

func (b *BankingRepository) TransactionList(accountID string, from time.Time, to time.Time) ([]repository.Transaction, error) {

	personal := mono.NewPersonal(b.token)

	user, err := personal.User(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var account mono.Account

	for _, acc := range user.Accounts {
		if acc.ID == accountID {
			account = acc
		}
	}

	if account.ID == "" {
		return nil, errors.New("Account not found by token")
	}

	transactions, err := personal.Transactions(context.Background(), account.ID, from, to)

	if err != nil {
		return nil, err
	}

	myTransactions := make([]repository.Transaction, len(transactions))

	for i, trx := range transactions {
		myTransaction := repository.Transaction{
			Time:        repository.Time(trx.Time),
			Description: trx.Description,
			Amount:      trx.Amount,
			Balance:     trx.Balance,
		}
		myTransactions[i] = myTransaction
	}

	return myTransactions, nil
}
