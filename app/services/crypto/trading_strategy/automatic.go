package trading_strategy

import (
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var baseCheckingTime time.Time
var altCheckingTime time.Time
var holdCount int64 = 0

type AutomaticTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
}

func (ats *AutomaticTradingStrategy) Execute(currentTime time.Time) {
	baseCheckingTime = currentTime

	maxHold := crypto.GetMaxHold()
	condition := map[string]interface{}{"is_on_hold": true}
	holdCount = repositories.CountNotifConfig(&condition)

	log.Println(fmt.Sprintf("execute automatic trading, with hold count: %d and maxHold %d", holdCount, maxHold))
	if holdCount > 0 {
		ats.cryptoHoldCoinPriceChan <- true
	}

	if holdCount < maxHold && ats.isTimeToCheckAltCoinPrice(currentTime) {
		ats.cryptoAltCoinPriceChan <- true
	}

}

func (ats *AutomaticTradingStrategy) InitService() {
	ats.cryptoHoldCoinPriceChan = make(chan bool)
	ats.cryptoAltCoinPriceChan = make(chan bool)

	go ats.startCheckHoldCoinPriceService(ats.cryptoHoldCoinPriceChan)
	go ats.startCheckAltCoinPriceService(ats.cryptoAltCoinPriceChan)
}

func (ats *AutomaticTradingStrategy) Shutdown() {
	close(ats.cryptoHoldCoinPriceChan)
	close(ats.cryptoAltCoinPriceChan)
}

func (AutomaticTradingStrategy) isTimeToCheckAltCoinPrice(currentTime time.Time) bool {
	minute := currentTime.Minute()
	if currentTime.Second() > 25 {
		return false
	}

	var listMinutes []int = []int{5, 20, 35, 50}
	for _, a := range listMinutes {
		if a == minute {
			return true
		}
	}

	return false
}

func (ats *AutomaticTradingStrategy) startCheckHoldCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		holdCoin := checkCryptoHoldCoinPrice(baseCheckingTime)
		sellTime := GetEndDate(baseCheckingTime, OPERATION_SELL)
		msg := ""
		if len(holdCoin) > 0 {

			tmpMsg := ""
			for _, coin := range holdCoin {
				currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
				if err != nil {
					log.Println(err.Error())
					continue
				}

				holdCoinMid := crypto.CheckCoin(currencyConfig.Symbol, "1h", 0, sellTime, 0, 0, 0)
				holdCoinLong := crypto.CheckCoin(currencyConfig.Symbol, "4h", 0, sellTime, 0, 0, 0)
				if holdCoinMid == nil || holdCoinLong == nil {
					log.Println("error hold coin nil. skip need to sell checking process")
					continue
				}

				if analysis.CheckIsNeedSellOnTrendUp(currencyConfig, coin) {
					bands := coin.Bands
					lastBand := bands[len(bands)-1]
					err = crypto.ReleaseCoin(*currencyConfig, lastBand.Candle)
					if err != nil {
						tmpMsg = err.Error()
					} else {
						tmpMsg = fmt.Sprintf("coin berikut akan dijual %d:\n", sellTime)
						tmpMsg += crypto.GenerateMsg(coin)
						tmpMsg += "\n"
						tmpMsg += crypto.HoldCoinMessage(*currencyConfig, &coin)
						tmpMsg += "\n"
						tmpMsg += "alasan dijual: " + analysis.GetSellReason() + "\n\n"

						balance := crypto.GetBalanceFromConfig()
						tmpMsg += fmt.Sprintf("saldo saat ini: %f\n", balance)
					}
					msg += tmpMsg
				}
			}
		}

		crypto.SendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) startCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		altCheckingTime = baseCheckingTime
		setLimitCheckOnTrendUp()
		msg := ""

		maxHold := crypto.GetMaxHold()
		if coin := ats.checkOnTrendUp(); coin != nil {
			if holdCount < maxHold {
				if ok, resMsg := holdAndGenerateMessage(coin); ok {
					msg += resMsg
					msg += "pattern: " + analysis.GetIgnoredReason() + " \n\n"

					holdCount++
				}
			}
		}

		crypto.SendNotif(msg)
	}
}

func setLimitCheckOnTrendUp() {
	var limit int = crypto.GetLimit()
	if limit > 40 {
		limit = 40
	}
	if limit < 4 {
		limit = 10
	}

	checkOnTrendUpLimit = limit
}

func holdAndGenerateMessage(coin *models.BandResult) (bool, string) {
	msg := ""
	hold := true
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
	if err == nil {
		bands := coin.Bands
		lastBand := bands[len(bands)-1]
		err = crypto.HoldCoin(*currencyConfig, lastBand.Candle)
		if err != nil {
			msg = err.Error()
			hold = false
		} else {
			msg += fmt.Sprintf("coin berikut telah dihold on %d:\n", altCheckingTime.Unix())
			msg += crypto.GenerateMsg(*coin)
			msg += "\n"
			msg += sendHoldMsg(coin)
			msg += "\n"
			msg += "coin mid interval:\n"
			msg += crypto.GenerateMsg(*coin.Mid)
			msg += "\n"
			msg += "coin long interval:\n"
			msg += crypto.GenerateMsg(*coin.Long)
			msg += "\n"
		}
	}
	return hold, msg
}

func (ats *AutomaticTradingStrategy) checkOnTrendUp() *models.BandResult {
	timeInMilli := GetEndDate(altCheckingTime, OPERATION_BUY)
	altCoins := checkCoinOnTrendUp(altCheckingTime)
	for _, coin := range altCoins {

		higest := analysis.GetHighestHightPriceByTime(altCheckingTime, coin.Bands, analysis.Time_type_1h, false)
		lowest := analysis.GetLowestLowPriceByTime(altCheckingTime, coin.Bands, analysis.Time_type_1h, false)
		resultMid := crypto.CheckCoin(coin.Symbol, "1h", 0, timeInMilli, coin.CurrentPrice, higest, lowest)

		if resultMid.Direction == analysis.BAND_DOWN {
			continue
		}

		coin.Mid = resultMid

		higest = analysis.GetHighestHightPriceByTime(altCheckingTime, resultMid.Bands, analysis.Time_type_4h, false)
		lowest = analysis.GetLowestLowPriceByTime(altCheckingTime, resultMid.Bands, analysis.Time_type_4h, false)
		resultLong := crypto.CheckCoin(coin.Symbol, "4h", 0, timeInMilli, resultMid.CurrentPrice, higest, lowest)

		if resultLong.Direction == analysis.BAND_DOWN {
			continue
		}

		coin.Long = resultLong

		if analysis.ApprovedPattern(coin, *resultMid, *resultLong, altCheckingTime) {
			return &coin
		}
	}

	return nil
}

func sendHoldMsg(result *models.BandResult) string {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(result.Symbol)
	if err != nil {
		return ""
	}
	return crypto.HoldCoinMessage(*currencyConfig, result)
}
