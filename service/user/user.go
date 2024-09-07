package user

import (
	"github.com/orenoid/telegram-account-bot/dal/user"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/orenoid/telegram-account-bot/utils/strings"
	"github.com/pkg/errors"
)

type Service struct {
	userRepo user.Repository
}

func (receiver *Service) CreateUser() (*models.User, error) {
	return receiver.userRepo.CreateUser()
}

func (receiver *Service) SetUserBalance(userID uint, balance float64) (float64, error) {
	userExists, err := receiver.userRepo.CheckUserExists(userID)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, errors.New("user not found")
	}
	newBalance, err := receiver.userRepo.SetUserBalance(userID, balance)
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (receiver *Service) GetUserBalance(userID uint) (float64, error) {
	userExists, err := receiver.userRepo.CheckUserExists(userID)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, errors.New("user not found")
	}
	balance, err := receiver.userRepo.GetUserBalance(userID)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (receiver *Service) CreateToken(userID uint) (string, error) {
	token, err := strings.GenerateToken()
	if err != nil {
		return "", err
	}
	err = receiver.userRepo.CreateToken(userID, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (receiver *Service) DisableAllTokens(userID uint) error {
	return receiver.userRepo.DisableAllTokens(userID)
}

func (receiver *Service) MustGetUserIDByToken(token string) (uint, error) {
	tokenRecord, err := receiver.userRepo.MustGetToken(token)
	if err != nil {
		return 0, err
	}
	return tokenRecord.UserID, nil
}

func NewUserService(userRepo user.Repository) *Service {
	return &Service{userRepo: userRepo}
}
