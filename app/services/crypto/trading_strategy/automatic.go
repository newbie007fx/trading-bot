package trading_strategy

import (
	"fmt"
	"sort"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var checkingTime time.Time

type AutomaticTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
	cryptoAltCoinDownChan   chan bool
}

func (ats *AutomaticTradingStrategy) Execute(currentTime time.Time) {
	checkingTime = currentTime

	condition := map[string]interface{}{"is_on_hold": true}
	holdCount := repositories.CountNotifConfig(&condition)
	if holdCount > 0 {
		ats.cryptoHoldCoinPriceChan <- true
	}

	maxHold := crypto.GetMaxHold()
	if holdCount < maxHold {
		if ats.isTimeToCheckAltCoinPrice(currentTime) {
			ats.cryptoAltCoinPriceChan <- true
		}

		waitMasterCoinProcessed()
		if (masterCoin.Trend != models.TREND_UP && masterCoinLongInterval.Trend == models.TREND_DOWN) || checkMasterDown() {
			ats.cryptoAltCoinDownChan <- true
		}
	}

}

func (ats *AutomaticTradingStrategy) InitService() {
	ats.cryptoHoldCoinPriceChan = make(chan bool)
	ats.cryptoAltCoinPriceChan = make(chan bool)
	ats.cryptoAltCoinDownChan = make(chan bool)

	go ats.startCheckHoldCoinPriceService(ats.cryptoHoldCoinPriceChan)
	go ats.startCheckAltCoinPriceService(ats.cryptoAltCoinPriceChan)
	go ats.startCheckAltCoinOnDownService(ats.cryptoAltCoinDownChan)
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
		waitMasterCoinProcessed()
		holdCoin := checkCryptoHoldCoinPrice()
		msg := ""
		if len(holdCoin) > 0 {

			tmpMsg := ""
			for _, coin := range holdCoin {
				if analysis.IsNeedToSell(coin, *masterCoin, ats.isTimeToCheckAltCoinPrice(checkingTime), masterCoinLongInterval.Trend) {
					currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
					if err == nil {
						bands := coin.Bands
						lastBand := bands[len(bands)-1]
						err = crypto.ReleaseCoin(*currencyConfig, lastBand.Candle)
						if err != nil {
							tmpMsg = err.Error()
						} else {
							tmpMsg = "coin berikut akan dijual:\n"
							tmpMsg += crypto.GenerateMsg(coin)
							tmpMsg += "\n"
							tmpMsg += crypto.HoldCoinMessage(*currencyConfig, &coin)
							tmpMsg += "\n"
							tmpMsg += "alasan dijual: " + analysis.GetSellReason() + "\n\n"

							balance := crypto.GetBalanceFromConfig()
							tmpMsg += fmt.Sprintf("saldo saat ini: %f\n", balance)
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

		crypto.SendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) startCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		waitMasterCoinProcessed()
		if masterCoin.Trend != models.TREND_UP && masterCoinLongInterval.Trend == models.TREND_DOWN {
			continue
		}
		altCoins := checkCryptoAltCoinPrice(&checkingTime)
		msg := ""
		if len(altCoins) > 0 {

			coin := ats.sortAndGetHigest(altCoins)
			if coin == nil {
				continue
			}

			currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
			if err == nil {
				bands := coin.Bands
				lastBand := bands[len(bands)-1]
				err = crypto.HoldCoin(*currencyConfig, lastBand.Candle)
				if err != nil {
					msg = err.Error()
				} else {
					msg = "coin berikut telah dihold:\n"
					msg += crypto.GenerateMsg(*coin)
					msg += fmt.Sprintf("weight: <b>%.2f</b>\n", coin.Weight)
					msg += "\n"
					msg += sendHoldMsg(coin)
					msg += "\n"

					if masterCoin != nil {
						msg += "untuk master coin:\n"
						msg += crypto.GenerateMsg(*masterCoin)
					}
				}
			}
		}

		crypto.SendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) startCheckAltCoinOnDownService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		altCoins := []models.BandResult{}

		endDate := GetEndDate(&checkingTime)

		responseChan := make(chan crypto.CandleResponse)

		limit := 120
		condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
		currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit)

		waitMasterCoinProcessed()
		for _, data := range *currency_configs {
			request := crypto.CandleRequest{
				Symbol:       data.Symbol,
				EndDate:      endDate,
				Limit:        40,
				Resolution:   "15m",
				ResponseChan: responseChan,
			}

			result := crypto.MakeCryptoRequest(data, request)
			if result == nil || result.Direction == analysis.BAND_DOWN || !(result.Position == models.BELOW_LOWER || result.Position == models.BELOW_SMA) {
				continue
			}

			if result.AllTrend.SecondTrend != models.TREND_UP {
				result.Weight = analysis.CalculateWeightOnDown(result)
				if result.Weight != 0 {
					altCoins = append(altCoins, *result)
				}
			}
		}
		msg := ""
		if len(altCoins) > 0 {

			sort.Slice(altCoins, func(i, j int) bool { return altCoins[i].Weight > altCoins[j].Weight })
			coin := altCoins[0]
			currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
			if err == nil {
				bands := coin.Bands
				lastBand := bands[len(bands)-1]
				err = crypto.HoldCoin(*currencyConfig, lastBand.Candle)
				if err != nil {
					msg = err.Error()
				} else {
					msg = "coin berikut telah dihold:\n"
					msg += crypto.GenerateMsg(coin)
					msg += fmt.Sprintf("weight: <b>%.2f</b>\n", coin.Weight)
					msg += "\n"
					msg += sendHoldMsg(&coin)
					msg += "\n"

					if masterCoin != nil {
						msg += "untuk master coin:\n"
						msg += crypto.GenerateMsg(*masterCoin)
					}
				}
			}
		}

		crypto.SendNotif(msg)
	}
}

func (ats *AutomaticTradingStrategy) sortAndGetHigest(altCoins []models.BandResult) *models.BandResult {
	results := []models.BandResult{}
	timeInMilli := GetEndDate(&checkingTime)
	for i := range altCoins {
		waitMasterCoinProcessed()
		altCoins[i].Weight += crypto.GetOnLongIntervalWeight(altCoins[i], *masterCoin, timeInMilli)
		if altCoins[i].Weight > 2.8 {
			results = append(results, altCoins[i])
		}
	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool { return results[i].Weight > results[j].Weight })

		return &results[0]
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

func checkMasterDown() bool {
	if masterCoin.Direction == analysis.BAND_DOWN {
		return false
	}

	masterLastBand := masterCoin.Bands[len(masterCoin.Bands)-2]
	masterSecondLastBand := masterCoin.Bands[len(masterCoin.Bands)-2]
	if masterSecondLastBand.Candle.Open > masterSecondLastBand.Candle.Close {
		secondLastBandPriceChanges := (masterSecondLastBand.Candle.Open - masterSecondLastBand.Candle.Close) / masterSecondLastBand.Candle.Open * 100
		if secondLastBandPriceChanges >= 0.55 {
			return masterCoin.PriceChanges > 0.26
		}
	}

	lastBandPercentChanges := (masterLastBand.Candle.Close - masterLastBand.Candle.Open) / masterLastBand.Candle.Open * 100
	return lastBandPercentChanges > 3.3 && checkingTime.Minute()/15 >= 10
}
