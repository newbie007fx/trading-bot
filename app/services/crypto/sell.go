package crypto

import (
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
		result, err := crypto.CreateSellOrder(config.Symbol, config.Balance)
		if err != nil {
			return err
		}
		SetBalance(balance + (result.Price * result.Quantity))
	} else {
		SetBalance(balance + totalBalance)
		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": 0})
	}

	return nil
}
