package analysis

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"time"
)

var reason string = ""

func IsNeedToSell(result models.BandResult, masterCoin models.BandResult, isCandleComplete bool, coinLongTrend, masterCoinLongTrend int8) bool {
	reason = ""
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(result.Symbol)
	if err != nil {
		log.Panicln(err.Error())
	}

	changes := result.CurrentPrice - currencyConfig.HoldPrice
	changesInPercent := changes / currencyConfig.HoldPrice * 100

	if isCandleComplete && masterCoin.Direction == BAND_DOWN {
		var masterDown, resultDown, safe = false, false, false
		for i := len(result.Bands) - 1; i >= len(result.Bands)-2; i-- {
			masterDown = masterCoin.Bands[i].Candle.Open > masterCoin.Bands[i].Candle.Close
			resultDown = result.Bands[i].Candle.Open > result.Bands[i].Candle.Close
			if !(masterDown && resultDown) {
				safe = true
				break
			}
		}

		if !safe && result.AllTrend.SecondTrend == models.TREND_DOWN {
			reason = "sell with criteria 0"
			return true
		}
	}

	if SellPattern(&result) && changesInPercent > 1 && isCandleComplete {
		reason = "sell up with criteria x1"
		return true
	}

	lastBand := result.Bands[len(result.Bands)-1]
	if currencyConfig.HoldPrice > result.CurrentPrice {
		if sellOnDown(result, currencyConfig, lastBand) {
			return true
		}
	} else {
		if changesInPercent > 3 && result.Direction == BAND_DOWN && masterCoin.Trend == models.TREND_DOWN && masterCoinLongTrend != models.TREND_UP {
			reason = "sell up with criteria x0"
			return true
		}

		if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend != models.TREND_UP && changesInPercent > 2.5 && isCandleComplete && result.Direction == BAND_DOWN {
			reason = "sell up with criteria x2"
			return true
		}

		if sellOnUp(result, currencyConfig, coinLongTrend, isCandleComplete, masterCoin.Trend, masterCoinLongTrend) {
			return true
		}
	}

	return isHoldedMoreThanDurationThreshold(currencyConfig, result, isCandleComplete)
}

func sellOnUp(result models.BandResult, currencyConfig *models.CurrencyNotifConfig, coinLongTrend int8, isCandleComplete bool, masterCoinTrend, masterCoinLongTrend int8) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	changes := result.CurrentPrice - currencyConfig.HoldPrice
	changesInPercent := changes / currencyConfig.HoldPrice * 100

	highest := getHigestPrice(result.Bands)
	highestChangePercent := changes / (highest - currencyConfig.HoldPrice) * 100

	highestHight := getHigestHightPrice(result.Bands)
	highestHightChangePercent := changes / (highestHight - currencyConfig.HoldPrice) * 100

	lastFiveData := result.Bands[len(result.Bands)-5 : len(result.Bands)]

	if checkOnTrendDown(result, coinLongTrend, masterCoinTrend, masterCoinLongTrend, changesInPercent, isCandleComplete) {
		reason = "sell with criteria y1"
		return true
	}

	lastBandPercentChanges := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100
	lastHightChangePercent := (lastBand.Candle.Close - lastBand.Candle.Open) / (lastBand.Candle.Hight - lastBand.Candle.Open) * 100
	specialTolerance := (changesInPercent > 10 && highestHightChangePercent <= 65) || (changesInPercent > 5 && lastBandPercentChanges > 5 && lastHightChangePercent <= 55 && isTimeBelowTenMinute())
	if !specialTolerance {
		if changesInPercent > 3.5 && !isCandleComplete && highestChangePercent > 55 && countDownCandleFromHighest(result.Bands) < 4 {
			return false
		}
	}

	condition := highestChangePercent <= 65 && changesInPercent >= 3
	if specialTolerance || (condition && CalculateTrendShort(lastFiveData) == models.TREND_DOWN && result.Direction == BAND_DOWN) {

		secondLastBand := result.Bands[len(result.Bands)-2]
		if result.Position == models.BELOW_LOWER {
			if lastBand.Candle.Open > float32(lastBand.Lower) && float32(lastBand.Lower) > result.CurrentPrice {
				changesOnLower := result.CurrentPrice - float32(lastBand.Lower)
				changesOnLowerPercent := changesOnLower / float32(lastBand.Lower) * 100
				if changesOnLowerPercent >= 3 {
					reason = "sell on up with criteria 1"
					return true
				}

				if changesInPercent > 3.7 {
					return false
				}
			}

			if secondLastBand.Candle.Close > float32(secondLastBand.Lower) || secondLastBand.Candle.Open > float32(secondLastBand.Lower) {
				if lastBand.Candle.Close < float32(lastBand.Lower) && lastBand.Candle.Open < float32(lastBand.Lower) && lastBand.Candle.Open > lastBand.Candle.Close {
					reason = "sell on up with criteria 2"
					return true
				}
			}
		} else if result.Position == models.BELOW_SMA {
			changesFromLower := result.CurrentPrice - float32(lastBand.Lower)
			changesFromLowerPercent := changesFromLower / result.CurrentPrice * 100
			if changesFromLowerPercent <= 1 {
				return false
			}

			if lastBand.Candle.Open > float32(lastBand.SMA) && float32(lastBand.SMA) > result.CurrentPrice {
				changesOnSMA := result.CurrentPrice - float32(lastBand.SMA)
				changesOnSMAPercent := changesOnSMA / float32(lastBand.SMA) * 100
				if changesOnSMAPercent >= 3 {
					reason = "sell on up with criteria 3"
					return true
				}

				if changesInPercent > 3.7 {
					return false
				}
			}

			if secondLastBand.Candle.Close > float32(secondLastBand.SMA) || secondLastBand.Candle.Open > float32(secondLastBand.SMA) {
				if lastBand.Candle.Close < float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Open > lastBand.Candle.Close {
					reason = "sell on up with criteria 4"
					return true
				}
			}
		} else if result.Position == models.ABOVE_SMA {
			changesFromSMA := result.CurrentPrice - float32(lastBand.SMA)
			changesFromSMAPercent := changesFromSMA / result.CurrentPrice * 100
			if changesFromSMAPercent <= 1 {
				return false
			}

			if lastBand.Candle.Open > float32(lastBand.Upper) && float32(lastBand.Upper) > result.CurrentPrice {
				changesOnUpper := result.CurrentPrice - float32(lastBand.Upper)
				changesOnUpperPercent := changesOnUpper / float32(lastBand.Upper) * 100
				if changesOnUpperPercent >= 3 {
					reason = "sell on up with criteria 5"
					return true
				}

				if changesInPercent > 3.7 {
					return false
				}
			}

			if secondLastBand.Candle.Close > float32(secondLastBand.Upper) || secondLastBand.Candle.Open > float32(secondLastBand.Upper) {
				if lastBand.Candle.Close < float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) && lastBand.Candle.Open > lastBand.Candle.Close {
					reason = "sell on up with criteria 6"
					return true
				}
			}
		} else if result.Position == models.ABOVE_UPPER {
			changesFromUpper := result.CurrentPrice - float32(lastBand.Upper)
			changesFromUpperPercent := changesFromUpper / result.CurrentPrice * 100
			if changesFromUpperPercent <= 1 {
				return false
			}
		}

		reason = "sell on up with criteria 7"
		return true

	}

	return false
}

