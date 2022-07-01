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

func NewService(billRepo bill.Repository, userRepo user.Repository) *Service {
	return &Service{billRepo: billRepo, userRepo: userRepo}
}