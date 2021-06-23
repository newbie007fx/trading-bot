package services

import (
	"telebot-trading/app/helper"
	"telebot-trading/app/repositories"
)

func GetConfigValueByName(name string) *string {
	simple_store := helper.GetSimpleStore()

	value := simple_store.Get(name)
	if value == nil {
		value = repositories.GetConfigValueByName(name)
		if value != nil {
			SaveConfig(name, *value)
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