func sellOnDown(result models.BandResult, currencyConfig *models.CurrencyNotifConfig, lastBand models.Band) bool {
	changes := currencyConfig.HoldPrice - result.CurrentPrice
	changesInPercent := changes / currencyConfig.HoldPrice * 100
	if changesInPercent >= 3 && result.Direction == BAND_DOWN {
		if result.Position == models.BELOW_LOWER {
			reason = "sell on down with criteria 1"
			return true
		} else if result.Position == models.BELOW_SMA {
			changesFromLower := currencyConfig.HoldPrice - float32(lastBand.Lower)
			changesFromLowerPercent := changesFromLower / currencyConfig.HoldPrice * 100
			if changesFromLowerPercent > 4 {
				reason = "sell on down with criteria 2"
				return true
			}
		} else if result.Position == models.ABOVE_SMA {
			changesFromSMA := currencyConfig.HoldPrice - float32(lastBand.SMA)
			changesFromSMAPercent := changesFromSMA / currencyConfig.HoldPrice * 100
			if changesFromSMAPercent > 4 {
				reason = "sell on down with criteria 3"
				return true
			}
		} else if result.Position == models.ABOVE_UPPER {
			changesFromUpper := currencyConfig.HoldPrice - float32(lastBand.Upper)
			changesFromUpperPercent := changesFromUpper / currencyConfig.HoldPrice * 100
			if changesFromUpperPercent > 4 {
				reason = "sell on down with criteria 4"
				return true
			}
		}
	} else if changesInPercent > 4 {
		reason = "sell on down with criteria 5"
		return true
	}

	return false
}

var highestIndex int = 0

func getHigestPrice(bands []models.Band) float32 {
	var highest float32 = 0
	for i, band := range bands {
		if highest < band.Candle.Close {
			highest = band.Candle.Close
			highestIndex = i
		}
	}

	return highest
}

func getHigestHightPrice(bands []models.Band) float32 {
	var highest float32 = 0
	for _, band := range bands {
		if highest < band.Candle.Close {
			highest = band.Candle.Hight
		}
	}

	return highest
}

func getLowestPrice(bands []models.Band) float32 {
	var lowest float32 = bands[0].Candle.Close
	for _, band := range bands {
		if lowest > band.Candle.Close {
			lowest = band.Candle.Close
		}
	}

	return lowest
}

func countDownCandleFromHighest(bands []models.Band) int {
	count := 0
	for i := highestIndex; i < len(bands); i++ {
		if bands[i].Candle.Close < bands[i].Candle.Open {
			count++
		}
	}
	return count
}

func GetSellReason() string {
	return reason
}

func checkOnTrendDown(result models.BandResult, coinLongTrend, masterCoinTrend, masterCoinLongIntervalTrend int8, priceChange float32, isCandleComplete bool) bool {
	if (masterCoinTrend != models.TREND_UP || coinLongTrend == models.TREND_DOWN) && masterCoinLongIntervalTrend == models.TREND_DOWN {
		if result.Direction == BAND_DOWN && result.AllTrend.SecondTrend != models.TREND_UP && isCandleComplete {
			lastBand := result.Bands[len(result.Bands)-1]
			lastBandOnUpper := lastBand.Candle.Low <= float32(lastBand.Upper) && lastBand.Candle.Hight >= float32(lastBand.Upper)
			if lastBandOnUpper {
				return true
			}
		}
	}

	return false
}

func isTimeBelowTenMinute() bool {
	currentTime := time.Now()

	return currentTime.Minute()%15 <= 10
}

func isHoldedMoreThanDurationThreshold(config *models.CurrencyNotifConfig, result models.BandResult, isCandleComplete bool) bool {
	currentTime := time.Now()
	durationOnOnePeriode := int64(20 * 4 * 60 * 60)
	maxThreshold := config.HoldedAt + durationOnOnePeriode

	if currentTime.Unix() > maxThreshold {
		if result.AllTrend.SecondTrend != models.TREND_UP && result.Direction == BAND_DOWN && isCandleComplete {
			reason = "sell after hold more than threshold"
			return true
		}
	}

	return false
}
