package analysis

import (
	"log"
	"telebot-trading/app/models"
	"time"
)

var reason string = ""

func IsNeedToSell(currencyConfig *models.CurrencyNotifConfig, result models.BandResult, requestTime time.Time, resultMid *models.BandResult) bool {
	reason = ""
	lastBand := result.Bands[len(result.Bands)-1]
	changes := result.CurrentPrice - currencyConfig.HoldPrice
	changesInPercent := changes / currencyConfig.HoldPrice * 100
	isCandleComplete := checkIsCandleComplete(requestTime, 15)

	if sellWhenDoubleUpTargetProfit(*currencyConfig, result, *resultMid, changesInPercent) {
		reason = "sell with double up target profit"
		return true
	}

	whenDown := result.AllTrend.Trend == models.TREND_DOWN || (lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA))
	if SellPattern(&result) && isCandleComplete && (changesInPercent > 1 || result.AllTrend.SecondTrend == models.TREND_DOWN) && !whenDown {
		midLastBand := resultMid.Bands[len(result.Bands)-1]
		midSecondLastBand := resultMid.Bands[len(result.Bands)-2]
		previousMidBandValid := false
		if midSecondLastBand.Candle.Close > midSecondLastBand.Candle.Open {
			halfMidSecondLastBand := (midSecondLastBand.Candle.Close-midSecondLastBand.Candle.Open)/2 + midSecondLastBand.Candle.Open
			previousMidBandValid = midLastBand.Candle.Close > halfMidSecondLastBand
		}
		lastbandCrossLower := lastBand.Candle.Low < float32(lastBand.Lower) && lastBand.Candle.Close > float32(lastBand.Lower)
		if changesInPercent < 3 && (resultMid.Direction != BAND_DOWN || previousMidBandValid || lastbandCrossLower) {

		} else {
			reason = "sell with criteria bearish engulfing"
			return true
		}
	}

	if currencyConfig.HoldPrice > result.CurrentPrice {
		if skipSell(*resultMid) {
			log.Println("skip sell gans")
			return false
		}

		if sellOnDown(result, currencyConfig, lastBand) {
			return true
		}
	} else {
		if changesInPercent > 3 && result.Direction == BAND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
			if changesInPercent <= 3.5 || isCandleComplete {
				reason = "sell up with criteria change > 3 and short trend = down"
				return true
			}
		}

		if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend != models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP && changesInPercent > 3 && isCandleComplete && result.Direction == BAND_DOWN {
			reason = "sell up with criteria trend not up after up"
			return true
		}

		if aboveUpperAndHeigt3xAvg(result) && changesInPercent > 5 && down25PercentFromHight(result, changesInPercent, currencyConfig.HoldPrice, 10) {
			reason = "above upper and heigt 3x avg"
			return true
		}

		if previousBandUpThenDown(result, changesInPercent, currencyConfig.HoldPrice) {
			reason = "previous band up then down"
			return true
		}

		if shortTrendOnPreviousBandNotUpAndDown25PercentFromHight(result, changesInPercent, currencyConfig.HoldPrice) {
			reason = "short trend on previous band not up and last band down 25 percent from hight"
			return true
		}

		if resultMid.Position == models.ABOVE_UPPER && !isHasCrossUpper(result.Bands[:len(result.Bands)-1], false) {
			if shortTrendOnPreviousBandNotUpAndDown25PercentFromHight(*resultMid, changesInPercent, currencyConfig.HoldPrice) {
				reason = "mid interval short trend on previous band not up and last band down 25 percent from hight"
				return true
			}
		}

		if firstCrossUpper(result, *resultMid, changesInPercent, currencyConfig) {
			reason = "first cross upper"
			return true
		}

		if changesInPercent > 3 && changesInPercent < 3.5 && aboveUpperAndMidIntervalCrossSMA(result, *resultMid) {
			reason = "above upper and mid interval cross sma"
			return true
		}

		midLastBand := resultMid.Bands[len(resultMid.Bands)-1]
		if changesInPercent > 3 && changesInPercent < 3.5 && midLastBand.Candle.Low > float32(midLastBand.Upper) {
			reason = "mid lastband low above upper"
			return true
		}

		if isHasOpenCloseAboveUpper(result.Bands[len(result.Bands)-3:]) && changesInPercent > 2.6 && changesInPercent < 3 {
			reason = "has Open close above upper"
			return true
		}

		if sellOnUp(result, currencyConfig, resultMid.AllTrend.Trend, isCandleComplete) {
			return true
		}
	}

	return isHoldedMoreThanDurationThreshold(currencyConfig, result, isCandleComplete)
}

