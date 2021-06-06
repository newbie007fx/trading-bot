package repositories

import (
	"telebot-trading/app/models"
	"telebot-trading/external/db"
)

func GetConfigValueByKey(key string) *string {
	site_config := new(models.SiteConfig)
	res := db.GetDB().Where("key = ?", key).First(site_config)
	if res.Error != nil {
		return nil
	}
	return &site_config.Value
}

func SaveConfig(data map[string]interface{}) error {
	result := db.GetDB().Table("site_configs").Create(data)
	return result.Error
}
