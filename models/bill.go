package models

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	UserID   uint            `gorm:"not null"`
	Amount   decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Category string          `gorm:"not null"`
	Name     sql.NullString
}
