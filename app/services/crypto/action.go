package crypto

import (
	"errors"
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

func HoldCoin(currencyConfig models.CurrencyNotifConfig, candleData *models.CandleData) error {
	if GetMode() != "manual" {
		err := Buy(currencyConfig, candleData)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	if !currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": true,
			"holded_at":  time.Now().Unix(),
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
	if GetMode() != "manual" {
		err := Sell(currencyConfig, candleData)
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

func GetCurrencyStatus(config models.CurrencyNotifConfig, resolution string, requestTime *time.Time) string {
	currentTime := time.Now()
	if requestTime != nil {
		currentTime = *requestTime
	}

	timeInMili := currentTime.Unix() * 1000

	responseChan := make(chan CandleResponse)
	request := CandleRequest{
		Symbol:       config.Symbol,
		EndDate:      timeInMili,
		Limit:        40,
		Resolution:   resolution,
		ResponseChan: responseChan,
	}

	result := MakeCryptoRequest(config, request)
	if result == nil {
		return "invalid requested date"
	}

	msg := GenerateMsg(*result)
	if config.IsOnHold {
		msg += HoldCoinMessage(config, result)
	}

	return msg
}

func GetWeightLog(config models.CurrencyNotifConfig, datetime time.Time) string {
	timeInMili := datetime.Unix() * 1000

	responseChan := make(chan CandleResponse)
	request := CandleRequest{
		Symbol:       config.Symbol,
		StartDate:    0,
		EndDate:      timeInMili,
		Limit:        int(models.CandleLimit),
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := MakeCryptoRequest(config, request)

	masterCoinConfig, _ := repositories.GetMasterCoinConfig()
	request = CandleRequest{
		Symbol:       masterCoinConfig.Symbol,
		StartDate:    0,
		EndDate:      timeInMili,
		Limit:        int(models.CandleLimit),
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	masterCoin := MakeCryptoRequest(*masterCoinConfig, request)
	weight := analysis.CalculateWeight(result, *masterCoin)
	msg := GenerateMsg(*result)
	msg += fmt.Sprintf("\nweight log %s for coin %s: %.2f", datetime.Format("January 2, 2006 15:04:05"), config.Symbol, weight)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	return msg
}

func GetBalances() string {
	msg := "Balance status: \nWallet Balance:\n"
	format := "Symbol: <b>%s</b> \nBalance: <b>%f</b> \nEstimation In USDT: <b>%f</b> \n"

	walletBalances := GetWalletBalance()
	var totalWalletBalance float32 = 0
	for _, walb := range walletBalances {
		msg += fmt.Sprintf(format, walb["symbol"], walb["balance"], walb["estimation_usdt"])
		totalWalletBalance += walb["estimation_usdt"].(float32)
	}

	currentBalance := GetBalanceFromConfig()
	msg += fmt.Sprintf("\n\nCurrent Balance: <b>%f</b>", currentBalance)

	msg += fmt.Sprintf("\n\nTotal Estimation Balance: <b>%f</b>", currentBalance+totalWalletBalance)

	return msg
}
