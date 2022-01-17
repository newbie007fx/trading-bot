package crypto

import (
	"errors"
	"fmt"
	"log"
	"math"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

const time_type_15m int = 1
const time_type_1h int = 2
const time_type_4h int = 3

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

	var closeBand models.Band
	var result *models.BandResult

	timeInMili := currentTime.Unix() * 1000

	isTimeOnFifteenMinute := currentTime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		result = CheckCoin(config, "15m", 0, timeInMili, 0, 0, 0)
		closeBand = result.Bands[len(result.Bands)-1]
	} else {
		oneMinuteResult := CheckCoin(config, "1m", 0, timeInMili, 0, 0, 0)
		closeBand = oneMinuteResult.Bands[len(oneMinuteResult.Bands)-1]
		higest := getHighestHightPrice(currentTime, oneMinuteResult.Bands, time_type_15m)
		lowest := getLowestLowPrice(currentTime, oneMinuteResult.Bands, time_type_15m)
		result = CheckCoin(config, "15m", 0, timeInMili, closeBand.Candle.Open, higest, lowest)
	}

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
		result = MakeCryptoRequestUpdateLasCandle(config, requestMid, closeBand.Candle.Close, closeBand.Candle.Hight, closeBand.Candle.Low)
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
	var result *models.BandResult

	isTimeOnFifteenMinute := datetime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		result = CheckCoin(config, "15m", 0, timeInMili, 0, 0, 0)
	} else {
		oneMinuteResult := CheckCoin(config, "1m", 0, timeInMili, 0, 0, 0)
		higest := getHighestHightPrice(datetime, oneMinuteResult.Bands, time_type_15m)
		lowest := getLowestLowPrice(datetime, oneMinuteResult.Bands, time_type_15m)
		result = CheckCoin(config, "15m", 0, timeInMili, oneMinuteResult.CurrentPrice, higest, lowest)
	}

	weight := analysis.CalculateWeight(result)
	msg := GenerateMsg(*result)
	msg += fmt.Sprintf("\nweight log %s for coin %s: %.2f", datetime.Format("January 2, 2006 15:04:05"), config.Symbol, weight)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	higest := getHighestHightPrice(datetime, result.Bands, time_type_1h)
	lowest := getLowestLowPrice(datetime, result.Bands, time_type_1h)
	resultMid := CheckCoin(config, "1h", 0, timeInMili, result.CurrentPrice, higest, lowest)
	weightMid := analysis.CalculateWeightLongInterval(resultMid)
	msg += fmt.Sprintf("\nweight midInterval for coin %s: %.2f", config.Symbol, weightMid)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetLongIntervalWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	higest = getHighestHightPrice(datetime, resultMid.Bands, time_type_4h)
	lowest = getLowestLowPrice(datetime, resultMid.Bands, time_type_4h)
	resultLong := CheckCoin(config, "4h", 0, timeInMili, resultMid.CurrentPrice, higest, lowest)
	weightLong := analysis.CalculateWeightLongInterval(resultLong)
	msg += fmt.Sprintf("\nweight long Interval for coin %s: %.2f", config.Symbol, weightLong)
	msg += "\n"
	msg += "detail weight: \n"
	for key, val := range analysis.GetLongIntervalWeightLogData() {
		msg += fmt.Sprintf("%s: %.2f\n", key, val)
	}

	shortIgnored := analysis.IsIgnored(result, datetime)
	msg += fmt.Sprintf("\nignord short interval: %t\n", shortIgnored)
	if shortIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	midIgnored := analysis.IsIgnoredMidInterval(resultMid, result)
	msg += fmt.Sprintf("ignord mid interval: %t\n", midIgnored)
	if midIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	longIgnored := analysis.IsIgnoredLongInterval(resultLong, result, resultMid)
	msg += fmt.Sprintf("ignord long interval: %t\n", longIgnored)
	if longIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	return msg
}

func GetSellLog(config models.CurrencyNotifConfig, datetime time.Time) string {
	timeInMili := datetime.Unix() * 1000
	var coin *models.BandResult

	isTimeOnFifteenMinute := datetime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		coin = CheckCoin(config, "15m", 0, timeInMili, 0, 0, 0)
	} else {
		oneMinuteResult := CheckCoin(config, "1m", 0, timeInMili, 0, 0, 0)
		higest := getHighestHightPrice(datetime, oneMinuteResult.Bands, time_type_15m)
		lowest := getLowestLowPrice(datetime, oneMinuteResult.Bands, time_type_15m)
		coin = CheckCoin(config, "15m", 0, timeInMili, oneMinuteResult.CurrentPrice, higest, lowest)
	}

	log.Println("close price", coin.Bands[len(coin.Bands)-1].Candle.Close)
	log.Println("higest price", coin.Bands[len(coin.Bands)-1].Candle.Hight)

	higest := getHighestHightPrice(datetime, coin.Bands, time_type_1h)
	lowest := getLowestLowPrice(datetime, coin.Bands, time_type_1h)
	coinMid := CheckCoin(config, "1h", 0, timeInMili, coin.CurrentPrice, higest, lowest)
	higest = getHighestHightPrice(datetime, coinMid.Bands, time_type_4h)
	lowest = getLowestLowPrice(datetime, coinMid.Bands, time_type_4h)
	coinLong := CheckCoin(config, "4h", 0, timeInMili, coinMid.CurrentPrice, higest, lowest)
	isNeedTosell := analysis.IsNeedToSell(&config, *coin, datetime, coinMid)
	if isNeedTosell || analysis.SpecialCondition(&config, coin.Symbol, *coin, *coinMid, *coinLong) {
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

func getHighestHightPrice(currentTime time.Time, bands []models.Band, timeType int) float32 {
	var numberBands int = 0
	var utcZone, _ = time.LoadLocation("UTC")
	currentTime = currentTime.In(utcZone)

	if timeType == time_type_15m {
		numberBands = (currentTime.Minute() + 1) % 15
		if numberBands == 0 {
			numberBands = 15
		}
	} else if timeType == time_type_1h {
		numberBands = int(math.Ceil(float64(currentTime.Minute()+1) / 15))
	} else {
		numberBands = (currentTime.Hour() + 1) % 4
		if numberBands == 0 {
			numberBands = 4
		}
	}

	return analysis.GetHigestHightPrice(bands[len(bands)-numberBands:])
}

func getLowestLowPrice(currentTime time.Time, bands []models.Band, timeType int) float32 {
	var numberBands int = 0
	var utcZone, _ = time.LoadLocation("UTC")
	currentTime = currentTime.In(utcZone)

	if timeType == time_type_15m {
		numberBands = (currentTime.Minute() + 1) % 15
		if numberBands == 0 {
			numberBands = 15
		}
	} else if timeType == time_type_1h {
		numberBands = int(math.Ceil(float64(currentTime.Minute()+1) / 15))
	} else {
		numberBands = (currentTime.Hour() + 1) % 4
		if numberBands == 0 {
			numberBands = 4
		}
	}

	return analysis.GetLowestLowPrice(bands[len(bands)-numberBands:])
}
