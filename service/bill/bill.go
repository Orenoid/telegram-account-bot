package bill

import (
	"github.com/orenoid/telegram-account-bot/dal/bill"
	"github.com/orenoid/telegram-account-bot/dal/user"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"time"
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

// CancelBillAndUpdateUserBalance 取消订单并更新用户余额
func (receiver *Service) CancelBillAndUpdateUserBalance(billID uint) error {
	return receiver.billRepo.DeleteBillAndUpdateUserBalance(billID)
}

type CreateBillDTO struct {
	Amount    float64
	Category  string
	Name      *string    // optional
	CreatedAt *time.Time // if not provided, then use current time as default
}

func (receiver *Service) CreateNewBills(userID uint, billDTOs []CreateBillDTO) error {
	userExists, err := receiver.userRepo.CheckUserExists(userID)
	if err != nil {
		return errors.WithStack(err)
	}
	if !userExists {
		return errors.New("user not exists")
	}
	createBillParams := make([]bill.CreateBillParams, 0, len(billDTOs))
	for _, billDTO := range billDTOs {
		createBillParams = append(createBillParams, bill.CreateBillParams{
			Amount:   billDTO.Amount,
			Category: billDTO.Category,
			CreateBillOptions: bill.CreateBillOptions{
				Name:      billDTO.Name,
				CreatedAt: billDTO.CreatedAt,
			},
		})
	}
	return receiver.billRepo.CreateBillsAndUpdateUserBalance(userID, createBillParams)
}

func NewService(billRepo bill.Repository, userRepo user.Repository) *Service {
	return &Service{billRepo: billRepo, userRepo: userRepo}
}
