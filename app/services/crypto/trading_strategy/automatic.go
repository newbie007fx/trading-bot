package trading_strategy

import (
	"fmt"
	"log"
	"sort"
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
	var listMinutes []int = []int{15, 30, 45, 0}
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

				holdCoinMid := crypto.CheckCoin(*currencyConfig, "1h", 0, sellTime, 0, 0, 0)
				holdCoinLong := crypto.CheckCoin(*currencyConfig, "4h", 0, sellTime, 0, 0, 0)
				if holdCoinMid == nil || holdCoinLong == nil {
					log.Println("error hold coin nil. skip need to sell checking process")
					continue
				}
				isNeedTosell := analysis.IsNeedToSell(currencyConfig, coin, baseCheckingTime, holdCoinMid)
				if isNeedTosell || analysis.SpecialCondition(currencyConfig, coin.Symbol, coin, *holdCoinMid, *holdCoinLong) {
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

		altCoins := checkCryptoAltCoinPrice(altCheckingTime)
		msg := ""
		if len(altCoins) > 0 {
			coins := ats.sortAndGetHigest(altCoins)
			if coins == nil {
				continue
			}

			maxHold := crypto.GetMaxHold()
			for _, coin := range *coins {
				if holdCount == maxHold {
					break
				}

				currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
				if err == nil {
					bands := coin.Bands
					lastBand := bands[len(bands)-1]
					err = crypto.HoldCoin(*currencyConfig, lastBand.Candle)
					if err != nil {
						msg = err.Error()
					} else {
						msg += fmt.Sprintf("coin berikut telah dihold on %d:\n", altCheckingTime.Unix())
						msg += crypto.GenerateMsg(coin)
						msg += fmt.Sprintf("weight: <b>%.2f</b>\n", coin.Weight)
						msg += "\n"
						msg += sendHoldMsg(&coin)
						msg += "\n"
						msg += "coin mid interval:\n"
						msg += crypto.GenerateMsg(*coin.Mid)
						msg += "\n"
						msg += "coin long interval:\n"
						msg += crypto.GenerateMsg(*coin.Long)
						msg += "\n"
						msg += "\n"

						holdCount++
					}
				}
			}
		}

		crypto.SendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) sortAndGetHigest(altCoins []models.BandResult) *[]models.BandResult {
	results := []models.BandResult{}
	timeInMilli := GetEndDate(altCheckingTime, OPERATION_BUY)
	for i := range altCoins {
		currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(altCoins[i].Symbol)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		higest := analysis.GetHighestHightPriceByTime(altCheckingTime, altCoins[i].Bands, analysis.Time_type_1h)
		lowest := analysis.GetLowestLowPriceByTime(altCheckingTime, altCoins[i].Bands, analysis.Time_type_1h)
		resultMid := crypto.CheckCoin(*currencyConfig, "1h", 0, timeInMilli, altCoins[i].CurrentPrice, higest, lowest)

		if resultMid.AllTrend.SecondTrend == models.TREND_DOWN && resultMid.AllTrend.ShortTrend == models.TREND_DOWN && resultMid.Direction == analysis.BAND_DOWN {
			continue
		}

		midWeight := getWeightCustomInterval(*resultMid, altCoins[i], "1h", nil)
		if midWeight == 0 {
			continue
		}
		altCoins[i].Mid = resultMid
		altCoins[i].Weight += midWeight

		higest = analysis.GetHighestHightPriceByTime(altCheckingTime, resultMid.Bands, analysis.Time_type_4h)
		lowest = analysis.GetLowestLowPriceByTime(altCheckingTime, resultMid.Bands, analysis.Time_type_4h)
		resultLong := crypto.CheckCoin(*currencyConfig, "4h", 0, timeInMilli, resultMid.CurrentPrice, higest, lowest)
		longWight := getWeightCustomInterval(*resultLong, altCoins[i], "4h", resultMid)
		if longWight == 0 {
			continue
		}
		altCoins[i].Long = resultLong
		altCoins[i].Weight += longWight
		results = append(results, altCoins[i])

	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool { return results[i].Weight > results[j].Weight })

		return &results
	}
	return nil
}

func getWeightCustomInterval(result, coin models.BandResult, interval string, previous *models.BandResult) float32 {
	ignored := false

	if interval == "1h" {
		ignored = analysis.IsIgnoredMidInterval(&result, &coin)
	} else {
		ignored = analysis.IsIgnoredLongInterval(&result, &coin, previous)
	}

	if ignored {
		return 0
	}

	return 1
}

func sendHoldMsg(result *models.BandResult) string {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(result.Symbol)
	if err != nil {
		return ""
	}
	return crypto.HoldCoinMessage(*currencyConfig, result)
}
