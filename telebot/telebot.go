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
	// 基础命令
	bot.Handle("/start", hub.HandleStartCommand)
	bot.Handle("/day", hub.HandleDayCommand)
	bot.Handle("/month", hub.HandleMonthCommand)
	bot.Handle(telebot.OnText, hub.HandleText)
	// 回调
	bot.Handle(&prevDayBillBtnTmpl, hub.HandleDayBillSelectionCallback)
	bot.Handle(&nextDayBillBtnTmpl, hub.HandleDayBillSelectionCallback)
	bot.Handle(&cancelBillBtnTmpl, hub.HandleCancelBillCallback)
}
