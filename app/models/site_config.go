package models

type SiteConfig struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"unique"`
	Value string `gorm:"size:200"`
}
