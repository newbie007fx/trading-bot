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

var checkingTime time.Time
var holdCount int64 = 0
var masterTmp models.BandResult
var masterMidTmp models.BandResult

type AutomaticTradingStrategy struct {
	cryptoHoldCoinPriceChan chan bool
	cryptoAltCoinPriceChan  chan bool
	cryptoAltCoinDownChan   chan bool
	masterCoinChan          chan bool
}

func (ats *AutomaticTradingStrategy) Execute(currentTime time.Time) {
	checkingTime = currentTime

	condition := map[string]interface{}{"is_on_hold": true}
	holdCount = repositories.CountNotifConfig(&condition)

	if ats.isTimeToCheckAltCoinPrice(currentTime) || holdCount > 0 || currentTime.Minute()%2 == 1 {
		ats.masterCoinChan <- true
	}

	if holdCount > 0 {
		ats.cryptoHoldCoinPriceChan <- true
	}

	maxHold := crypto.GetMaxHold()
	if holdCount < maxHold {
		if ats.isTimeToCheckAltCoinPrice(currentTime) {
			ats.cryptoAltCoinPriceChan <- true
		} else {
			minuteLeft := checkingTime.Minute() % 15
			if (minuteLeft > 5 && minuteLeft <= 14) && checkingTime.Minute()%2 == 1 {
				waitMasterCoinProcessed()
				if checkMasterDown() && masterCoin.Direction == analysis.BAND_UP {
					ats.cryptoAltCoinDownChan <- true
				}
			}
		}
	}

}

func (ats *AutomaticTradingStrategy) InitService() {
	ats.cryptoHoldCoinPriceChan = make(chan bool)
	ats.cryptoAltCoinPriceChan = make(chan bool)
	ats.cryptoAltCoinDownChan = make(chan bool)
	ats.masterCoinChan = make(chan bool)

	go ats.startCheckHoldCoinPriceService(ats.cryptoHoldCoinPriceChan)
	go ats.startCheckAltCoinPriceService(ats.cryptoAltCoinPriceChan)
	go ats.startCheckAltCoinOnDownService(ats.cryptoAltCoinDownChan)
	go StartCheckMasterCoinPriceService(ats.masterCoinChan)
}

