package models

import (
	"time"
)

type CurrencyNotifConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Symbol    string `gorm:"size:200;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
