package telegram

import (
	"github.com/orenoid/telegram-account-bot/dal/telegram"
	"github.com/orenoid/telegram-account-bot/models"
)

type Service struct {
	teleRepo telegram.Repository
}

func (s *Service) CreateOrUpdateTelegramUser(userID int64, userName string, chatID int64) (*models.TelegramUser, error) {
	return s.teleRepo.CreateOrUpdateTelegramUser(userID, userName, chatID)
}

func NewService(teleRepo telegram.Repository) *Service {
	return &Service{
		teleRepo: teleRepo,
	}
}
