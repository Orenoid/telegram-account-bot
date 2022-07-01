package telebot

import (
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillCreatedTemplate_Render(t *testing.T) {
	template := &BillCreatedTemplate{&models.Bill{Category: "购物", Amount: decimal.NewFromFloat(23)}}
	assert.Equal(t, "账单金额：23 元，类别：购物\n点击 \"撤销\" 可撤回该账单并回滚余额", template.Render())

	template = &BillCreatedTemplate{&models.Bill{Category: "购物", Amount: decimal.NewFromFloat(23.33)}}
	assert.Equal(t, "账单金额：23.33 元，类别：购物\n点击 \"撤销\" 可撤回该账单并回滚余额", template.Render())
}