func sellOnUp(result models.BandResult, currencyConfig *models.CurrencyNotifConfig, coinLongTrend int8, isCandleComplete bool) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	changes := result.CurrentPrice - currencyConfig.HoldPrice
	changesInPercent := changes / currencyConfig.HoldPrice * 100

	highest := getHigestPrice(result.Bands)
	highestChangePercent := changes / (highest - currencyConfig.HoldPrice) * 100

	highestHight := GetHigestHightPrice(result.Bands)
	highestHightChangePercent := changes / (highestHight - currencyConfig.HoldPrice) * 100

	lastFiveData := result.Bands[len(result.Bands)-5 : len(result.Bands)]

	lastBandPercentChangesDown := (lastBand.Candle.Open - lastBand.Candle.Close) / lastBand.Candle.Close * 100
	lastBandPercentChanges := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100
	lastHightChangePercent := (lastBand.Candle.Close - lastBand.Candle.Open) / (lastBand.Candle.Hight - lastBand.Candle.Open) * 100
	changeFivePercentAndDownFromHight := (changesInPercent > 5 && lastBandPercentChanges > 5 && lastHightChangePercent <= 55 && isTimeBelowTenMinute())
	changeFivePercentAndDownFromPreviousHight := changesInPercent > 5 && highestChangePercent <= 65 && lastBandPercentChangesDown > 2
	specialTolerance := (changesInPercent > 10 && highestHightChangePercent <= 65) || changeFivePercentAndDownFromHight || changeFivePercentAndDownFromPreviousHight
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
	secondLastBand := result.Bands[len(result.Bands)-2]
	if changesInPercent >= 3 && result.Direction == BAND_DOWN {
		if result.Position == models.BELOW_LOWER && (changesInPercent > 3.3 || secondLastBand.Candle.Close < float32(secondLastBand.Lower)) {
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

func sellWhenDoubleUpTargetProfit(config models.CurrencyNotifConfig, result models.BandResult, midResult models.BandResult, changeInPercent float32) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if changeInPercent >= 3 && changeInPercent < 3.847 && result.Direction == BAND_DOWN {
		timeFirstAboveMargin := getTimeFirstAboveMargin(config.HoldPrice, midResult.Bands, config.HoldedAt)
		if lastBand.Candle.CloseTime > timeFirstAboveMargin {
			fourHourInSecond := 60 * 60 * 4
			diff := lastBand.Candle.CloseTime - timeFirstAboveMargin
			return diff/int64(fourHourInSecond) >= 1
		}
	}

	return false
}

func getTimeFirstAboveMargin(holdPrice float32, bands []models.Band, openTime int64) int64 {
	index := getIndexMoreThanLimitOpenTime(bands, openTime)
	for _, band := range bands[index:] {
		percent := (band.Candle.Hight - holdPrice) / holdPrice * 100
		if percent > 3 {
			return band.Candle.CloseTime
		}
	}

	return 0
}

func getIndexMoreThanLimitOpenTime(bands []models.Band, openTime int64) int {
	for i, band := range bands {
		if band.Candle.CloseTime > openTime {
			return i
		}
	}
	return 0
}

func skipSell(resultMid models.BandResult) bool {
	lastBand := resultMid.Bands[len(resultMid.Bands)-1]
	if resultMid.AllTrend.SecondTrend == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.SMA) {
		if countAboveSMA(resultMid.Bands[len(resultMid.Bands)/2:]) >= len(resultMid.Bands)/2 && countDownBand(resultMid.Bands[len(resultMid.Bands)-3:]) < 3 {
			return true
		}
	}

	return false
}

