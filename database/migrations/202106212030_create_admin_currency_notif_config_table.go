package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateCurrencyNotifConfigTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202106212030",
		Migrate: func(tx *gorm.DB) error {
			type CurrencyNotifConfig struct {
				ID        uint   `gorm:"primaryKey"`
				Symbol    string `gorm:"size:200;unique;not null"`
				CreatedAt time.Time
				UpdatedAt time.Time
			}
			return tx.AutoMigrate(&CurrencyNotifConfig{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("currency_notif_configs")
		},
	}
}
