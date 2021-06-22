package repositories

import (
	"telebot-trading/app/models"
	"telebot-trading/external/db"
)

func GetCurrencyNotifConfigs() *[]models.CurrencyNotifConfig {
	notifConfigs := []models.CurrencyNotifConfig{}
	res := db.GetDB().Find(&notifConfigs)
	if res.Error != nil {
		return nil
	}
	return &notifConfigs
}
