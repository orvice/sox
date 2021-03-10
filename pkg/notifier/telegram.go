package notifier

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramNotifier struct {
	cli    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramNotifier{cli: bot, chatID: chatID}, nil
}

func (t *TelegramNotifier) Send(ctx context.Context,message string) error {
	msg := tgbotapi.NewMessage(t.chatID,message)
	_,err := t.cli.Send(msg)
	return err
}
