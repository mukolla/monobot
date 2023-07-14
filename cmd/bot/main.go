package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mukolla/monobot/pkg/config"
	"github.com/mukolla/monobot/pkg/repository"
	"github.com/mukolla/monobot/pkg/repository/boltdb"
	"github.com/mukolla/monobot/pkg/telegram"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	db, err := initDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)
	transactionRepository := boltdb.NewCashRepository(db)

	telegramBot := telegram.NewBot(bot, tokenRepository, transactionRepository, cfg.Message)

	telegramBot.Start()
}

func initDb(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.TransactionCash))
		if err != nil {
			return err
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.UserInfoCash))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return db, err
}
