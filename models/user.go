package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Balance decimal.NullDecimal `gorm:"type:decimal(10,2)"`
}
