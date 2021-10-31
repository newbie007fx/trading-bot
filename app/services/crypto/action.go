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
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := MakeCryptoRequest(config, request)
	if result == nil {
		return "invalid requested date"
	}

	if resolution != "15m" {
		responseChanMid := make(chan CandleResponse)
		requestMid := CandleRequest{
			Symbol:       config.Symbol,
			StartDate:    0,
			EndDate:      timeInMili,
			Limit:        int(models.CandleLimit),
			Resolution:   resolution,
			ResponseChan: responseChanMid,
		}
		closeBand := result.Bands[len(result.Bands)-1]
		result = MakeCryptoRequestUpdateLasCandle(config, requestMid, closeBand.Candle.Close, closeBand.Candle.Hight)
		if result == nil {
			return "invalid requested date"
		}
	}

	msg := GenerateMsg(*result)
	if config.IsOnHold {
		msg += HoldCoinMessage(config, result)
	}

	return msg
}

func GetWeightLog(config models.CurrencyNotifConfig, datetime time.Time) string {
	timeInMili := datetime.Unix() * 1000
	var closeBand models.Band
	var closeBandMaster models.Band
	var result *models.BandResult
	var masterCoin *models.BandResult

	masterCoinConfig, _ := repositories.GetMasterCoinConfig()

	isTimeOnFifteenMinute := datetime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		result = CheckCoin(config, "15m", 0, timeInMili, nil)
		closeBand = result.Bands[len(result.Bands)-1]

		masterCoin = CheckCoin(*masterCoinConfig, "15m", 0, timeInMili, nil)
		closeBandMaster = masterCoin.Bands[len(masterCoin.Bands)-1]
	} else {
		oneMinuteResult := CheckCoin(config, "1m", 0, timeInMili, nil)
		closeBand = oneMinuteResult.Bands[len(oneMinuteResult.Bands)-1]
		result = CheckCoin(config, "15m", 0, timeInMili, &closeBand)

		oneMinuteMaster := CheckCoin(*masterCoinConfig, "1m", 0, timeInMili, nil)
		closeBandMaster = oneMinuteMaster.Bands[len(oneMinuteMaster.Bands)-1]
		masterCoin = CheckCoin(*masterCoinConfig, "15m", 0, timeInMili, &closeBandMaster)
	}

	masterCoinMid := CheckCoin(*masterCoinConfig, "1h", 0, timeInMili, &closeBandMaster)

	weight := analysis.CalculateWeight(result, *masterCoin)
	msg := GenerateMsg(*result)
	msg += fmt.Sprintf("\nweight log %s for coin %s: %.2f", datetime.Format("January 2, 2006 15:04:05"), config.Symbol, weight)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	resultMid := CheckCoin(config, "1h", 0, timeInMili, &closeBand)
	weightMid := analysis.CalculateWeightLongInterval(resultMid, masterCoin.Trend)
	msg += fmt.Sprintf("\nweight midInterval for coin %s: %.2f", config.Symbol, weightMid)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetLongIntervalWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	resultLong := CheckCoin(config, "4h", 0, timeInMili, &closeBand)
	weightLong := analysis.CalculateWeightLongInterval(resultLong, masterCoin.Trend)
	msg += fmt.Sprintf("\nweight long Interval for coin %s: %.2f", config.Symbol, weightLong)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetLongIntervalWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	shortIgnored := analysis.IsIgnored(result, masterCoin, datetime)
	msg += fmt.Sprintf("\nignord short interval: %t\n", shortIgnored)
	if shortIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	midIgnored := analysis.IsIgnoredMidInterval(resultMid, result)
	msg += fmt.Sprintf("ignord mid interval: %t\n", midIgnored)
	if midIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	longIgnored := analysis.IsIgnoredLongInterval(resultLong, result, resultMid, masterCoin.Trend, masterCoinMid.Trend)
	msg += fmt.Sprintf("ignord long interval: %t\n", longIgnored)
	if longIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	masterDownIgnored := analysis.IsIgnoredMasterDown(result, resultMid, resultLong, masterCoin, datetime)
	msg += fmt.Sprintf("ignord  master down: %t\n", masterDownIgnored)
	if masterDownIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	return msg
}

func GetSellLog(config models.CurrencyNotifConfig, datetime time.Time) string {
	timeInMili := datetime.Unix() * 1000
	var closeBand models.Band
	var closeBandMaster models.Band
	var coin *models.BandResult
	var masterCoin *models.BandResult

	masterCoinConfig, _ := repositories.GetMasterCoinConfig()

	isTimeOnFifteenMinute := datetime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		coin = CheckCoin(config, "15m", 0, timeInMili, nil)
		closeBand = coin.Bands[len(coin.Bands)-1]

		masterCoin = CheckCoin(*masterCoinConfig, "15m", 0, timeInMili, nil)
		closeBandMaster = masterCoin.Bands[len(masterCoin.Bands)-1]
	} else {
		oneMinuteResult := CheckCoin(config, "1m", 0, timeInMili, nil)
		closeBand = oneMinuteResult.Bands[len(oneMinuteResult.Bands)-1]
		coin = CheckCoin(config, "15m", 0, timeInMili, &closeBand)

		oneMinuteMaster := CheckCoin(*masterCoinConfig, "1m", 0, timeInMili, nil)
		closeBandMaster = oneMinuteMaster.Bands[len(oneMinuteMaster.Bands)-1]
		masterCoin = CheckCoin(*masterCoinConfig, "15m", 0, timeInMili, &closeBandMaster)
	}

	masterCoinMid := CheckCoin(*masterCoinConfig, "1h", 0, timeInMili, &closeBandMaster)

	coinMid := CheckCoin(config, "1h", 0, timeInMili, &closeBand)
	coinLong := CheckCoin(config, "1h", 0, timeInMili, &closeBand)
	isNeedTosell := analysis.IsNeedToSell(*coin, *masterCoin, isTimeOnFifteenMinute, coinMid.Trend, masterCoinMid.Trend)
	if isNeedTosell || analysis.SpecialCondition(coin.Symbol, *coin, *coinMid, *coinLong) {

		msg := fmt.Sprintf("sell log on %s:\n", datetime.Format("January 2, 2006 15:04:05"))
		msg += GenerateMsg(*coin)
		msg += "\n"
		msg += "alasan dijual: " + analysis.GetSellReason() + "\n\n"
		return msg
	}

	return ""
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
