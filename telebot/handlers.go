package telebot

import (
	"github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type HandlersHub struct {
	teleService *telegram.Service
}

func NewHandlerHub(teleService *telegram.Service) *HandlersHub {
	return &HandlersHub{teleService: teleService}
}

func (hub *HandlersHub) HandleStartCommand(ctx telebot.Context) error {
	chat := ctx.Chat()
	if chat == nil {
		return errors.New("nil chat of context")
	}

	sender := ctx.Sender()
	_, err := hub.teleService.CreateOrUpdateTelegramUser(sender.ID, sender.Username, chat.ID)
	if err != nil {
		return err
	}
	// TODO set default keyboard
	err = ctx.Send("hello") // TODO send help message
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
