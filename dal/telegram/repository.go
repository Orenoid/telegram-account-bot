package telegram

import "github.com/orenoid/telegram-account-bot/models"

type Repository interface {
	CreateOrUpdateTelegramUser(userID int64, userName string, chatID int64) (*models.TelegramUser, error)
	GetUser(teleUserID int64) (*models.User, error)
}