func aboveUpperAndMidIntervalCrossSMA(resultShort, resultMid models.BandResult) bool {
	midLastBand := resultMid.Bands[len(resultMid.Bands)-1]
	if resultShort.Position == models.ABOVE_UPPER && resultMid.Position == models.ABOVE_SMA && midLastBand.Candle.Open < float32(midLastBand.SMA) {
		return countBelowSMA(resultMid.Bands[len(resultMid.Bands)-6:len(resultMid.Bands)-1], false) == 5 && !isHasCrossLower(resultMid.Bands[len(resultMid.Bands)-6:len(resultMid.Bands)-1], false)
	}

	return false
}

func countAboveSMA(bands []models.Band) int {
	var count int = 0
	for _, band := range bands {
		if band.Candle.Open > float32(band.SMA) && band.Candle.Close > float32(band.SMA) {
			count++
		}
	}
	return count
}

func countAboveUpper(bands []models.Band) int {
	var count int = 0
	for _, band := range bands {
		if band.Candle.Open > float32(band.Upper) {
			count++
		}
	}
	return count
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

func GetHigestHightPrice(bands []models.Band) float32 {
	var highest float32 = 0
	for _, band := range bands {
		if highest < band.Candle.Hight {
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

func GetLowestLowPrice(bands []models.Band) float32 {
	var lowest float32 = bands[0].Candle.Low
	for _, band := range bands {
		if lowest > band.Candle.Low {
			lowest = band.Candle.Low
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

func SellPattern(bandResult *models.BandResult) bool {
	lastBand := bandResult.Bands[len(bandResult.Bands)-1]
	secondLastBand := bandResult.Bands[len(bandResult.Bands)-2]
	if BearishEngulfing(bandResult.Bands[len(bandResult.Bands)-3:]) {
		return lastBand.Candle.Close < lastBand.Candle.Open && secondLastBand.Candle.Close < secondLastBand.Candle.Open
	}

	return false
}

func SpecialCondition(currencyConfig *models.CurrencyNotifConfig, symbol string, shortInterval, midInterval, longInterval models.BandResult) bool {
	lastBand := shortInterval.Bands[len(shortInterval.Bands)-1]
	changes := lastBand.Candle.Close - currencyConfig.HoldPrice
	changesInPercent := changes / currencyConfig.HoldPrice * 100

	if isLastBandCrossUpperAndPreviousBandNot(shortInterval.Bands) && changesInPercent >= 3 {
		if isLastBandCrossUpperAndPreviousBandNot(midInterval.Bands) {
			if isLastBandCrossUpperAndPreviousBandNot(longInterval.Bands) {
				reason = "last band cross upper and previous band not"
				return true
			}
		}
	}

	if aboveUpperAndOtherIntervalAboveSMA(shortInterval, midInterval, longInterval, changesInPercent, currencyConfig.HoldPrice) && changesInPercent >= 3 {
		reason = "above upper and other interval above sma"
		return true
	}

	midLastBand := midInterval.Bands[len(midInterval.Bands)-1]
	longLastBand := longInterval.Bands[len(longInterval.Bands)-1]
	if currencyConfig.HoldPrice < lastBand.Candle.Close && changesInPercent > 3 {
		if lastBand.Candle.Open > float32(lastBand.Upper) && lastBand.Candle.Close > float32(lastBand.Upper) {
			if midLastBand.Candle.Open < float32(midLastBand.Upper) && midLastBand.Candle.Close > float32(midLastBand.Upper) {
				if longLastBand.Candle.Open < float32(longLastBand.Upper) && longLastBand.Candle.Close > float32(longLastBand.Upper) {
					reason = "open close above upper, mid cross upper, long cross upper"
					return true
				}
			}
		}
	}

	if longInterval.AllTrend.FirstTrend == models.TREND_DOWN && longInterval.AllTrend.SecondTrend == models.TREND_DOWN && countBelowSMA(longInterval.Bands, true) > len(longInterval.Bands)/2 {
		if midInterval.AllTrend.FirstTrend != models.TREND_UP && midInterval.AllTrend.SecondTrend == models.TREND_UP && isHasCrossLower(midInterval.Bands, false) {
			if shortInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands) == 1 && changesInPercent >= 2.5 && changesInPercent < 3 {
				reason = "first up on trend down"
				return true
			}
		}
	}

	if longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
		hightIndex := getHighestIndex(midInterval.Bands[len(midInterval.Bands)/2:]) + len(midInterval.Bands)/2
		if hightIndex >= len(midInterval.Bands)/2-3 && midInterval.Bands[hightIndex].Candle.Hight > float32(midInterval.Bands[hightIndex].Upper) {
			if countBelowSMA(midInterval.Bands[hightIndex:], false) > 0 && midLastBand.Candle.Hight >= midInterval.Bands[hightIndex].Candle.Close {
				if changesInPercent > 3 && changesInPercent < 3.5 {
					reason = "long interval down, reversal from sma"
					return true
				}
			}
		}
	}

	if longInterval.AllTrend.SecondTrend != models.TREND_UP && longInterval.AllTrend.ShortTrend != models.TREND_UP {
		midHigestBand := getHighestIndex(midInterval.Bands)
		midSecondHight := getHigestIndexSecond(midInterval.Bands)
		if midSecondHight == len(midInterval.Bands)-1 && midSecondHight-midHigestBand > 5 && !diffInThresholdFromHigest(midInterval.Bands[midHigestBand], midInterval.Bands[midSecondHight]) {
			shortHigestBand := getHighestIndex(shortInterval.Bands)
			if shortHigestBand == len(shortInterval.Bands)-1 && changesInPercent > 3 && changesInPercent < 3.2 {
				reason = "long interval down, got new hight"
				return true
			}
		}
	}

	if longInterval.AllTrend.FirstTrend == models.TREND_DOWN && longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
		if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-1:], false) || isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-1:], true) {
			if lastBand.Candle.Close > float32(lastBand.Upper) && changesInPercent > 2.5 && changesInPercent < 3 {
				reason = "long interval down, mid cross sma and short cross upper"
				return true
			}
		}
	}

	return false
}

