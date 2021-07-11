package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddVolumeColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		Volume float32 `gorm:"default:0;index"`
	}

	return &gormigrate.Migration{
		ID: "202107111120",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "Volume")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "Volume")

			return err
		},
	}
}
