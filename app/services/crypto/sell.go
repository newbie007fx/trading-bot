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
	crypto := driver.GetCrypto()
	balance := GetBalance()

	if candleData == nil {
		currentTime := time.Now()
		timeInMili := currentTime.Unix() * 1000

		candlesData, err := crypto.GetCandlesData(config.Symbol, 1, timeInMili, "15m")
		if err != nil {
			return err
		}
		candleData = &candlesData[0]
	}

	totalBalance := config.Balance * candleData.Close
	if GetMode() == "automatic" {
		result, err := crypto.CreateSellOrder(config.Symbol, totalBalance)
		if err != nil {
			return err
		}

		log.Println(fmt.Sprintf("coin sell, symbol %s, balance %f, price %f", result.Symbol, result.Quantity, result.Price))

		SetBalance(balance + (result.Price * result.Quantity))
		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": config.Balance - result.Quantity})
	} else {
		SetBalance(balance + totalBalance)
		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": 0})
	}

	return nil
}
