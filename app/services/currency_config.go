package services

import (
	"errors"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"time"
)

func HoldCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
	if getMode() != "manual" {
		crypto.Buy(currencyConfig, candleData)
	}

	if !currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": true,
		}
		err := repositories.UpdateCurrencyNotifConfig(currencyConfig.ID, data)
		if err != nil {
			log.Println(err.Error())
			return errors.New("error waktu update lur")
		}
	}

	return nil
}

func ReleaseCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
	if getMode() != "manual" {
		crypto.Sell(currencyConfig, candleData)
	}

	if currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": false,
		}
		err := repositories.UpdateCurrencyNotifConfig(currencyConfig.ID, data)
		if err != nil {
			log.Println(err.Error())
			return errors.New("error waktu update lur")
		}
	}

	return nil
}

func GetCurrencyStatus(config models.CurrencyNotifConfig) string {
	currentTime := time.Now()
	timeInMili := currentTime.Unix() * 1000

	responseChan := make(chan crypto.CandleResponse)
	request := crypto.CandleRequest{
		Symbol:       config.Symbol,
		EndDate:      timeInMili,
		Limit:        33,
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := crypto.MakeCryptoRequest(config, request)

	return crypto.GenerateMsg(*result)
}

func getMode() string {
	mode := "manual"

	result := repositories.GetConfigValueByName("mode")
	if result != nil {
		mode = *result
	}

	return mode
}
