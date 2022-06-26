package bill

import "github.com/orenoid/telegram-account-bot/models"

type Repository interface {
	// CreateBillAndUpdateUserBalance 为用户创建一个账单，并更新用户余额（若用户余额不为空）
	CreateBillAndUpdateUserBalance(userID uint, amount float64, category string, opts ...CreateBillOptions) (*models.Bill, error)
}

type CreateBillOptions struct {
	Name *string
}
