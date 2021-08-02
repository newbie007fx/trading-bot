package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddBalanceColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		Balance float32 `gorm:"default:0;index"`
	}

	return &gormigrate.Migration{
		ID: "202108021329",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "Balance")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "Balance")

			return err
		},
	}
}
