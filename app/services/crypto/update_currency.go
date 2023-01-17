package crypto

import (
	"encoding/json"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
)

func StartUpdateCurrencyService(updateCurrencyChan chan bool) {
	log.Println("update currency service is up")
	for <-updateCurrencyChan {
		log.Println("starting update currency worker ")
		UpdateCurrency()
		log.Println("update currency worker done")
		UpdateVolume()
	}
}

func UpdateCurrency() {
	condition := map[string]interface{}{
		"status": models.STATUS_ACTIVE,
	}
	repositories.UpdateCurrencyNotifConfigAll(map[string]interface{}{
		"status": models.STATUS_MARKET_OFF,
	}, &condition)

	cryptoDriver := driver.GetCrypto()
	symbols, err := cryptoDriver.GetExchangeInformation()
	if err == nil {
		for _, symbol := range *symbols {
			if symbol.QuoteAsset != "USDT" {
				continue
			}

			config := ""
			configByte, err := json.Marshal(symbol)
			if err == nil {
				config = string(configByte)
			}

			status := models.STATUS_ACTIVE
			if symbol.Status != "TRADING" {
				status = models.STATUS_MARKET_OFF
			}

			currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol.Symbol)
			if err == nil {
				repositories.UpdateCurrencyNotifConfig(currencyConfig.ID, map[string]interface{}{
					"status": status,
					"config": config,
				})
			} else if symbol.Status == "TRADING" {
				repositories.SaveCurrencyNotifConfig(map[string]interface{}{
					"symbol": symbol.Symbol,
					"status": status,
					"config": config,
				})
			}
		}
	} else {
		log.Println(err)
	}
}
