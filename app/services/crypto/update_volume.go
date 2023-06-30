package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
)

func UpdateVolume() {
	log.Println("starting update volume worker ")

	symbolsLocal := getListSymbol()

	cryptoDriver := driver.GetCrypto()
	priceChanges, err := cryptoDriver.ListPriceChangeStats(&symbolsLocal)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, data := range *priceChanges {
		err := repositories.UpdateCurrencyNotifConfigBySymbol(data.Symbol, map[string]interface{}{
			"volume":        data.Volume,
			"price_changes": 0,
		})

		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}
	}

	log.Println("update volume worker done")
}

func getListActiveSymbol() []string {
	symbols := []string{}
	condition := map[string]interface{}{
		"status": models.STATUS_ACTIVE,
	}

	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil, nil, nil)
	for _, data := range *currency_configs {
		symbols = append(symbols, data.Symbol)
	}

	return symbols
}
