package bill

import (
	"github.com/orenoid/telegram-account-bot/dal/bill"
	"github.com/orenoid/telegram-account-bot/dal/user"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
)

type Service struct {
	billRepo bill.Repository
	userRepo user.Repository
}

func (receiver *Service) CreateNewBill(userID uint, amount float64, category string, opts ...bill.CreateBillOptions) (*models.Bill, error) {
	userExists, err := receiver.userRepo.CheckUserExists(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !userExists {
		return nil, errors.New("user not exists")
	}
	return receiver.billRepo.CreateBillAndUpdateUserBalance(userID, amount, category, opts...)
}

// GetUserBillsByCreateTime 获取用户在指定时间范围内的账单列表，若 opts 为空，则返回账单（opts 只取列表第一个作为查询参数）
func (receiver *Service) GetUserBillsByCreateTime(userID uint, opts ...bill.GetUserBillsByCreateTimeOptions) ([]*models.Bill, error) {
	return receiver.billRepo.GetUserBillsByCreateTime(userID, opts...)
}

func NewService(billRepo bill.Repository, userRepo user.Repository) *Service {
	return &Service{billRepo: billRepo, userRepo: userRepo}
}
