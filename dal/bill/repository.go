package bill

import (
	"github.com/orenoid/telegram-account-bot/models"
	"time"
)

type Repository interface {
	// CreateBillAndUpdateUserBalance 为用户创建一个账单，并更新用户余额（若用户余额不为空）
	CreateBillAndUpdateUserBalance(userID uint, amount float64, category string, opts ...CreateBillOptions) (*models.Bill, error)
	// GetUserBillsByCreateTime 获取用户在指定时间范围内的账单列表，若 opts 为空，则返回账单（opts 只取列表第一个作为查询参数）
	GetUserBillsByCreateTime(userID uint, opts ...GetUserBillsByCreateTimeOptions) ([]*models.Bill, error)
}

type CreateBillOptions struct {
	Name *string
}

type GetUserBillsByCreateTimeOptions struct {
	GreaterThan    time.Time // 时间范围区间左侧
	GreaterOrEqual bool      // 是否为闭区间
	LessThan       time.Time // 时间范围区间右侧
	LessOrEqual    bool      // 是否为闭区间
}
