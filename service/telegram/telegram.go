package telegram

import (
	"github.com/orenoid/telegram-account-bot/dal/telegram"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/pkg/errors"
)

type Service struct {
	teleRepo telegram.Repository
}

func (s *Service) CreateOrUpdateTelegramUser(userID int64, userName string, chatID int64) (*models.TelegramUser, error) {
	return s.teleRepo.CreateOrUpdateTelegramUser(userID, userName, chatID)
}

func (s *Service) GetBaseUserID(teleUserID int64) (uint, error) {
	baseUser, err := s.teleRepo.GetUser(teleUserID)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return baseUser.ID, nil
}

func NewService(teleRepo telegram.Repository) *Service {
	return &Service{
		teleRepo: teleRepo,
	}
}
