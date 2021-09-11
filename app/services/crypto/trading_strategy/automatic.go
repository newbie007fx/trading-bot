package trading_strategy

import (
	"fmt"
	"sort"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var checkingTime time.Time

type AutomaticTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
}

func (ats *AutomaticTradingStrategy) Execute(currentTime time.Time) {
	checkingTime = currentTime

	condition := map[string]interface{}{"is_on_hold": true}
	holdCount := repositories.CountNotifConfig(&condition)
	if holdCount > 0 {
		ats.cryptoHoldCoinPriceChan <- true
	}

	maxHold := crypto.GetMaxHold()
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
		holdCoin := checkCryptoHoldCoinPrice()
		msg := ""
		if len(holdCoin) > 0 {

			waitMasterCoinProcessed()

			tmpMsg := ""
			for _, coin := range holdCoin {
				if analysis.IsNeedToSell(coin, *masterCoin, ats.isTimeToCheckAltCoinPrice(checkingTime)) {
					tmpMsg = "coin berikut akan dijual:\n"
					tmpMsg += crypto.GenerateMsg(coin)
					tmpMsg += "\n"
					tmpMsg += "alasan dijual: " + analysis.GetSellReason() + "\n\n"

					balance := crypto.GetBalance()
					tmpMsg += fmt.Sprintf("saldo saat ini: %f\n", balance)

					currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
					if err == nil {
						bands := coin.Bands
						lastBand := bands[len(bands)-1]
						err = services.ReleaseCoin(*currencyConfig, lastBand.Candle)
						if err != nil {
							tmpMsg = err.Error()
						}
					}
					msg += tmpMsg
				}
			}

			if masterCoin != nil && msg != "" {
				msg += "untuk master coin:\n"
				msg += crypto.GenerateMsg(*masterCoin)
			}
		}

		sendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) startCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		altCoins := checkCryptoAltCoinPrice(&checkingTime)
		msg := ""
		if len(altCoins) > 0 {

			coin := ats.sortAndGetHigest(altCoins)
			if coin == nil {
				continue
			}

			msg = "coin berikut telah dihold:\n"
			msg += crypto.GenerateMsg(*coin)
			msg += fmt.Sprintf("weight: <b>%.2f</b>\n", coin.Weight)
			msg += "\n"

			if masterCoin != nil {
				msg += "untuk master coin:\n"
				msg += crypto.GenerateMsg(*masterCoin)
			}

			currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
			if err == nil {
				bands := coin.Bands
				lastBand := bands[len(bands)-1]
				err = services.HoldCoin(*currencyConfig, lastBand.Candle)
				if err != nil {
					msg = err.Error()
				}
			}
		}

		sendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) sortAndGetHigest(altCoins []models.BandResult) *models.BandResult {
	results := []models.BandResult{}
	timeInMilli := GetEndDate(&checkingTime)
	for i := range altCoins {
		waitMasterCoinProcessed()
		altCoins[i].Weight += crypto.GetOnLongIntervalWeight(altCoins[i], *masterCoin, timeInMilli)
		if altCoins[i].Weight > 2.95 {
			results = append(results, altCoins[i])
		}
	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool { return results[i].Weight > results[j].Weight })

		return &results[0]
	}
	return nil
}
