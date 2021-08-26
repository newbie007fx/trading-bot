package trading_strategy

import (
	"sort"
	"telebot-trading/app/models"
	"telebot-trading/app/services/crypto"
	"time"
)

type ManualTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
}

func (mts *ManualTradingStrategy) Execute(currentTime time.Time) {
	minute := currentTime.Minute()
	if minute%5 == 0 {
		mts.cryptoHoldCoinPriceChan <- true
	}

	if mts.isTimeToCheckAltCoinPrice(currentTime) {
		mts.cryptoAltCoinPriceChan <- true
	}
}

func (mts *ManualTradingStrategy) InitService() {
	mts.cryptoHoldCoinPriceChan = make(chan bool)
	mts.cryptoAltCoinPriceChan = make(chan bool)

	go mts.startCheckHoldCoinPriceService(mts.cryptoHoldCoinPriceChan)
	go mts.startCheckAltCoinPriceService(mts.cryptoAltCoinPriceChan)
}

func (mts *ManualTradingStrategy) Shutdown() {
	close(mts.cryptoHoldCoinPriceChan)
	close(mts.cryptoAltCoinPriceChan)
}

func (ManualTradingStrategy) isTimeToCheckAltCoinPrice(time time.Time) bool {
	minute := time.Minute()
	var listMinutes []int = []int{15, 30, 45, 0}
	for _, a := range listMinutes {
		if a == minute {
			return true
		}
	}

	return false
}

func (ManualTradingStrategy) startCheckHoldCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		holdCoin := checkCryptoHoldCoinPrice()
		msg := ""
		if len(holdCoin) > 0 {
			msg = "List coin yang dihold:\n"
			haveNote := false
			for _, coin := range holdCoin {
				if coin.Note != "" {
					msg += crypto.GenerateMsg(coin)
					msg += "\n"
					haveNote = true
				}
			}
			if !haveNote {
				msg = ""
			}

			waitMasterCoinProcessed()
			if masterCoin != nil && msg != "" {
				msg += "untuk master coin:\n"
				msg += crypto.GenerateMsg(*masterCoin)
			}
		}

		sendNotif(msg)
	}
}

func (ManualTradingStrategy) startCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		altCoin := checkCryptoAltCoinPrice()
		msg := ""
		if len(altCoin) > 0 {
			if len(altCoin) > 5 {
				altCoin = sortAndGetTopFive(altCoin)
			}

			msg += "top gain coin:\n"
			for _, coin := range altCoin {
				msg += crypto.GenerateMsg(coin)
				msg += "\n"
			}

			if masterCoin != nil {
				msg += "untuk master coin:\n"
				msg += crypto.GenerateMsg(*masterCoin)
			}
		}

		sendNotif(msg)
	}
}

func sortAndGetTopFive(coins []models.BandResult) []models.BandResult {
	sort.Slice(coins, func(i, j int) bool { return coins[i].Weight > coins[j].Weight })

	return coins[0:5]
}
