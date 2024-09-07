package user

import "github.com/orenoid/telegram-account-bot/models"

type Repository interface {
	// user

	CreateUser() (*models.User, error)
	CheckUserExists(userID uint) (bool, error)
	SetUserBalance(userID uint, balance float64) (float64, error)
	GetUserBalance(userID uint) (float64, error)

	// auth

	CreateToken(userID uint, token string) error
	MustGetToken(token string) (*models.Token, error)
	DisableToken(userID uint, token string) error
	DisableAllTokens(userID uint) error
}
