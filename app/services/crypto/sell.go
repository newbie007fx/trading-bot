package crypto

import (
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
	"time"
)

func Sell(config models.CurrencyNotifConfig, candleData *models.CandleData) error {
	cryptoDriver := driver.GetCrypto()
	balance := GetBalanceFromConfig()

	if candleData == nil {
		currentTime := time.Now()
		timeInMili := (currentTime.Unix() - 1) * 1000

		candlesData, err := cryptoDriver.GetCandlesData(config.Symbol, 1, 0, timeInMili, "15m")
		if err != nil {
			return err
		}
		candleData = &candlesData[0]
	}

	totalBalance := config.Balance * candleData.Close
	if GetMode() == "automatic" {
		result, err := sellWithRetry(config.Symbol, totalBalance)
		if err != nil {
			return fmt.Errorf("error when try to sell coin %s with amount %.2f", config.Symbol, config.Balance)
		}

		log.Printf("coin sell, symbol %s, balance %f, price %f, status %s\n", result.Symbol, result.Quantity, result.Price, result.Status)

		SetBalance(balance + (result.Price * result.Quantity))
		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": config.Balance - result.Quantity, "price_changes": 0})
		RequestSyncBalance()
	} else {
		fmt.Printf("sell coin %s with amount %.2f, price %f, total balance %f\n", config.Symbol, config.Balance, candleData.Close, totalBalance)
		SetBalance(balance + totalBalance)
		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": 0, "price_changes": 0})
	}

	return nil
}

func sellWithRetry(symbol string, baseBalance float32) (result *models.CreateOrderResponse, err error) {
	cryptoDriver := driver.GetCrypto()
	totalBalance := baseBalance
	for i := 0; i < 3; i++ {
		totalBalance -= 0.1
		result, err = cryptoDriver.CreateSellOrder(symbol, (totalBalance))
		if err != nil {
			log.Println("error when try to sell coin, msg: ", err.Error())
		} else {
			break
		}
	}

	return result, err
}
