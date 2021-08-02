package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
	"time"
)

func HoldCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
	if getMode() != "manual" {
		buy(currencyConfig, candleData)
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
		sell(currencyConfig, candleData)
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

func buy(config models.CurrencyNotifConfig, candleData *models.CandleData) error {
	balance := getBalance()
	if candleData == nil {
		endTime := time.Now().Unix()
		startTime := endTime - (60 * 15)
		crypto := driver.GetCrypto()
		candlesData, err := crypto.GetCandlesData(config.Symbol, startTime, endTime, "15")
		if err != nil {
			return err
		}
		candleData = &candlesData[len(candlesData)-1]
	}

	totalCoin := balance / candleData.Close
	repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": totalCoin})
	SetBalance(balance - (totalCoin * candleData.Close))

	return nil
}

func sell(config models.CurrencyNotifConfig, candleData *models.CandleData) error {
	balance := getBalance()
	if candleData == nil {
		endTime := time.Now().Unix()
		startTime := endTime - (60 * 15)
		crypto := driver.GetCrypto()
		candlesData, err := crypto.GetCandlesData(config.Symbol, startTime, endTime, "15")
		if err != nil {
			return err
		}
		candleData = &candlesData[len(candlesData)-1]
	}

	totalBalance := config.Balance * candleData.Close
	repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": 0})
	SetBalance(balance + totalBalance)

	return nil
}

func getBalance() float32 {
	var balance float32 = 0

	result := repositories.GetConfigValueByName("balance")
	if result != nil {
		resultFloat, err := strconv.ParseFloat(*result, 32)
		if err == nil {
			balance = float32(resultFloat)
		}
	}

	return balance
}

func SetBalance(balance float32) error {
	s := fmt.Sprintf("%f", balance)
	return repositories.SetConfigByName("balance", s)
}

func getMode() string {
	mode := "manual"

	result := repositories.GetConfigValueByName("mode")
	if result != nil {
		mode = *result
	}

	return mode
}
