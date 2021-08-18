package trading_strategy

import (
	"sort"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

type AutomaticTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
}

func (ats *AutomaticTradingStrategy) Execute(currentTime time.Time) {
	condition := map[string]interface{}{"is_on_hold": true}
	hold_count := repositories.CountNotifConfig(&condition)
	if hold_count > 0 {
		ats.cryptoHoldCoinPriceChan <- true
	} else {
		if ats.isTimeToCheckAltCoinPrice(currentTime) {
			ats.cryptoAltCoinPriceChan <- true
		}
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

func (AutomaticTradingStrategy) isTimeToCheckAltCoinPrice(time time.Time) bool {
	minute := time.Minute()
	var listMinutes []int = []int{5, 15, 20, 30, 35, 45, 50, 0}
	for _, a := range listMinutes {
		if a == minute {
			return true
		}
	}

	return false
}

func (AutomaticTradingStrategy) startCheckHoldCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		holdCoin := checkCryptoHoldCoinPrice()
		msg := ""
		if len(holdCoin) > 0 {
			for _, coin := range holdCoin {
				if analysis.IsNeedToSell(coin) {
					msg += "coin berikut akan dijual:\n"
					msg += generateMsg(coin)
					msg += "\n"

					currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
					if err == nil {
						bands := coin.Bands
						lastBand := bands[len(bands)-1]
						services.ReleaseCoin(*currencyConfig, lastBand.Candle)
					}
				}
			}

			waitMasterCoinProcessed()
			if masterCoin != nil && msg != "" {
				msg += "untuk master coin:\n"
				msg += generateMsg(*masterCoin)
			}
		}

		sendNotif(msg)
	}
}

func (AutomaticTradingStrategy) startCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		altCoins := checkCryptoAltCoinPrice()
		msg := ""
		if len(altCoins) > 0 {
			sort.Slice(altCoins, func(i, j int) bool { return altCoins[i].Weight > altCoins[j].Weight })
			coin := altCoins[0]

			msg = "coin berikut telah dihold:\n"
			msg += generateMsg(coin)
			msg += "\n"

			if masterCoin != nil {
				msg += "untuk master coin:\n"
				msg += generateMsg(*masterCoin)
			}

			currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
			if err == nil {
				bands := coin.Bands
				lastBand := bands[len(bands)-1]
				services.HoldCoin(*currencyConfig, lastBand.Candle)
			}
		}

		sendNotif(msg)
	}
}
