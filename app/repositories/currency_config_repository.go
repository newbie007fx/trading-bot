package repositories

import (
	"telebot-trading/app/models"
	"telebot-trading/external/db"
	"time"
)

func GetCurrencyNotifConfigs(condition *map[string]interface{}, limit *int) *[]models.CurrencyNotifConfig {
	notifConfigs := []models.CurrencyNotifConfig{}
	res := db.GetDB().Order("is_master desc, is_on_hold desc, volume desc, id asc")
	if limit != nil {
		res.Limit(*limit)
	}
	if condition != nil {
		res.Where(*condition)
	}
	res.Find(&notifConfigs)
	if res.Error != nil {
		return nil
	}
	return &notifConfigs
}

func CountNotifConfig(condition *map[string]interface{}) (count int64) {
	query := db.GetDB().Table("currency_notif_configs")
	if condition != nil {
		query.Where(*condition)
	}
	query.Count(&count)
	return
}

func SaveCurrencyNotifConfig(data map[string]interface{}) error {
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()

	result := db.GetDB().Table("currency_notif_configs").Create(data)
	return result.Error
}

func GetCurrencyNotifConfig(id uint) (*models.CurrencyNotifConfig, error) {
	currencyConfig := models.CurrencyNotifConfig{}
	result := db.GetDB().Table("currency_notif_configs").Where("id = ?", id).Take(&currencyConfig)
	return &currencyConfig, result.Error
}

func GetMasterCoinConfig() (*models.CurrencyNotifConfig, error) {
	currencyConfig := models.CurrencyNotifConfig{}
	result := db.GetDB().Table("currency_notif_configs").Where("is_master = ?", true).Take(&currencyConfig)
	return &currencyConfig, result.Error
}

func GetCurrencyNotifConfigBySymbol(symbol string) (*models.CurrencyNotifConfig, error) {
	currencyConfig := models.CurrencyNotifConfig{}
	result := db.GetDB().Table("currency_notif_configs").Where("symbol = ?", symbol).Take(&currencyConfig)
	return &currencyConfig, result.Error
}

func UpdateCurrencyNotifConfig(id uint, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	result := db.GetDB().Table("currency_notif_configs").Where("id = ?", id).Updates(data)
	return result.Error
}

func DeleteCurrencyNotifConfig(id uint) error {
	currencyConfig, err := GetCurrencyNotifConfig(id)
	if err != nil {
		return err
	}

	result := db.GetDB().Delete(currencyConfig)
	return result.Error
}

func SetMaster(id uint) error {
	db.GetDB().Table("currency_notif_configs").Where("is_master = ?", 1).Updates(map[string]interface{}{"updated_at": time.Now(), "is_master": false})
	result := db.GetDB().Table("currency_notif_configs").Where("id = ?", id).Updates(map[string]interface{}{"updated_at": time.Now(), "is_master": true})

	return result.Error
}
