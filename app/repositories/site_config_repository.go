package repositories

import (
	"telebot-trading/app/models"
	"telebot-trading/external/db"

	"gorm.io/gorm/clause"
)

func GetConfigValueByName(name string) *string {
	site_config := new(models.SiteConfig)
	res := db.GetDB().Where("name = ?", name).First(site_config)
	if res.Error != nil {
		return nil
	}
	return &site_config.Value
}

func SaveConfig(data map[string]interface{}) error {
	result := db.GetDB().Table("site_configs").Create(data)
	return result.Error
}

func SetConfigByName(name, value string) error {
	data := map[string]interface{}{
		"name":  name,
		"value": value,
	}
	result := db.GetDB().Table("site_configs").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(data)

	return result.Error
}
