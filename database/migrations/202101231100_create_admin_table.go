package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateAdminTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202101231100",
		Migrate: func(tx *gorm.DB) error {
			type Admin struct {
				ID        uint   `gorm:"primaryKey"`
				Name      string `gorm:"size:200"`
				Email     string `gorm:"unique"`
				Password  string `gorm:"size:200"`
				Role      string `gorm:"size:200"`
				CreatedAt time.Time
				UpdatedAt time.Time
			}
			return tx.AutoMigrate(&Admin{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("admins")
		},
	}
}
