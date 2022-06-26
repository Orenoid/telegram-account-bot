package bill

import (
	"database/sql"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func (receiver *mysqlRepo) CreateBillAndUpdateUserBalance(userID uint, amount float64, category string, opts ...CreateBillOptions) (*models.Bill, error) {
	// TODO 解决并发更新余额问题，以及查询用户与更新余额的非原子操作场景
	var newBill *models.Bill
	err := receiver.db.Transaction(func(tx *gorm.DB) error {
		userModel := &models.User{}
		result := tx.Where("id = ?", userID).First(userModel)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}

		newBill = &models.Bill{UserID: userID, Amount: decimal.NewFromFloat(amount), Category: category}
		for _, opt := range opts {
			if opt.Name != nil {
				newBill.Name = sql.NullString{String: *opt.Name, Valid: true}
			}
		}

		result = tx.Create(newBill)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		if result.RowsAffected == 0 {
			return errors.New("failed to create bill")
		}

		if userModel.Balance.Valid {
			userModel.Balance.Decimal = userModel.Balance.Decimal.Add(newBill.Amount)
			result := tx.Save(userModel)
			if result.Error != nil {
				return errors.WithStack(result.Error)
			}
			if result.RowsAffected != 1 {
				return errors.New("failed to update user balance")
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return newBill, nil
}

func NewMysqlRepo(dsn string) (*mysqlRepo, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &mysqlRepo{db}, nil
}