func (ats *AutomaticTradingStrategy) Shutdown() {
	close(ats.cryptoHoldCoinPriceChan)
	close(ats.cryptoAltCoinPriceChan)
	close(ats.cryptoAltCoinDownChan)
	close(ats.masterCoinChan)
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
		holdCoin := checkCryptoHoldCoinPrice(checkingTime)
		msg := ""
		if len(holdCoin) > 0 {

			tmpMsg := ""
			for _, coin := range holdCoin {
				waitMasterCoinProcessed()
				if masterCoin == nil {
					continue
				}

				holdCoinMid := crypto.CheckCoin(coin.Symbol, "1h", 0, GetEndDate(checkingTime))
				holdCoinLong := crypto.CheckCoin(coin.Symbol, "4h", 0, GetEndDate(checkingTime))
				isNeedTosell := analysis.IsNeedToSell(coin, *masterCoin, ats.isTimeToCheckAltCoinPrice(checkingTime), holdCoinMid.Trend, masterCoinLongInterval.Trend)
				if isNeedTosell || analysis.SpecialCondition(coin.Symbol, coin, *holdCoinMid, *holdCoinLong) {
					currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
					if err == nil {
						bands := coin.Bands
						lastBand := bands[len(bands)-1]
						err = crypto.ReleaseCoin(*currencyConfig, lastBand.Candle)
						if err != nil {
							tmpMsg = err.Error()
						} else {
							tmpMsg = fmt.Sprintf("coin berikut akan dijual %d:\n", GetEndDate(checkingTime))
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
		if skippedProcess() {
			log.Println("checking alt coin skipped")
			continue
		}
		masterTmp = *masterCoin
		masterMidTmp = *masterCoinLongInterval
		altCoins := checkCryptoAltCoinPrice(checkingTime)
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
						msg += fmt.Sprintf("coin berikut telah dihold on %d:\n", checkingTime.Unix())
						msg += crypto.GenerateMsg(coin)
						msg += fmt.Sprintf("weight: <b>%.2f</b>\n", coin.Weight)
						msg += "\n"
						msg += sendHoldMsg(&coin)
						msg += "\n"

						if checkIsOnLongIntervalChangePeriode() {
							break
						}

						holdCount++
					}
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

func (ats *AutomaticTradingStrategy) startCheckAltCoinOnDownService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		log.Println("executing alt check on master down")

		altCoins := []models.BandResult{}

		endDate := GetEndDate(checkingTime)

		responseChan := make(chan crypto.CandleResponse)

		limit := 100
		condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
		currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit)

		waitMasterCoinProcessed()
		if analysis.BearishEngulfing(masterCoin.Bands[len(masterCoin.Bands)-4:]) && masterCoin.Direction == analysis.BAND_DOWN {
			return
		}

		for _, data := range *currency_configs {
			request := crypto.CandleRequest{
				Symbol:       data.Symbol,
				EndDate:      endDate,
				Limit:        40,
				Resolution:   "15m",
				ResponseChan: responseChan,
			}

			result := crypto.MakeCryptoRequest(data, request)
			if result == nil || result.Direction == analysis.BAND_DOWN {
				continue
			}

			result.Weight = analysis.CalculateWeightOnDown(result)
			if result.Weight == 0 {
				continue
			}

			midInterval := crypto.CheckCoin(result.Symbol, "1h", 0, GetEndDate(checkingTime))
			if !analysis.IsIgnoredMasterDown(result, midInterval, masterCoin, checkingTime) {
				altCoins = append(altCoins, *result)
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
					msg = fmt.Sprintf("coin berikut telah dihold on %d:\n", checkingTime.Unix())
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

func (ats *AutomaticTradingStrategy) sortAndGetHigest(altCoins []models.BandResult) *[]models.BandResult {
	results := []models.BandResult{}
	timeInMilli := GetEndDate(checkingTime)
	for i := range altCoins {
		resultMid := crypto.CheckCoin(altCoins[i].Symbol, "1h", 0, timeInMilli)
		midWeight := getWeightCustomInterval(*resultMid, altCoins[i], "1h", nil)
		if midWeight == 0 {
			continue
		}
		altCoins[i].Weight += midWeight
		if altCoins[i].Weight > 1.7 {
			resultLong := crypto.CheckCoin(altCoins[i].Symbol, "4h", 0, timeInMilli)
			longWight := getWeightCustomInterval(*resultLong, altCoins[i], "4h", resultMid)
			if longWight == 0 {
				continue
			}
			altCoins[i].Weight += longWight
			if altCoins[i].Weight > 2.5 {
				results = append(results, altCoins[i])
			}
		}
	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool { return results[i].Weight > results[j].Weight })

		return &results
	}
	return nil
}

func getWeightCustomInterval(result, coin models.BandResult, interval string, previous *models.BandResult) float32 {
	weight := analysis.CalculateWeightLongInterval(&result, masterTmp.Trend)
	ignored := false

	if interval == "1h" {
		ignored = analysis.IsIgnoredMidInterval(&result, &coin)
	} else {
		ignored = analysis.IsIgnoredLongInterval(&result, &coin, previous, masterTmp.Trend, masterMidTmp.Trend)
	}

	if ignored || result.Direction == analysis.BAND_DOWN {
		return 0
	}

	return weight
}

func sendHoldMsg(result *models.BandResult) string {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(result.Symbol)
	if err != nil {
		return ""
	}
	return crypto.HoldCoinMessage(*currencyConfig, result)
}

func checkMasterDown() bool {
	if masterCoin.Trend == models.TREND_DOWN && masterCoinLongInterval.Trend == models.TREND_DOWN {
		return true
	}

	return false
}

func skippedProcess() bool {
	if masterCoin.Trend == models.TREND_DOWN && masterCoinLongInterval.Trend != models.TREND_DOWN {
		return true
	}

	if analysis.BearishEngulfing(masterCoin.Bands[len(masterCoin.Bands)-3:]) && masterCoin.Direction == analysis.BAND_DOWN {
		return true
	}

	return masterCoin.Direction == analysis.BAND_DOWN
}

func checkIsOnLongIntervalChangePeriode() bool {
	hours := []int{0, 4, 8, 12, 16, 20}
	hour := checkingTime.Hour()
	minute := checkingTime.Minute()
	for _, data := range hours {
		if data == hour && minute == 15 {
			return true
		}
	}

	return false
}
