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

func GetWeightLog(symbol string, datetime time.Time, modeChecking string) string {
	var utcZone, _ = time.LoadLocation("UTC")
	datetime = datetime.In(utcZone)
	timeInMili := datetime.Unix() * 1000

	var result *models.BandResult

	isTimeOnFifteenMinute := datetime.Unix()%(15*60) == 0
	if isTimeOnFifteenMinute {
		result = CheckCoin(symbol, "15m", 0, timeInMili, 0, 0, 0)
	} else {
		oneMinuteResult := CheckCoin(symbol, "1m", 0, timeInMili, 0, 0, 0)
		if oneMinuteResult != nil {
			higest := analysis.GetHighestHightPriceByTime(datetime, oneMinuteResult.Bands, analysis.Time_type_15m, true)
			lowest := analysis.GetLowestLowPriceByTime(datetime, oneMinuteResult.Bands, analysis.Time_type_15m, true)
			result = CheckCoin(symbol, "15m", 0, timeInMili, oneMinuteResult.CurrentPrice, higest, lowest)
		}
	}

	msg := "some issue happend. please check"

	if result != nil {
		msg = GenerateMsg(*result)

		higest := analysis.GetHighestHightPriceByTime(datetime, result.Bands, analysis.Time_type_1h, true)
		lowest := analysis.GetLowestLowPriceByTime(datetime, result.Bands, analysis.Time_type_1h, true)
		resultMid := CheckCoin(symbol, "1h", 0, timeInMili, result.CurrentPrice, higest, lowest)

		if resultMid != nil {
			msg += "\n" + GenerateMsg(*resultMid)

			higest = analysis.GetHighestHightPriceByTime(datetime, resultMid.Bands, analysis.Time_type_4h, true)
			lowest = analysis.GetLowestLowPriceByTime(datetime, resultMid.Bands, analysis.Time_type_4h, true)
			resultLong := CheckCoin(symbol, "4h", 0, timeInMili, resultMid.CurrentPrice, higest, lowest)

			if resultLong != nil {
				msg += "\n" + GenerateMsg(*resultLong)

				if analysis.ApprovedPattern(*result, *resultMid, *resultLong, datetime, modeChecking) {
					msg += "\npattern: " + analysis.GetMatchPattern()
				} else {
					msg += "\nignored reason: " + analysis.GetIgnoredReason()
				}
			}
		}
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

	if analysis.CheckIsNeedSellOnTrendUp(&config, *coin, datetime) {
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
	format := "Symbol: <b>%s</b> \nBalance: <b>%f</b> \nEstimation In USDT: <b>%f</b> \nPercent Changes: <b>%f</b> \nHold Duration: <b>%d minutes</b> \n"

	walletBalances := GetWalletBalance()
	var totalWalletBalance float32 = 0
	for _, walb := range walletBalances {

		msg += fmt.Sprintf(format, walb["symbol"], walb["balance"], walb["estimation_usdt"], walb["percent_change"], walb["hold_duration"])
		totalWalletBalance += walb["estimation_usdt"].(float32)
	}

	currentBalance := GetBalanceFromConfig()
	msg += fmt.Sprintf("\n\nCurrent Balance: <b>%f</b>", currentBalance)

	msg += fmt.Sprintf("\n\nTotal Estimation Balance: <b>%f</b>", currentBalance+totalWalletBalance)

	return msg
}
