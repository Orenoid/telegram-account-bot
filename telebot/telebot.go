package telebot

import (
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

func NewBot(settings telebot.Settings, hub *HandlersHub) (*telebot.Bot, error) {
	bot, err := telebot.NewBot(settings)
	if err != nil {
		return &telebot.Bot{}, errors.WithStack(err)
	}

	RegisterHandlers(bot, hub)

	return bot, nil
}

func RegisterHandlers(bot *telebot.Bot, hub *HandlersHub) {
	bot.Handle("/start", hub.HandleStartCommand)
	bot.Handle(telebot.OnText, hub.HandleText)
}
