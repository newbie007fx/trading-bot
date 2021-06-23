package models

type SiteConfig struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Value string `gorm:"size:200"`
}
