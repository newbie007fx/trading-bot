package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddStatusAndConfigColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		Status int8 `gorm:"default:10;index"`
		Config string
	}

	return &gormigrate.Migration{
		ID: "202212191501",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "Status")
			if err == nil {
				err = tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "Config")
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "Status")
			if err == nil {
				err = tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "Config")
			}

			return err
		},
	}
}
