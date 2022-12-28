package repositories

import (
	"telebot-trading/app/models"
	"telebot-trading/external/db"
	"time"
)

func GetCurrencyNotifConfigsIgnoredCoins(condition *map[string]interface{}, limit *int, ignoredCoins *[]string, order *string) *[]models.CurrencyNotifConfig {
	notifConfigs := []models.CurrencyNotifConfig{}

	defaultOrder := "is_master desc, is_on_hold desc, volume desc, price_changes desc"
	if order == nil {
		order = &defaultOrder
	}

	res := db.GetDB().Order(*order)
	if limit != nil {
		res.Limit(*limit)
	}

	if ignoredCoins != nil {
		res.Not(map[string]interface{}{"symbol": ignoredCoins})
	}

	if condition != nil {
		for key, value := range *condition {
			res.Where(key, value)
		}
	}

	res.Find(&notifConfigs)
	if res.Error != nil {
		return nil
	}

	return &notifConfigs
}

func GetCurrencyNotifConfigs(condition *map[string]interface{}, limit *int, offset *int, order *string) *[]models.CurrencyNotifConfig {
	notifConfigs := []models.CurrencyNotifConfig{}

	defaultOrder := "is_master desc, is_on_hold desc, volume desc, price_changes desc"
	if order == nil {
		order = &defaultOrder
	}

	res := db.GetDB().Order(*order)
	if limit != nil {
		res.Limit(*limit)
	}

	if offset != nil {
		res.Offset(*offset)
	}

	if condition != nil {
		for key, value := range *condition {
			res.Where(key, value)
		}
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

func UpdateCurrencyNotifConfigBySymbol(symbol string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	result := db.GetDB().Table("currency_notif_configs").Where("symbol = ?", symbol).Updates(data)
	return result.Error
}

func UpdateCurrencyNotifConfigAll(data map[string]interface{}, condition *map[string]interface{}) error {
	data["updated_at"] = time.Now()
	query := db.GetDB().Table("currency_notif_configs")
	if condition != nil {
		for key, value := range *condition {
			query.Where(key, value)
		}
	}
	result := query.Updates(data)
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