func diffInThresholdFromHigest(higest, close models.Band) bool {
	if higest.Candle.Close < close.Candle.Hight {
		return false
	}

	percent := (higest.Candle.Close - close.Candle.Hight) / close.Candle.Hight * 100
	return percent > 0.5
}

func aboveUpperAndOtherIntervalAboveSMA(shortInterval, midInterval, longInterval models.BandResult, changeInpercent, holdPrice float32) bool {
	if shortInterval.Position == models.ABOVE_UPPER && midInterval.Position == models.ABOVE_SMA && longInterval.Position == models.ABOVE_SMA {
		return down25PercentFromHight(shortInterval, changeInpercent, holdPrice, 5)
	}

	return false
}

func aboveUpperAndHeigt3xAvg(result models.BandResult) bool {
	if result.Position != models.ABOVE_UPPER || result.Direction != BAND_UP {
		return false
	}

	var total float32
	mid := len(result.Bands) / 2
	bands := result.Bands[mid : len(result.Bands)-1]
	for _, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			total += band.Candle.Close - band.Candle.Open
		} else {
			total += band.Candle.Open - band.Candle.Close
		}
	}
	average := total / float32(len(bands))
	lastBand := result.Bands[len(result.Bands)-1]
	lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open

	return average*3 < lastBandHeight
}

