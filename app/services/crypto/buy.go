package crypto

import (
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
	"time"
)

func Buy(config models.CurrencyNotifConfig, candleData *models.CandleData) error {
	crypto := driver.GetCrypto()
	balance := GetBalanceFromConfig()

	if candleData == nil {
		currentTime := time.Now()
		timeInMili := (currentTime.Unix() - 1) * 1000

		candlesData, err := crypto.GetCandlesData(config.Symbol, 1, 0, timeInMili, "15m")
		if err != nil {
			return err
		}
		candleData = &candlesData[0]
	}

	condition := map[string]interface{}{"is_on_hold": true}
	holdCount := repositories.CountNotifConfig(&condition)

	maxHold := GetMaxHold()
	if maxHold-holdCount == 1 {
		balance -= 0.1
	}

	coinBalance := balance / (float32(maxHold) - float32(holdCount))

	if GetMode() == "automatic" {
		result, err := crypto.CreateBuyOrder(config.Symbol, coinBalance)
		if err != nil {
			log.Println(err.Error())
			return fmt.Errorf("error when try to buy coin %s with amount %.2f", config.Symbol, coinBalance)
		}

		fmt.Printf("coin buy, symbol %s, balance %f, price %f, status %s\n", result.Symbol, result.Quantity, result.Price, result.Status)

		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": config.Balance + result.Quantity, "hold_price": result.Price, "reach_target_profit_at": 0})
		SetBalance(balance - (result.Quantity * result.Price))
		RequestSyncBalance()
	} else {
		totalCoin := coinBalance / candleData.Close
		fmt.Printf("buy coin %s with amount %.2f, price %f, total coin %f\n", config.Symbol, coinBalance, candleData.Close, totalCoin)

		SetBalance(balance - (totalCoin * candleData.Close))
		repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": totalCoin, "hold_price": candleData.Close, "reach_target_profit_at": 0})
	}

	return nil
}
