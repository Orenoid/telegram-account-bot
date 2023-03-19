package user

import "github.com/orenoid/telegram-account-bot/models"

type Repository interface {
	CreateUser() (*models.User, error)
	CheckUserExists(userID uint) (bool, error)
	SetUserBalance(userID uint, balance float64) (float64, error)
	GetUserBalance(userID uint) (float64, error)
}
