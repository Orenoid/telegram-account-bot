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
	bot.Handle("/help", hub.HandleHelpCommand)
	bot.Handle("/start", hub.HandleStartCommand)
	bot.Handle("/day", hub.HandleDayCommand)
	bot.Handle("/month", hub.HandleMonthCommand)
	bot.Handle("/cancel", hub.HandleCancelCommand)
	bot.Handle("/set_keyboard", hub.HandleSetKeyboardCommand)
	bot.Handle("/set_balance", hub.HandleSetBalanceCommand)
	bot.Handle(telebot.OnText, hub.HandleText)
	// 回调
	bot.Handle(&prevDayBillBtnTmpl, hub.HandleDayBillSelectionCallback)
	bot.Handle(&nextDayBillBtnTmpl, hub.HandleDayBillSelectionCallback)
	bot.Handle(&prevMonthBtnTmpl, hub.HandleMonthBillSelectionCallback)
	bot.Handle(&nextMonthBtnTmpl, hub.HandleMonthBillSelectionCallback)
	bot.Handle(&cancelBillBtnTmpl, hub.HandleCancelBillCallback)
}
