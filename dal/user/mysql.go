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

func (receiver *mysqlRepo) MustGetToken(token string) (*models.Token, error) {
	var foundToken models.Token

	// 查找指定 token 的记录
	err := receiver.db.Where("token = ?", token).First(&foundToken).Error
	if err != nil {
		// 如果未找到记录，或者发生其他错误，返回 nil 和错误信息
		return nil, errors.WithStack(err)
	}

	// 返回找到的 token 记录
	return &foundToken, nil
}

func (receiver *mysqlRepo) CreateToken(userID uint, token string) error {
	newTokenRecord := &models.Token{
		UserID: userID,
		Token:  token,
	}
	result := receiver.db.Create(newTokenRecord)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}
	return nil
}

func (receiver *mysqlRepo) ValidateToken(token string) (bool, error) {
	var existingToken models.Token

	// 查找用户的 token
	err := receiver.db.Where("token = ?", token).First(&existingToken).Error
	if err != nil {
		// 如果没有找到 token，返回 false 和 nil 错误
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		// 其他数据库错误
		return false, errors.WithStack(err)
	}

	// 如果找到了 token，返回 true
	return true, nil
}

func (receiver *mysqlRepo) DisableToken(userID uint, token string) error {
	if err := receiver.db.Where(
		"user_id = ? AND token = ?", userID, token).Delete(&models.Token{}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (receiver *mysqlRepo) DisableAllTokens(userID uint) error {
	if err := receiver.db.Where("user_id = ?", userID).Delete(&models.Token{}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
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
