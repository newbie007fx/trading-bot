package crypto

import (
	"encoding/json"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
	"time"
)

func StartUpdateCurrencyService(updateCurrencyChan chan bool) {
	log.Println("update currency service is up")
	for <-updateCurrencyChan {
		currentTime := time.Now()
		time.Sleep(10 * time.Second)
		if currentTime.Day()%5 == 1 {
			log.Println("starting check update worker")
			UpdateVolume()
		} else {
			log.Println("starting check currency worker ")
			CheckCurrency()
		}
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
	symbols, err := cryptoDriver.GetExchangeInformation(nil)
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
					"status":        status,
					"config":        config,
					"price_changes": 0,
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

func CheckCurrency() {
	condition := map[string]interface{}{
		"status": models.STATUS_ACTIVE,
	}
	repositories.UpdateCurrencyNotifConfigAll(map[string]interface{}{
		"status": models.STATUS_MARKET_OFF,
	}, &condition)

	symbolsLocal := getListSymbol()

	cryptoDriver := driver.GetCrypto()
	symbols, err := cryptoDriver.GetExchangeInformation(&symbolsLocal)
	if err == nil {
		availableSymbols := []string{}
		for _, symbol := range *symbols {
			config := ""
			configByte, err := json.Marshal(symbol)
			if err == nil {
				config = string(configByte)
			}

			status := models.STATUS_ACTIVE
			if symbol.Status != "TRADING" {
				status = models.STATUS_MARKET_OFF
			}

			repositories.UpdateCurrencyNotifConfigBySymbol(symbol.Symbol, map[string]interface{}{
				"status":        status,
				"config":        config,
				"price_changes": 0,
			})

			availableSymbols = append(availableSymbols, symbol.Symbol)
		}

		err := repositories.DeleteCurrencyNotifConfigSymbolNotIn(availableSymbols)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}

func getListSymbol() []string {
	symbols := []string{}
	currency_configs := repositories.GetCurrencyNotifConfigs(nil, nil, nil, nil)
	for _, data := range *currency_configs {
		symbols = append(symbols, data.Symbol)
	}

	return symbols
}
