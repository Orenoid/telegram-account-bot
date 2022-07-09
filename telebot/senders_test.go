package telebot

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
	"reflect"
	"testing"
)

func TestDateBillsSender_Send(t *testing.T) {
	bills := []*models.Bill{{Amount: decimal.NewFromFloat(-1), UserID: 1, Category: "杂项"}}
	sender := &DateBillsSender{bills, 2022, 12, 1, false}

	type Send struct {
		called        bool
		paramTo       telebot.Recipient
		paramWhat     interface{}
		paramOpts     []interface{}
		returnMessage *telebot.Message
		returnErr     error
	}
	send := &Send{returnMessage: &telebot.Message{Text: "some special text"}, returnErr: errors.New("send error")}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&telebot.Bot{}), "Send",
		func(bot *telebot.Bot, to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
			send.called = true
			send.paramTo = to
			send.paramWhat = what
			send.paramOpts = opts
			return send.returnMessage, send.returnErr
		})
	defer patches.Reset()

	bot := &telebot.Bot{Me: &telebot.User{ID: 314}}
	recipient := &telebot.User{ID: 618}
	msg, err := sender.Send(bot, recipient, nil)
	titleTemplate := DateTitleTemplate{sender.Year, sender.Month, sender.Day, sender.ShowYear}
	billsTemplate := BillListTemplate{Bills: sender.Bills, MergeCategory: false}
	assert.Equal(t, titleTemplate.Render()+"\n\n"+billsTemplate.Render(), send.paramWhat)
	assert.Equal(t, recipient, send.paramTo)
	assert.Equal(t, send.returnErr, err)
	assert.Equal(t, send.returnMessage, msg)
}
