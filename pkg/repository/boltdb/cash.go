package boltdb

import (
	"bytes"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/mukolla/monobot/pkg/repository"
)

type CashRepository struct {
	db                *bolt.DB
	TransactionBucket []byte
	UserInfoBucket    []byte
}

func NewCashRepository(db *bolt.DB) *CashRepository {
	return &CashRepository{
		db:                db,
		TransactionBucket: []byte(repository.TransactionCash),
		UserInfoBucket:    []byte(repository.UserInfoCash),
	}
}

func (t *CashRepository) TransactionGet(cacheKey string) ([]repository.Transaction, error) {
	var transactions []repository.Transaction

	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.TransactionBucket)
		value := b.Get(stringToBytes(cacheKey))
		if value != nil {
			var err error
			transactions, err = decodeTransactions(value)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *CashRepository) TransactionSave(cacheKey string, transaction []repository.Transaction) error {
	if len(transaction) == 0 {
		return t.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(t.TransactionBucket)
			return b.Put(stringToBytes(cacheKey), []byte("flag"))
		})
	} else {

		return t.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(t.TransactionBucket)
			encoded, err := encodeTransactions(transaction)
			if err != nil {
				return err
			}
			return b.Put(stringToBytes(cacheKey), encoded)
		})
	}
}

func (t *CashRepository) TransactionEmptyExists(cacheKey string) (bool, error) {
	var exists bool

	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.TransactionBucket)
		value := b.Get(stringToBytes(cacheKey))
		if value != nil {
			if bytes.Equal(value, []byte("flag")) {
				exists = true //The data is there, but it is empty
			} else {
				exists = false //Data is not empty
			}
		} else {
			exists = false //No data
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (t *CashRepository) UserInfoGet(cacheKey string) (repository.UserInfo, error) {
	var userInfo repository.UserInfo

	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.UserInfoBucket)
		value := b.Get(stringToBytes(cacheKey))
		if value != nil {
			var err error
			userInfo, err = decodeUserInfo(value)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return repository.UserInfo{}, err
	}

	return userInfo, nil
}

func (t *CashRepository) UserInfoSave(cacheKey string, userInfo repository.UserInfo) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.UserInfoBucket)
		encoded, err := encodeUserInfo(userInfo)
		if err != nil {
			return err
		}
		return b.Put(stringToBytes(cacheKey), encoded)
	})
}

func decodeTransactions(data []byte) ([]repository.Transaction, error) {
	var transactions []repository.Transaction
	err := json.Unmarshal(data, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func encodeTransactions(transaction []repository.Transaction) ([]byte, error) {
	data, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func decodeUserInfo(data []byte) (repository.UserInfo, error) {
	var userInfo repository.UserInfo
	err := json.Unmarshal(data, &userInfo)
	if err != nil {
		return repository.UserInfo{}, err
	}
	return userInfo, nil
}

func encodeUserInfo(userInfo repository.UserInfo) ([]byte, error) {
	data, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func stringToBytes(s string) []byte {
	return []byte(s)
}
