package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddOnHoldColumnCurrencyNotifTable() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		IsOnHold bool `gorm:"default:0"`
		IsMaster bool `gorm:"default:0"`
	}

	return &gormigrate.Migration{
		ID: "202106252134",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "IsOnHold")
			if err == nil {
				err = tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "IsMaster")
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "IsOnHold")
			if err == nil {
				err = tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "IsMaster")
			}

			return err
		},
	}
}
