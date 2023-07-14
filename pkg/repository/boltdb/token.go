package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/mukolla/monobot/pkg/repository"
	"log"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(chatID int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))
	})
}

func (r *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}

func (r *TokenRepository) GetAll(bucket repository.Bucket) (map[string]string, error) {
	results := make(map[string]string)
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()

		for key, value := c.First(); key != nil; key, value = c.Next() {
			results[string(key)] = string(value)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no data found in the bucket")
	}

	// Виведення результатів у логи
	for key, value := range results {
		log.Printf("Key: %s, Value: %s\n", key, value)
	}

	return results, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