func previousBandUpThenDown(result models.BandResult, changeInPercent float32, holdPrice float32) bool {
	heightIndex := getIndexHigestCrossUpper(result.Bands)
	if heightIndex == -1 || heightIndex > len(result.Bands)-5 || result.Direction != BAND_DOWN {
		return false
	}

	heightBand := result.Bands[heightIndex]
	percentFromHeight := (heightBand.Candle.Close - holdPrice) / holdPrice * 100
	if percentFromHeight > 4 && percentFromHeight < 7 {
		return changeInPercent/percentFromHeight*100 < 76 && changeInPercent >= 3
	}
	return false
}

func shortTrendOnPreviousBandNotUpAndDown25PercentFromHight(result models.BandResult, changeInPercent float32, holdPrice float32) bool {
	shortTrendPreviousBand := CalculateTrendShort(result.Bands[len(result.Bands)-6 : len(result.Bands)-1])
	shortTrendPreviousBand2 := CalculateTrendShort(result.Bands[len(result.Bands)-7 : len(result.Bands)-2])
	secondLastBand := result.Bands[len(result.Bands)-2]
	if (shortTrendPreviousBand2 != models.TREND_UP || shortTrendPreviousBand != models.TREND_UP || secondLastBand.Candle.Close < secondLastBand.Candle.Open) && result.AllTrend.ShortTrend == models.TREND_UP {
		return down25PercentFromHight(result, changeInPercent, holdPrice, 7) && changeInPercent >= 3
	}

	return false
}

func down25PercentFromHight(result models.BandResult, changeInPercent float32, holdPrice float32, maxFromHight int) bool {
	heightPrice := result.Bands[len(result.Bands)-1].Candle.Hight
	percentFromHeight := (heightPrice - holdPrice) / holdPrice * 100
	if percentFromHeight < float32(maxFromHight) {
		percent := changeInPercent / percentFromHeight * 100
		return percent < 82
	}

	return false
}

func checkIsCandleComplete(requestTime time.Time, intervalMinute int) bool {
	minute := requestTime.Minute()
	return minute%intervalMinute == 0
}

func firstCrossUpper(shortInterval, midInterval models.BandResult, changeInPercent float32, currencyNotifConfig *models.CurrencyNotifConfig) bool {
	midLastBand := midInterval.Bands[len(midInterval.Bands)-1]
	shortLastBand := shortInterval.Bands[len(shortInterval.Bands)-1]

	if midInterval.Position == models.ABOVE_UPPER && shortInterval.Position == models.ABOVE_UPPER {
		if !isHasCrossUpper(midInterval.Bands[:len(midInterval.Bands)-1], true) && isHasCrossLower(midInterval.Bands[:len(midInterval.Bands)-1], false) {
			return changeInPercent > 3 && changeInPercent < 3.5
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && ((midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Close > float32(midLastBand.SMA)) || (midLastBand.Candle.Open < float32(midLastBand.Upper) && midLastBand.Candle.Close > float32(midLastBand.Upper))) {
		if !isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-6:len(shortInterval.Bands)-1], true) || isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-6:len(midInterval.Bands)-1], true) {
			return changeInPercent > 3 && changeInPercent < 3.5
		}
	}

	if shortInterval.Position == models.ABOVE_SMA && isHasCrossSMA(shortInterval.Bands[len(shortInterval.Bands)-2:], false) {
		if isHasBelowLower(midInterval.Bands[len(midInterval.Bands)-3:]) {
			return changeInPercent > 3 && changeInPercent < 3.5
		}
	}

	percentFromHigest := (shortLastBand.Candle.Hight - currencyNotifConfig.HoldPrice) / currencyNotifConfig.HoldPrice * 100
	if shortLastBand.Candle.Open > float32(shortLastBand.Upper) && shortLastBand.Candle.Close > float32(shortLastBand.Upper) {
		if isHasCrossSMA([]models.Band{midLastBand}, false) || isHasCrossUpper([]models.Band{midLastBand}, true) {
			return percentFromHigest > 3 && changeInPercent > 2.5
		}
	}

	return false
}

func isHasBelowLower(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open < float32(band.Lower) && band.Candle.Close < float32(band.Lower) {
			return true
		}
	}

	return false
}
