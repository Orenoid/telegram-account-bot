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

func (receiver *mysqlRepo) CreateBillsAndUpdateUserBalance(userID uint, billParams []CreateBillParams) error {
	err := receiver.db.Transaction(func(tx *gorm.DB) error {
		userModel := &models.User{}
		result := tx.Where("id = ?", userID).First(userModel)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}

		billRecords := make([]models.Bill, 0, len(billParams))
		for _, billParam := range billParams {
			billRecord := &models.Bill{
				UserID:   userID,
				Amount:   decimal.NewFromFloat(billParam.Amount),
				Category: billParam.Category,
			}
			if billParam.Name != nil {
				billRecord.Name = sql.NullString{String: *billParam.Name, Valid: true}
			}
			if billParam.CreatedAt != nil {
				billRecord.CreatedAt = *billParam.CreatedAt
				billRecord.UpdatedAt = *billParam.CreatedAt
			}
			billRecords = append(billRecords, *billRecord)
		}

		result = tx.Create(billRecords)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		if result.RowsAffected != int64(len(billRecords)) {
			return errors.New("failed to create bill")
		}

		if userModel.Balance.Valid {
			for _, billRecord := range billRecords {
				userModel.Balance.Decimal = userModel.Balance.Decimal.Add(billRecord.Amount)
			}
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
		return errors.WithStack(err)
	}
	return nil
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

func (receiver *mysqlRepo) GetUserBillsByCreateTime(userID uint, opts ...GetUserBillsByCreateTimeOptions) ([]*models.Bill, error) {
	var bills []*models.Bill
	query := receiver.db.Where("user_id = ?", userID)
	if len(opts) > 0 {
		opt := opts[0]
		if opt.GreaterOrEqual {
			query = query.Where("created_at >= ?", opt.GreaterThan)
		} else {
			query = query.Where("created_at = ?", opt.GreaterThan)
		}
		if opt.LessOrEqual {
			query = query.Where("created_at <= ?", opt.LessThan)
		} else {
			query = query.Where("created_at < ?", opt.LessThan)
		}
	}
	result := query.Find(&bills)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return bills, nil
}

func (receiver *mysqlRepo) DeleteBillAndUpdateUserBalance(billID uint) error {
	err := receiver.db.Transaction(func(tx *gorm.DB) error {
		// 删除订单
		billModel := &models.Bill{}
		result := tx.Where("id = ?", billID).First(billModel)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		result = tx.Delete(billModel)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		// 更新用户余额
		userModel := &models.User{}
		result = tx.Where("id = ?", billModel.UserID).First(userModel)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		if userModel.Balance.Valid {
			userModel.Balance.Decimal = userModel.Balance.Decimal.Sub(billModel.Amount)
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
	return err
}

func NewMysqlRepo(dsn string) (*mysqlRepo, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &mysqlRepo{db}, nil
}
