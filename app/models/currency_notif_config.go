package models

import (
	"time"
)

const STATUS_ACTIVE = 10
const STATUS_MARKET_OFF = 5
const STATUS_NONACTIVE = 0

type CurrencyNotifConfig struct {
	ID                  uint    `gorm:"primaryKey"`
	Symbol              string  `gorm:"size:200;unique"`
	IsOnHold            bool    `gorm:"default:0"`
	IsMaster            bool    `gorm:"default:0"`
	Volume              float32 `gorm:"default:0;index"`
	Balance             float32 `gorm:"default:0;index"`
	HoldPrice           float32 `gorm:"default:0"`
	HoldedAt            int64   `gorm:"default:0"`
	ReachTargetProfitAt int64   `gorm:"default:0"`
	PriceChanges        float32 `gorm:"default:0"`
	Status              int8    `gorm:"default:0;index"`
	Config              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (cnc CurrencyNotifConfig) GetFormattedUpdatedAt() string {
	return cnc.UpdatedAt.Format("2006-01-02 15:04:05")
}
