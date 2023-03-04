package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddCreatedAtAndUpdatedAtOnCurrentyConfig() *gormigrate.Migration {
	type CurrencyNotifConfig struct {
		CreatedAt int64
		UpdatedAt int64
	}

	return &gormigrate.Migration{
		ID: "202303041329",
		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "CreatedAt")
			if err == nil {
				err = tx.Migrator().AddColumn(&CurrencyNotifConfig{}, "UpdatedAt")
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "CreatedAt")
			if err == nil {
				err = tx.Migrator().DropColumn(&CurrencyNotifConfig{}, "UpdatedAt")
			}

			return err
		},
	}
}
