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
	template := &BillCreatedTemplate{sender.bill}
	menu := &telebot.ReplyMarkup{}
	menu.Inline(menu.Row(CancelBillBtn(sender.bill.ID)))
	return bot.Send(recipient, template.Render(), menu)
}

type DateBillsSender struct {
	Bills            []*models.Bill
	Year, Month, Day int
	ShowYear         bool
}

func (sender *DateBillsSender) Send(bot *telebot.Bot, recipient telebot.Recipient, _ *telebot.SendOptions) (*telebot.Message, error) {
	return bot.Send(recipient, sender.Text(), sender.ReplyMarkup())
}

func (sender *DateBillsSender) Text() string {
	titleTemplate := DateTitleTemplate{sender.Year, sender.Month, sender.Day, sender.ShowYear}
	billsTemplate := BillListTemplate{Bills: sender.Bills, MergeCategory: false}
	return titleTemplate.Render() + "\n\n" + billsTemplate.Render()
}

func (sender *DateBillsSender) ReplyMarkup() *telebot.ReplyMarkup {
	selector := &telebot.ReplyMarkup{ResizeKeyboard: true}
	selector.Inline(selector.Row(
		PrevDayBillBtn(sender.Year, sender.Month, sender.Day),
		NextDayBillBtn(sender.Year, sender.Month, sender.Day),
	))
	return selector
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
