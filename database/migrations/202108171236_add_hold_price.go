package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddHoldPriceColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		HoldPrice float32 `gorm:"default:0"`
	}

	return &gormigrate.Migration{
		ID: "202108171236",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "HoldPrice")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "HoldPrice")

			return err
		},
	}
}
