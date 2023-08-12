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
var modeChecking string = ""
var isSkiped bool = false
var checkLimitHistory []int = []int{}

type AutomaticTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
}

func (ats *AutomaticTradingStrategy) Execute(currentTime time.Time) {
	baseCheckingTime = currentTime

	maxHold := crypto.GetMaxHold()
	condition := map[string]interface{}{"is_on_hold": true}
	holdCount = repositories.CountNotifConfig(&condition)

	if holdCount > 0 {
		ats.cryptoHoldCoinPriceChan <- true
	}

	if crypto.IsProfitMoreThanThreshold() {
		if !isSkiped {
			log.Println("skipped: profit is more than threshold")
			isSkiped = true
		}
		return
	}

	isSkiped = false
	if holdCount < maxHold && ats.isTimeToCheckAltCoinPrice(currentTime) {
		log.Printf("execute automatic trading, with hold count: %d and maxHold %d", holdCount, maxHold)
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
	if currentTime.Second() > 5 {
		return false
	}

	return minute == 0 || minute == 15 || minute == 30 || minute == 45
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

				if analysis.CheckIsNeedSellOnTrendUp(currencyConfig, coin, baseCheckingTime) {
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

						var changes float32 = 0
						if currencyConfig.HoldPrice < coin.CurrentPrice {
							changes = (coin.CurrentPrice - currencyConfig.HoldPrice) / currencyConfig.HoldPrice * 100
						} else {
							changes = -((currencyConfig.HoldPrice - coin.CurrentPrice) / currencyConfig.HoldPrice * 100)
						}
						crypto.SetProfit(changes)
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

		coins := ats.checkOnTrendUp()
		if len(coins) > 0 {
			msg := generateNeedToCheckMessage(coins)

			crypto.SendNotif(msg)
			crypto.NotifyCounterDecrement()
		}
	}
}

func setLimitCheckOnTrendUp() {
	var limit int = crypto.GetLimit()
	if limit > 60 {
		limit = 60
	}

	if len(checkLimitHistory) < 13 {
		checkLimitHistory = append(checkLimitHistory, limit)
	} else {
		checkLimitHistory = append(checkLimitHistory[1:], limit)
		log.Printf("check limit history: %v\n", checkLimitHistory)
	}

	modeChecking = crypto.GetModeChecking()
	checkOnTrendUpLimit = limit

	if isNoNeedDoubleCheck() {
		crypto.SetNotifyCounter(1)
	}
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

func generateNeedToCheckMessage(coins []models.BandResult) string {
	msg := fmt.Sprintf("silahkan cek list koin berikut %d:\n", altCheckingTime.Unix())
	for i, coin := range coins {
		msg += fmt.Sprintf("%d. %s", i+1, crypto.GenerateMsg(coin))
		msg += "\n"
	}

	return msg
}

func (ats *AutomaticTradingStrategy) checkOnTrendUp() []models.BandResult {
	timeInMilli := GetEndDate(altCheckingTime, OPERATION_BUY)
	altCoins := checkCoinOnTrendUp(altCheckingTime)
	sort.SliceStable(altCoins, func(i, j int) bool {
		return altCoins[i].PriceChanges > altCoins[j].PriceChanges
	})

	var coins []models.BandResult = []models.BandResult{}
	if crypto.IsShouldNotify() {
		coinsString := ""
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

			if analysis.ApprovedPattern(coin, *resultMid, *resultLong, altCheckingTime, isNoNeedDoubleCheck()) {
				coins = append(coins, coin)
				coinsString += ", " + coin.Symbol
			}

			if len(coins) == 3 {
				break
			}
		}
		log.Println("found coins: ", coinsString)
	} else {
		log.Println("notify counter is zero, skipping checking detail")
	}

	return coins
}

func sendHoldMsg(result *models.BandResult) string {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(result.Symbol)
	if err != nil {
		return ""
	}
	return crypto.HoldCoinMessage(*currencyConfig, result)
}

func isNoNeedDoubleCheck() bool {
	if len(checkLimitHistory) >= 13 {
		last := checkLimitHistory[len(checkLimitHistory)-1]
		if last > 5 {
			total := 0
			for _, limit := range checkLimitHistory[:len(checkLimitHistory)-1] {
				total += limit
			}
			average := total / (len(checkLimitHistory) - 1)

			return last > (average * 2)
		}
	}

	return false
}
