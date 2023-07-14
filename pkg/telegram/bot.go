package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mukolla/monobot/pkg/config"
	"github.com/mukolla/monobot/pkg/repository"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	tokenRepository repository.TokenRepository
	cash            repository.CashRepository
	redirectURL     string
	messages        config.Message
}

func NewBot(bot *tgbotapi.BotAPI, tr repository.TokenRepository, transactionRepository repository.CashRepository, messages config.Message) *Bot {
	return &Bot{bot: bot, tokenRepository: tr, cash: transactionRepository, messages: messages}
}

func (b *Bot) Start() error {
	updates, err := b.initUpdatesChanel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return err
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if _, err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) initUpdatesChanel() (tgbotapi.UpdatesChannel, error) {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60
	return b.bot.GetUpdatesChan(config)
}
