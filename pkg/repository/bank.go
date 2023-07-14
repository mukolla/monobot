package repository

import (
	"strconv"
	"time"
)

type Bank string

type Time struct {
	time.Time
}

// MarshalJSON is used to convert the timestamp to JSON
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*t = Time{time.Unix(q, 0).UTC()}
	return nil
}

type Transaction struct {
	Time        Time   `json:"time"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Balance     int64  `json:"balance"`
}

const (
	BankMonoBank Bank = "monobank"
)

type BankRepository interface {
	Account() (UserInfo, error)
	Balance() (string, error)
	TransactionList(accountID string) ([]Transaction, error)
}

type Account struct {
	ID           string   `json:"id"`           // Account identifier.
	Balance      int      `json:"balance"`      // Balance is minimal units (cents).
	CurrencyCode int32    `json:"currencyCode"` // Currency code in ISO4217.
	IBAN         string   `json:"iban"`         // IBAN.
	MaskedPan    []string `json:"maskedPan"`
}

type UserInfo struct {
	Name       string    `json:"name"`       // User name.
	WebHookURL string    `json:"webHookUrl"` // URL for receiving new transactions.
	Accounts   []Account `json:"accounts"`   // List of available accounts
}
