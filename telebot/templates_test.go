package telebot

import (
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestBillCreatedTemplate_Render(t *testing.T) {
	template := &BillCreatedTemplate{&models.Bill{Category: "购物", Amount: decimal.NewFromFloat(23)}}
	assert.Equal(t, "账单金额：23 元，类别：购物\n点击 \"撤销\" 可撤回该账单并回滚余额", template.Render())

	template = &BillCreatedTemplate{&models.Bill{Category: "购物", Amount: decimal.NewFromFloat(23.33)}}
	assert.Equal(t, "账单金额：23.33 元，类别：购物\n点击 \"撤销\" 可撤回该账单并回滚余额", template.Render())
}

func TestBillListTemplate_Render(t *testing.T) {
	// 测试合并统计支出
	bills := []*models.Bill{
		{Category: "饮食", Amount: decimal.NewFromFloat(-10)},
		{Category: "饮食", Amount: decimal.NewFromFloat(-20.22)},
		{Category: "出行", Amount: decimal.NewFromFloat(-10)},
	}

	template := BillListTemplate{Bills: bills, MergeCategory: true}
	assert.Equal(t, `合计支出：40.22 元

       - 饮食：30.22 元
       - 出行：10 元`, template.Render())

	// 测试合并统计收入
	bills = []*models.Bill{
		{Category: "工资", Amount: decimal.NewFromFloat(10000)},
		{Category: "咸鱼", Amount: decimal.NewFromFloat(100.01)},
		{Category: "咸鱼", Amount: decimal.NewFromFloat(50.01)},
	}
	template = BillListTemplate{Bills: bills, MergeCategory: true}
	assert.Equal(t, `合计收入：10150.02 元

       - 工资：10000 元
       - 咸鱼：150.02 元`, template.Render())

	// 测试同时有支出和收入的字符串拼接行为
	bills = []*models.Bill{
		{Category: "饮食", Amount: decimal.NewFromFloat(-10)},
		{Category: "饮食", Amount: decimal.NewFromFloat(-20.22)},
		{Category: "出行", Amount: decimal.NewFromFloat(-10)},
		{Category: "工资", Amount: decimal.NewFromFloat(10000)},
		{Category: "咸鱼", Amount: decimal.NewFromFloat(100.01)},
		{Category: "咸鱼", Amount: decimal.NewFromFloat(50.01)},
	}
	template = BillListTemplate{Bills: bills, MergeCategory: true}
	assert.Equal(t,
		`合计支出：40.22 元

       - 饮食：30.22 元
       - 出行：10 元

合计收入：10150.02 元

       - 工资：10000 元
       - 咸鱼：150.02 元`, template.Render())

	template = BillListTemplate{Bills: nil}
	assert.Equal(t, "暂无收支记录", template.Render())

	// 测试订单单独展示行为（按创建时间排序）
	now := time.Now()
	bills = []*models.Bill{
		{Category: "饮食", Amount: decimal.NewFromFloat(-10), Model: gorm.Model{CreatedAt: now}},
		{Category: "饮食", Amount: decimal.NewFromFloat(-20.22), Model: gorm.Model{CreatedAt: now.Add(2 * time.Second)}},
		{Category: "出行", Amount: decimal.NewFromFloat(-10), Model: gorm.Model{CreatedAt: now.Add(-2 * time.Second)}},

		{Category: "工资", Amount: decimal.NewFromFloat(10000), Model: gorm.Model{CreatedAt: now}},
		{Category: "咸鱼", Amount: decimal.NewFromFloat(100.01), Model: gorm.Model{CreatedAt: now.Add(-2 * time.Second)}},
	}
	template = BillListTemplate{Bills: bills, MergeCategory: false}
	assert.Equal(t,
		`合计支出：40.22 元

       - 出行：10 元
       - 饮食：10 元
       - 饮食：20.22 元

合计收入：10100.01 元

       - 咸鱼：100.01 元
       - 工资：10000 元`, template.Render())

	template = BillListTemplate{Bills: nil}
	assert.Equal(t, "暂无收支记录", template.Render())
}

func TestMonthTitleTemplate_Render(t *testing.T) {
	template := &MonthTitleTemplate{Year: 2022, Month: 1}
	assert.Equal(t, "2022年1月收支统计", template.Render())
}

func TestDateTitleTemplate_Render(t *testing.T) {
	template := &DateTitleTemplate{2022, 12, 25, true}
	assert.Equal(t, "2022年12月25日收支统计", template.Render())

	template = &DateTitleTemplate{2022, 12, 25, false}
	assert.Equal(t, "12月25日收支统计", template.Render())
}
