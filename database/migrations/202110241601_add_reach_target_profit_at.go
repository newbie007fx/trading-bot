package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddReachTargetProfitAtColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		ReachTargetProfitAt int64 `gorm:"default:0"`
	}

	return &gormigrate.Migration{
		ID: "202110241601",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "ReachTargetProfitAt")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "ReachTargetProfitAt")

			return err
		},
	}
}
