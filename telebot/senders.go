package telebot

import (
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type NewBillSender struct {
	bill *models.Bill
}

func (sender *NewBillSender) Send(bot *telebot.Bot, recipient telebot.Recipient, _ *telebot.SendOptions) (*telebot.Message, error) {
	if sender.bill == nil {
		return nil, errors.New("bill should be nil")
	}
	// TODO 撤销功能
	template := &BillCreatedTemplate{sender.bill}
	button := telebot.InlineButton{Unique: "cancelBill", Text: "撤销"}
	opts := &telebot.SendOptions{
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{button}},
		},
	}
	return bot.Send(recipient, template.Render(), opts)
}
