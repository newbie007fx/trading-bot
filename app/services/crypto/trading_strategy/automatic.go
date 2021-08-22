package trading_strategy

import (
	"fmt"
	"sort"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var lastCheckCoin *[]models.BandResult
var timeToBuy bool

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
		if ats.isTimeToCheckAltCoinPrice(currentTime) || ats.isTimeToBuykAltCoinPrice(currentTime) {
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

func (AutomaticTradingStrategy) isTimeToBuykAltCoinPrice(currentTime time.Time) bool {
	minute := currentTime.Minute()
	var listMinutes []int = []int{5, 20, 35, 50}
	for _, a := range listMinutes {
		if a == minute {
			timeToBuy = true
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

					balance := services.GetBalance()
					msg += fmt.Sprintf("saldo saat ini: %f\n", balance)
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

func (ats *AutomaticTradingStrategy) startCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		altCoins := checkCryptoAltCoinPrice()
		msg := ""
		if len(altCoins) > 0 {

			coin := ats.sortAndGetHigest(altCoins)
			if coin == nil {
				continue
			}

			msg = "coin berikut telah dihold:\n"
			msg += generateMsg(*coin)
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

func (ats *AutomaticTradingStrategy) sortAndGetHigest(altCoins []models.BandResult) *models.BandResult {
	if timeToBuy {
		for i, _ := range altCoins {
			if lastCheckCoin == nil {
				return nil
			}
			for _, coin := range *lastCheckCoin {
				if altCoins[i].Symbol == coin.Symbol {
					altCoins[i].Weight += 0.5
				}
			}
		}

		lastCheckCoin = nil
		timeToBuy = false
	} else {
		lastCheckCoin = &altCoins
		return nil
	}

	sort.Slice(altCoins, func(i, j int) bool { return altCoins[i].Weight > altCoins[j].Weight })

	return &altCoins[0]
}
