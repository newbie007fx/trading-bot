package models

import (
	"time"
)

type CurrencyNotifConfig struct {
	ID        uint    `gorm:"primaryKey"`
	Symbol    string  `gorm:"size:200;unique"`
	IsOnHold  bool    `gorm:"default:0"`
	IsMaster  bool    `gorm:"default:0"`
	Volume    float32 `gorm:"default:0;index"`
	Balance   float32 `gorm:"default:0;index"`
	HoldPrice float32 `gorm:"default:0"`
	HoldedAt  int64   `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (cnc CurrencyNotifConfig) GetFormattedUpdatedAt() string {
	return cnc.UpdatedAt.Format("2006-01-02 15:04:05")
}
