package user

import (
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func (receiver *mysqlRepo) GetUserBalance(userID uint) (float64, error) {
	userModel := &models.User{}
	result := receiver.db.First(userModel, userID)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}
	return userModel.Balance.Decimal.InexactFloat64(), nil
}

func (receiver *mysqlRepo) CreateUser() (*models.User, error) {
	userModel := &models.User{}
	result := receiver.db.Create(userModel)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return userModel, nil
}

// CheckUserExists check if user exists
func (receiver *mysqlRepo) CheckUserExists(userID uint) (bool, error) {
	var count int64
	result := receiver.db.Model(&models.User{}).Where("id = ?", userID).Count(&count)
	if result.Error != nil {
		return false, errors.WithStack(result.Error)
	}
	return count > 0, nil
}

func (receiver *mysqlRepo) SetUserBalance(userID uint, balanceFloat float64) (float64, error) {
	balanceDecimal := decimal.NewNullDecimal(decimal.NewFromFloat(balanceFloat))

	result := receiver.db.Model(&models.User{}).
		Where("id = ?", userID).Update("balance", balanceDecimal)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.New("failed to set user balance")
	}

	userModel := &models.User{}
	result = receiver.db.First(userModel, userID)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}
	return userModel.Balance.Decimal.InexactFloat64(), nil
}

func NewMysqlRepo(dsn string) (*mysqlRepo, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &mysqlRepo{db}, nil
}
