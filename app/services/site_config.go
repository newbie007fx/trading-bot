package services

import (
	"telebot-trading/app/helper"
	"telebot-trading/app/repositories"
)

func GetConfigValueByKey(key string) *string {
	simple_store := helper.GetSimpleStore()

	value := simple_store.Get(key)
	if value == nil {
		value = repositories.GetConfigValueByKey(key)
		if value != nil {
			SaveConfig(key, *value)
		}
	}

	return value
}

func SaveConfig(key, value string) error {
	data := map[string]interface{}{
		"key":   key,
		"value": value,
	}

	return repositories.SaveConfig(data)
}
