package telebot

import (
	"fmt"
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

type BillCreatedTemplate struct {
	bill *models.Bill
}

func (template *BillCreatedTemplate) Render() string {
	format :=
		`账单金额：%s 元，类别：%s
点击 "撤销" 可撤回该账单并回滚余额`
	return fmt.Sprintf(format, template.bill.Amount.String(), template.bill.Category)
}
