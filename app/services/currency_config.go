package services

import (
	"errors"
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

func HoldCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
	if crypto.GetMode() != "manual" {
		err := crypto.Buy(currencyConfig, candleData)
		if err != nil {
			log.Println(err.Error())
			return err
		}
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
	if crypto.GetMode() != "manual" {
		err := crypto.Sell(currencyConfig, candleData)
		if err != nil {
			log.Println(err.Error())
			return err
		}
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
		Limit:        40,
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

func GetWeightLog(config models.CurrencyNotifConfig, datetime time.Time) string {
	timeInMili := datetime.Unix() * 1000

	responseChan := make(chan crypto.CandleResponse)
	request := crypto.CandleRequest{
		Symbol:       config.Symbol,
		EndDate:      timeInMili,
		Limit:        40,
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := crypto.MakeCryptoRequest(config, request)

	masterCoinConfig, _ := repositories.GetMasterCoinConfig()
	request = crypto.CandleRequest{
		Symbol:       masterCoinConfig.Symbol,
		EndDate:      timeInMili,
		Limit:        40,
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	masterCoin := crypto.MakeCryptoRequest(*masterCoinConfig, request)
	weight := analysis.CalculateWeight(result, *masterCoin)
	msg := fmt.Sprintf("weight log %s for coin %s: %.2f", datetime.Format("January 2, 2006 15:04:05"), config.Symbol, weight)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}
	weightLongInterval := crypto.GetOnLongIntervalWeight(*result, *masterCoin, timeInMili)
	msg += fmt.Sprintf("weight on long interval : %.2f", weightLongInterval)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetLongIntervalWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	return msg
}

func GetBalance() string {
	msg := "Balance status: \nWallet Balance:\n"
	format := "Symbol: <b>%s</b> \nBalance: <b>%f</b> \nEstimation In USDT: <b>%f</b> \n"

	walletBalances := crypto.GetWalletBalance()
	var totalWalletBalance float32 = 0
	for _, walb := range walletBalances {
		msg += fmt.Sprintf(format, walb["symbol"], walb["balance"], walb["estimation_usdt"])
		totalWalletBalance += walb["estimation_usdt"].(float32)
	}

	currentBalance := crypto.GetBalance()
	msg += fmt.Sprintf("\n\nCurrent Balance: <b>%f</b>", currentBalance)

	msg += fmt.Sprintf("\n\nTotal Estimation Balance: <b>%f</b>", currentBalance+totalWalletBalance)

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
