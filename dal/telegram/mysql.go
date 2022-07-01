package telegram

import (
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func (repo *mysqlRepo) CreateOrUpdateTelegramUser(userID int64, userName string, chatID int64) (*models.TelegramUser, error) {
	telegramUser := &models.TelegramUser{}
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		result := tx.Model(&models.TelegramUser{}).Where("id = ? and chat_id = ?", userID, chatID).Count(&count)
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
		if count > 0 {
			result := tx.Model(&models.TelegramUser{}).Where("id = ? and chat_id = ?", userID, chatID).
				Updates(map[string]interface{}{"user_name": userName})
			if result.Error != nil {
				return errors.WithStack(result.Error)
			}
			if result.RowsAffected == 0 {
				return errors.New("failed to update user")
			}
		} else {
			newBaseUser := &models.User{}
			result := tx.Create(newBaseUser)
			if result.Error != nil {
				return errors.WithStack(result.Error)
			}
			newTelegramUser := &models.TelegramUser{BaseUserID: newBaseUser.ID, UserName: userName, ChatID: chatID}
			newTelegramUser.ID = uint(userID)
			result = tx.Create(newTelegramUser)
			if result.Error != nil {
				return errors.WithStack(result.Error)
			}
			if result.RowsAffected == 0 {
				return errors.New("failed to create user")
			}
		}
		result = tx.First(telegramUser, userID)
		return result.Error
	})
	if err != nil {
		return nil, err
	}
	return telegramUser, nil
}

func (repo *mysqlRepo) GetUser(teleUserID int64) (*models.User, error) {
	var user models.User
	result := repo.db.Model(&user).Select("users.id, users.created_at, users.updated_at, users.balance").
		Joins("join telegram_users tu on tu.base_user_id = users.id").
		Where("tu.id = ?", teleUserID).Scan(&user)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	return &user, nil
}

func NewMysqlRepo(dsn string) (*mysqlRepo, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &mysqlRepo{db}, nil
}
