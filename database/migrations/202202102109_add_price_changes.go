package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddPriceChangeColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		PriceChanges float32 `gorm:"default:0;index"`
	}

	return &gormigrate.Migration{
		ID: "202202102109",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "PriceChanges")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "PriceChanges")

			return err
		},
	}
}
