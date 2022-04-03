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

	var closeBand models.Band
	var result *models.BandResult

	timeInMili := currentTime.Unix() * 1000

	isTimeOnFifteenMinute := currentTime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		result = CheckCoin(config.Symbol, "15m", 0, timeInMili, 0, 0, 0)
		closeBand = result.Bands[len(result.Bands)-1]
	} else {
		oneMinuteResult := CheckCoin(config.Symbol, "1m", 0, timeInMili, 0, 0, 0)
		closeBand = oneMinuteResult.Bands[len(oneMinuteResult.Bands)-1]
		higest := analysis.GetHighestHightPriceByTime(currentTime, oneMinuteResult.Bands, analysis.Time_type_15m, true)
		lowest := analysis.GetLowestLowPriceByTime(currentTime, oneMinuteResult.Bands, analysis.Time_type_15m, true)
		result = CheckCoin(config.Symbol, "15m", 0, timeInMili, closeBand.Candle.Open, higest, lowest)
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
		result = MakeCryptoRequestUpdateLasCandle(requestMid, closeBand.Candle.Close, closeBand.Candle.Hight, closeBand.Candle.Low)
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

func GetWeightLog(symbol string, datetime time.Time) string {
	timeInMili := datetime.Unix() * 1000

	result := CheckCoin(symbol, "15m", 0, timeInMili, 0, 0, 0)
	msg := GenerateMsg(*result)

	higest := analysis.GetHighestHightPriceByTime(datetime, result.Bands, analysis.Time_type_1h, true)
	lowest := analysis.GetLowestLowPriceByTime(datetime, result.Bands, analysis.Time_type_1h, true)
	resultMid := CheckCoin(symbol, "1h", 0, timeInMili, result.CurrentPrice, higest, lowest)

	msg += "\n" + GenerateMsg(*resultMid)

	higest = analysis.GetHighestHightPriceByTime(datetime, resultMid.Bands, analysis.Time_type_4h, true)
	lowest = analysis.GetLowestLowPriceByTime(datetime, resultMid.Bands, analysis.Time_type_4h, true)
	resultLong := CheckCoin(symbol, "4h", 0, timeInMili, resultMid.CurrentPrice, higest, lowest)

	msg += "\n" + GenerateMsg(*resultLong)

	shortIgnored := analysis.IgnoredOnUpTrendShort(*result)
	msg += fmt.Sprintf("\nignord short interval: %t\n", shortIgnored)
	if shortIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	midIgnored := analysis.IgnoredOnUpTrendMid(*resultMid, *result)
	msg += fmt.Sprintf("ignord mid interval: %t\n", midIgnored)
	if midIgnored {
		msg += fmt.Sprintf("ignord reason: %s\n", analysis.GetIgnoredReason())
	}

	longIgnored := analysis.IgnoredOnUpTrendLong(*resultLong, *resultMid, *result)
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
		coin = CheckCoin(config.Symbol, "15m", 0, timeInMili, 0, 0, 0)
	} else {
		oneMinuteResult := CheckCoin(config.Symbol, "1m", 0, timeInMili, 0, 0, 0)
		higest := analysis.GetHighestHightPriceByTime(datetime, oneMinuteResult.Bands, analysis.Time_type_15m, true)
		lowest := analysis.GetLowestLowPriceByTime(datetime, oneMinuteResult.Bands, analysis.Time_type_15m, true)
		coin = CheckCoin(config.Symbol, "15m", 0, timeInMili, oneMinuteResult.CurrentPrice, higest, lowest)
	}

	log.Println("close price", coin.Bands[len(coin.Bands)-1].Candle.Close)
	log.Println("higest price", coin.Bands[len(coin.Bands)-1].Candle.Hight)

	higest := analysis.GetHighestHightPriceByTime(datetime, coin.Bands, analysis.Time_type_1h, true)
	lowest := analysis.GetLowestLowPriceByTime(datetime, coin.Bands, analysis.Time_type_1h, true)
	coinMid := CheckCoin(config.Symbol, "1h", 0, timeInMili, coin.CurrentPrice, higest, lowest)
	higest = analysis.GetHighestHightPriceByTime(datetime, coinMid.Bands, analysis.Time_type_4h, true)
	lowest = analysis.GetLowestLowPriceByTime(datetime, coinMid.Bands, analysis.Time_type_4h, true)
	coinLong := CheckCoin(config.Symbol, "4h", 0, timeInMili, coinMid.CurrentPrice, higest, lowest)
	if analysis.CheckIsNeedSellOnTrendUp(&config, *coin, *coinMid, *coinLong, datetime) {
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
	format := "Symbol: <b>%s</b> \nBalance: <b>%f</b> \nEstimation In USDT: <b>%f</b> \nPercent Changes: <b>%f</b> \n"

	walletBalances := GetWalletBalance()
	var totalWalletBalance float32 = 0
	for _, walb := range walletBalances {
		msg += fmt.Sprintf(format, walb["symbol"], walb["balance"], walb["estimation_usdt"], walb["percent_change"])
		totalWalletBalance += walb["estimation_usdt"].(float32)
	}

	currentBalance := GetBalanceFromConfig()
	msg += fmt.Sprintf("\n\nCurrent Balance: <b>%f</b>", currentBalance)

	msg += fmt.Sprintf("\n\nTotal Estimation Balance: <b>%f</b>", currentBalance+totalWalletBalance)

	return msg
}
