package models

import (
	"time"
)

type CurrencyNotifConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Symbol    string `gorm:"size:200;unique"`
	IsOnHold  bool   `gorm:"default:0"`
	IsMaster  bool   `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
