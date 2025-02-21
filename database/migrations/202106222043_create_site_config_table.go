package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateSiteConfigTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202106222043",
		Migrate: func(tx *gorm.DB) error {
			type SiteConfig struct {
				ID    uint   `gorm:"primaryKey"`
				Name  string `gorm:"unique"`
				Value string `gorm:"size:200"`
			}
			return tx.AutoMigrate(&SiteConfig{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("site_configs")
		},
	}
}
