package user

import "github.com/orenoid/account-bot/models"

type Repository interface {
	CreateUser() (*models.User, error)
	CheckUserExists(userID uint) (bool, error)
	// SetUserBalance 更新用户余额，若为空，则覆盖
	SetUserBalance(userID uint, balance float64) (float64, error)
}
