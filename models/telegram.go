package models

import "gorm.io/gorm"

type TelegramUser struct {
	gorm.Model
	BaseUserID uint   `gorm:"not null;unique"`
	UserName   string `gorm:"not null;unique"`
	ChatID     int64  `gorm:"not null"`
}
