package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddHoldedAtColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		HoldedAt int64 `gorm:"default:0"`
	}

	return &gormigrate.Migration{
		ID: "202109181354",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "HoldedAt")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "HoldedAt")

			return err
		},
	}
}
