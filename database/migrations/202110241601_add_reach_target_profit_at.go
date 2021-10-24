package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddReachTargetProfitAtColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		ReachTargerProfitAt int64 `gorm:"default:0"`
	}

	return &gormigrate.Migration{
		ID: "202110241601",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "ReachTargerProfitAt")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "ReachTargerProfitAt")

			return err
		},
	}
}
