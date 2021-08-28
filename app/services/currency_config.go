package services

import (
	"errors"
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"time"
)

func HoldCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
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

	if getMode() != "manual" {
		crypto.Buy(currencyConfig, candleData)
	}

	return nil
}

func ReleaseCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
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

	if getMode() != "manual" {
		crypto.Sell(currencyConfig, candleData)
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

	msg := crypto.GenerateMsg(*result)
	if config.IsOnHold {
		msg += holdCoinMessage(config, result)
	}

	return msg
}

func holdCoinMessage(config models.CurrencyNotifConfig, result *models.BandResult) string {
	var changes float32

	if config.HoldPrice < result.CurrentPrice {
		changes = (result.CurrentPrice - config.HoldPrice) / config.HoldPrice * 100
	} else {
		changes = (config.HoldPrice - result.CurrentPrice) / config.HoldPrice * 100
	}

	format := "Hold status: \nHold price: <b>%f</b> \nBalance: <b>%f</b> \nCurrent price: <b>%f</b> \nChanges: <b>%.2f%%</b> \nEstimation in USDT: <b>%f</b> \n"
	msg := fmt.Sprintf(format, config.HoldPrice, config.Balance, result.CurrentPrice, changes, (result.CurrentPrice * config.Balance))

	return msg
}

func getMode() string {
	mode := "manual"

	result := repositories.GetConfigValueByName("mode")
	if result != nil {
		mode = *result
	}

	return mode
}
