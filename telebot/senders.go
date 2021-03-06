package telebot

import (
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

var _ telebot.Sendable = (*NewBillSender)(nil)
var _ telebot.Sendable = (*DateBillsSender)(nil)

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

type DateBillsSender struct {
	Bills            []*models.Bill
	Year, Month, Day int
	ShowYear         bool
}

func (sender *DateBillsSender) Send(bot *telebot.Bot, recipient telebot.Recipient, _ *telebot.SendOptions) (*telebot.Message, error) {
	titleTemplate := DateTitleTemplate{sender.Year, sender.Month, sender.Day, sender.ShowYear}
	billsTemplate := BillListTemplate{Bills: sender.Bills, MergeCategory: false}
	return bot.Send(recipient, titleTemplate.Render()+"\n\n"+billsTemplate.Render())
}

type MonthBillsSender struct {
	Bills       []*models.Bill
	Year, Month int
}

func (sender *MonthBillsSender) Send(bot *telebot.Bot, recipient telebot.Recipient, _ *telebot.SendOptions) (*telebot.Message, error) {
	titleTemplate := MonthTitleTemplate{Year: sender.Year, Month: sender.Month}
	billsTemplate := BillListTemplate{Bills: sender.Bills, MergeCategory: true}
	return bot.Send(recipient, titleTemplate.Render()+"\n\n"+billsTemplate.Render())
}
