package analysis

import (
	"fmt"
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""

func IsIgnored(result, masterCoin *models.BandResult, requestTime time.Time) bool {
	if isInAboveUpperBandAndDownTrend(result) {
		ignoredReason = "isInAboveUpperBandAndDownTrend"
		return true
	}

	if lastBandHeadDoubleBody(result) {
		ignoredReason = "lastBandHeadDoubleBody"
		return true
	}

	if isContaineBearishEngulfing(result) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if isPosititionBellowUpperMarginBellowThreshold(result) {
		ignoredReason = "isPosititionBellowUpperMarginBellowThreshold"
		return true
	}

	if whenHeightTripleAverage(result) {
		ignoredReason = "whenHeightTripleAverage"
		return true
	}

	if lastFourCandleNotUpTrend(result.Bands) && !reversalFromLower(*result) {
		ignoredReason = "lastFourCandleNotUpTrend"
		return true
	}

	isHaveCountDown := countDownBand(result.Bands[len(result.Bands)-5:]) > 0
	if result.AllTrend.Trend == models.TREND_UP {
		secondDown := result.Bands[len(result.Bands)-2].Candle.Close < result.Bands[len(result.Bands)-2].Candle.Open
		thirdDown := result.Bands[len(result.Bands)-3].Candle.Close < result.Bands[len(result.Bands)-3].Candle.Open
		if CountUpBand(result.Bands[len(result.Bands)-5:]) < 3 || (secondDown && thirdDown) {
			ignoredReason = fmt.Sprintf("count up when trend up: %d", CountUpBand(result.Bands[len(result.Bands)-5:]))
			return true
		}
	} else {
		if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-4:]) < 3 && isHaveCountDown {
			ignoredReason = "count up when trend not up"
			return true
		}
	}

	// if isUpThreeOnMidIntervalChange(result, requestTime) {
	// 	ignoredReason = "isUpThreeOnMidIntervalChange"
	// 	return true
	// }

	// secondLastBand := result.Bands[len(result.Bands)-2]
	// if secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Open > float32(secondLastBand.Upper) {
	// 	ignoredReason = "previous band down from upper"
	// 	return true
	// }

	lastBand := result.Bands[len(result.Bands)-1]
	secondLastBand := result.Bands[len(result.Bands)-2]
	if IsHammer(result.Bands) && result.AllTrend.SecondTrend == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "hammer pattern"
		return true
	}

	if IsDoji(lastBand, true) && result.AllTrend.SecondTrend == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "doji pattern"
		return true
	}

	if ThreeBlackCrowds((result.Bands[len(result.Bands)-4:])) {
		ignoredReason = "three black crowds pattern"
		return true
	}

	buffer := lastBand.Candle.Close * 0.1 / 100
	if lastBand.Candle.Open >= float32(lastBand.Upper)-buffer || (secondLastBand.Candle.Open >= float32(secondLastBand.Upper) && secondLastBand.Candle.Close >= float32(secondLastBand.Upper)) {
		ignoredReason = "lasband or previous band open close above upper"
		return true
	}

	if isUpSignificanAndNotUp(result) {
		ignoredReason = "after up significan and trend not up"
		return true
	}

	isUp := result.AllTrend.Trend == models.TREND_UP && !isHaveCountDown
	if result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrendPercent > 70 && !isReversal(result.Bands) && !isUp {
		ignoredReason = "above sma and just minor up"
		return true
	}

	if isLastBandChangeMoreThan5AndHeadMoreThan3(lastBand) {
		ignoredReason = "last band change more than 5 and head more than 3"
		return true
	}

	if secondLastBand.Candle.Open > float32(secondLastBand.Upper) && secondLastBand.Candle.Close < float32(secondLastBand.Upper) {
		if lastBand.Candle.Close < float32(lastBand.Upper) {
			ignoredReason = "below sma and previous band down from upper"
			return true
		}
	}

	if lastBand.Candle.Hight > float32(lastBand.Upper) && !isHasCrossUpper(result.Bands[len(result.Bands)-6:len(result.Bands)-1], false) {
		ignoredReason = "above upper and just one"
		return true
	}

	if countBelowSMA(result.Bands[len(result.Bands)-7:], true) == 7 {
		if !isHasCrossLower(result.Bands[len(result.Bands)-7:], false) {
			ignoredReason = "last seven band below sma and not cross lower"
			return true
		}
	}

	percentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if isHasCrossUpper(result.Bands[:len(result.Bands)/2], true) && result.AllTrend.Trend == models.TREND_UP && result.Position == models.ABOVE_SMA {
		if !isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) && countBelowSMA(result.Bands[len(result.Bands)-5:], true) >= 1 {
			if result.PriceChanges > 3 && percentFromUpper < 1.5 {
				ignoredReason = "up down, above sma and margin < 1.5"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && (result.AllTrend.FirstTrend == models.TREND_SIDEWAY || result.AllTrend.SecondTrend == models.TREND_SIDEWAY) {
		if isHasCrossUpper(result.Bands[len(result.Bands)-15:], true) && isHasCrossSMA(result.Bands[len(result.Bands)-15:], false) && countBelowSMA(result.Bands[len(result.Bands)-15:], true) == 0 && percentFromUpper < 3 {
			ignoredReason = "above sma and sideway"
			return true
		}
	}

	return ignored(result, masterCoin)
}

func IsIgnoredMidInterval(result *models.BandResult, shortInterval *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) && !isLastBandOrPreviousBandCrossSMA(result.Bands) && !reversalFromLower(*shortInterval) {
		ignoredReason = "isInAboveUpperBandAndDownTrend"
		return true
	}

	if lastBandHeadDoubleBody(result) {
		ignoredReason = "lastBandHeadDoubleBody"
		return true
	}

	lastThreeExceptLastBand := result.Bands[len(result.Bands)-4 : len(result.Bands)-1]
	shortUp := CalculateTrendShort(result.Bands[len(result.Bands)-5:]) == models.TREND_UP || (CalculateTrendShort(result.Bands[len(result.Bands)-4:]) == models.TREND_UP && (isHasCrossSMA(lastThreeExceptLastBand, true) || isHasCrossLower(lastThreeExceptLastBand, false)))
	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && !shortUp {
		ignoredReason = "first trend up and second sideway"
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		ignoredReason = "first trend up and second down"
		return true
	}

	if result.Position == models.ABOVE_UPPER && shortInterval.AllTrend.Trend != models.TREND_UP {
		if result.AllTrend.SecondTrend != models.TREND_UP {
			ignoredReason = "position above upper trend not up"
			return true
		}

		if (CountUpBand(result.Bands[len(result.Bands)-3:]) < 2 || !isHasCrossUpper(result.Bands[len(result.Bands)-4:len(result.Bands)-1], false)) && (result.AllTrend.SecondTrend != models.TREND_UP) {
			ignoredReason = "position above upper but previous band not upper or count up bellow 3"
			return true
		}
	}

	lastBand := result.Bands[len(result.Bands)-1]
	if shortInterval.Position == models.ABOVE_UPPER {
		if CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP || shortInterval.AllTrend.Trend != models.TREND_UP {
			ignoredReason = "short interval above upper and mid trend down or trend not up"
			return true
		}

		checkTrend := (result.AllTrend.FirstTrend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN) || (result.AllTrend.FirstTrend == models.TREND_SIDEWAY && result.AllTrend.SecondTrend == models.TREND_SIDEWAY)
		numberBelowSMA := countBelowSMA(result.Bands[len(result.Bands)/2:], true)
		if numberBelowSMA < 2 && checkTrend && !(isHasCrossLower(result.Bands[len(result.Bands)/2:], false) && lastBand.Candle.Close > float32(lastBand.SMA) && result.AllTrend.Trend != models.TREND_DOWN) {
			ignoredReason = "short interval above upper and mid inteval trend not up"
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	if isContaineBearishEngulfing(result) && !isLastBandOrPreviousBandCrossSMA(result.Bands) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if isLastBandOrPreviousBandCrossSMA(result.Bands) && result.AllTrend.SecondTrend == models.TREND_UP {
		if shortInterval.Position == models.ABOVE_SMA {
			shortLastBand := shortInterval.Bands[len(shortInterval.Bands)-1]
			percent := (shortLastBand.Upper - float64(shortLastBand.Candle.Close)) / float64(shortLastBand.Candle.Close) * 100
			if percent < 1.3 {
				ignoredReason = "mid cross band and short candle near upper band"
				return true
			}
		}

		if isLastBandOrPreviousBandCrossSMA(shortInterval.Bands) {
			ignoredReason = "mid cross band and short candle cross band"
			return true
		}
	}

	if isLastBandOrPreviousBandCrossSMA(shortInterval.Bands) && !isReversal(shortInterval.Bands) {
		if shortInterval.AllTrend.Trend == models.TREND_DOWN || result.AllTrend.Trend == models.TREND_DOWN {
			ignoredReason = "short interval cross  sma and mid interval or short interval trend down"
			return true
		}
	}

	secondLastBand := result.Bands[len(result.Bands)-2]
	if (lastBand.Candle.Open >= float32(lastBand.Upper) && lastBand.Candle.Close > float32(lastBand.Upper)) || (secondLastBand.Candle.Close >= float32(secondLastBand.Upper) && secondLastBand.Candle.Open >= float32(secondLastBand.Upper)) {
		ignoredReason = fmt.Sprintf("open close above upper, %.4f, %.4f", lastBand.Candle.Open, lastBand.Upper)
		return true
	}

	if shortInterval.Position == models.ABOVE_SMA && result.AllTrend.Trend != models.TREND_UP && !isReversal(result.Bands) {
		ignoredReason = "above sma and trend not up"
		return true
	}

	if shortInterval.AllTrend.FirstTrend != models.TREND_UP && shortInterval.AllTrend.SecondTrend == models.TREND_UP && (result.AllTrend.Trend != models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP) {
		shortIntervalHalfBands := shortInterval.Bands[:len(shortInterval.Bands)/2]
		if (shortInterval.AllTrend.FirstTrendPercent <= 50 || isHasCrossLower(shortIntervalHalfBands, false)) && shortInterval.AllTrend.SecondTrendPercent <= 50 {
			if shortInterval.Position == models.ABOVE_SMA {
				shortLastBand := shortInterval.Bands[len(shortInterval.Bands)-1]
				shortPercentFromUpper := (float32(shortLastBand.Upper) - shortLastBand.Candle.Close) / shortLastBand.Candle.Close * 100
				if result.Position == models.BELOW_SMA {
					percentFromSMA := (float32(lastBand.SMA) - lastBand.Candle.Close) / lastBand.Candle.Close * 100
					if shortPercentFromUpper < 3 && percentFromSMA < 3 {
						ignoredReason = "up down, and margin form up bellow threshold"
						return true
					}
				}

				if result.Position == models.ABOVE_SMA {
					percentFromUpper := (float32(lastBand.Upper) - lastBand.Candle.Close) / lastBand.Candle.Close * 100
					if shortPercentFromUpper < 3 && percentFromUpper < 3 {
						ignoredReason = "up down, and margin form up bellow threshold"
						return true
					}
				}
			}
		}
	}

	highest := getHigestHightPrice(result.Bands)
	lowest := getLowestLowPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent <= 2 {
		ignoredReason = "hight and low bellow 3"
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if isHasCrossUpper(result.Bands[:len(result.Bands)/2], false) && result.Position == models.BELOW_SMA {
			ignoredReason = "down-down from upper and position bellow sma"
			return true
		}
	}

	if result.AllTrend.Trend == models.TREND_UP && result.AllTrend.SecondTrendPercent < 5 {
		longestIndex := getLongestCandleIndex(result.Bands[len(result.Bands)/2:])
		secondWaveTrendDetail := CalculateTrendsDetail(result.Bands[longestIndex+len(result.Bands)/2:])
		if secondWaveTrendDetail.FirstTrend != models.TREND_UP || secondWaveTrendDetail.SecondTrend != models.TREND_UP {
			ignoredReason = "after significan up and not up up"
			return true
		}
	}

	if result.AllTrend.FirstTrend != models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent > 20 {
		if lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA) {
			ignoredReason = "on down, not significan and cross sma"
			return true
		}
	}

	if isDoubleUp(result.Bands) {
		ignoredReason = "has double up"
		return true
	}

	if result.AllTrend.SecondTrendPercent < 7 && lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		if CalculateTrendShort(result.Bands[len(result.Bands)-5:]) != models.TREND_UP && !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) {
			ignoredReason = "significan down and below sma, last 5 not up trend, not cross lower"
			return true
		}
	}

	if previousBandDownAndLastBandHammerOrDoji(result.Bands) {
		ignoredReason = "previous band down and last band hammer or doji"
		return true
	}

	shortLastBand := shortInterval.Bands[len(shortInterval.Bands)-1]
	shortPercentFromUpper := (shortLastBand.Upper - float64(shortLastBand.Candle.Close)) / float64(shortLastBand.Candle.Close) * 100
	if result.Position == models.ABOVE_UPPER && !isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-5:], true) {
		if shortPercentFromUpper < 3 {
			ignoredReason = "above upper and short interval 5 band not cross upper and margin below 3"
			return true
		}
	}

	if result.AllTrend.SecondTrendPercent < 15 && (secondLastBand.Candle.Open < float32(secondLastBand.Lower) || secondLastBand.Candle.Close < float32(secondLastBand.Lower)) {
		if CalculateTrendShortAvg(result.Bands[len(result.Bands)-3:]) != models.TREND_UP && lastBand.Candle.Low < float32(lastBand.Lower) {
			ignoredReason = "significan down, last 3 avg not up trend"
			return true
		}
	}

	significanDown := result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 10 && result.AllTrend.SecondTrend == models.TREND_DOWN
	secondWaweBelowSMA := result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 15 && countBelowSMA(result.Bands[len(result.Bands)/2:], false) >= len(result.Bands)/2
	if significanDown || secondWaweBelowSMA {
		if !isHasCrossSMA(result.Bands[len(result.Bands)-5:], true) && !isHasCrossLower(result.Bands[len(result.Bands)-5:], false) {
			ignoredReason = "first wave significan down, second still donw, below sma, not cross sma or lower"
			return true
		}
	}

	percentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if shortLastBand.Candle.Hight > float32(shortLastBand.Upper) && result.Position == models.ABOVE_SMA {
		if !isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) && isHasCrossLower(result.Bands[len(result.Bands)-5:], false) {
			if percentFromUpper < 3 {
				ignoredReason = "short interval cross upper, mid not cross upper but lower and margin below 3"
				return true
			}
		}
	}

	if crossSMAAndPreviousBandNotHaveAboveSMA(result.Bands) {
		ignoredReason = "cross sma and all previous band do not above sma"
		return true
	}

	if isHasCrossUpper(result.Bands[:len(result.Bands)/2], true) && !isHasCrossUpper(result.Bands[len(result.Bands)/2:], true) {
		if result.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_UPPER && percentFromUpper < 3 {
			ignoredReason = "short interval above upper, mid interval above sma and margin below 3"
			return true
		}
	}

	if getHighestIndex(result.Bands) == len(result.Bands)-1 && percentFromUpper < 3 && result.Position == models.ABOVE_SMA {
		if isHasCrossLower(shortInterval.Bands[:len(shortInterval.Bands)/2], false) && isHasCrossUpper(shortInterval.Bands[:len(shortInterval.Bands)/2], true) {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)/2:], true) {
				ignoredReason = "short interval above upper, mid interval above sma and margin below 3 - second"
				return true
			}
		}
	}

	percentFromSMA := (lastBand.SMA - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if result.AllTrend.Trend == models.TREND_DOWN && countBelowSMA(result.Bands[len(result.Bands)/2:], false) >= len(result.Bands)/2 {
		if lastBand.Candle.Open < float32(lastBand.Lower) && lastBand.Candle.Close < float32(lastBand.Lower) && percentFromSMA < 3.5 {
			ignoredReason = "open close below lower and margin below 3.5"
			return true
		}
	}

	if result.Position == models.BELOW_SMA && result.AllTrend.Trend == models.TREND_DOWN && countBelowSMA(result.Bands[len(result.Bands)-7:], false) == 7 {
		if shortLastBand.Candle.Close > float32(shortLastBand.SMA) && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-7:], true) {
			if !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) && percentFromSMA < 3 {
				ignoredReason = "down trend last 7 band below SMA and percent from sma < 3"
				return true
			}
		}
	}

	return false
}

func IsIgnoredLongInterval(result *models.BandResult, shortInterval *models.BandResult, midInterval *models.BandResult, masterTrend, masterMidTrend int8) bool {
	if isInAboveUpperBandAndDownTrend(result) && (CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 || isHasDoji(result.Bands[len(result.Bands)-5:])) {
		if !isLastBandOrPreviousBandCrossSMA(result.Bands) {
			ignoredReason = "isInAboveUpperBandAndDownTrend"
			return true
		}
	}

	lastBand := result.Bands[len(result.Bands)-1]
	if isInAboveUpperBandAndDownTrend(result) && lastBand.Candle.Hight > float32(lastBand.Upper) {
		if lastBand.Candle.Hight-lastBand.Candle.Close > lastBand.Candle.Close-lastBand.Candle.Open {
			ignoredReason = "isInAboveUpperBandAndDownTrend and down from upper"
			return true
		}
	}

	if !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) && !isReversal(shortInterval.Bands) {
		if result.AllTrend.Trend == models.TREND_DOWN && shortInterval.AllTrend.Trend == models.TREND_DOWN && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) == models.TREND_DOWN {
			ignoredReason = "first trend down and seconddown"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend != models.TREND_UP {
		if CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP && getHighestIndex(result.Bands) > 5 {
			ignoredReason = "first trend up and second not up"
			return true
		}
	}

	if isContaineBearishEngulfing(result) && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if result.Position == models.ABOVE_UPPER && result.AllTrend.ShortTrend != models.TREND_UP {
		ignoredReason = "position above upper and not trend up"
		return true
	}

	isMidIntervalTrendNotUp := CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-4:]) == models.TREND_DOWN
	isLongIntervalTrendNotUp := CalculateTrendShort(result.Bands[len(result.Bands)-3:]) == models.TREND_DOWN
	if shortInterval.Position == models.ABOVE_UPPER || midInterval.Position == models.ABOVE_UPPER || result.Position == models.ABOVE_UPPER {
		if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN || isMidIntervalTrendNotUp || isLongIntervalTrendNotUp {
			ignoredReason = "when above upper or master trend down and trend not up"
			return true
		}

		lenData := len(result.Bands)
		hight := getHigestHightPrice(result.Bands[lenData-lenData/4:])
		low := getLowestPrice(result.Bands[lenData-lenData/4:])
		difference := hight - low
		percent := float32(difference) / float32(low) * 100
		if percent > 45 {
			ignoredReason = "up more than 45"
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	shortLastBand := shortInterval.Bands[len(result.Bands)-1]
	percentshort := (float32(shortLastBand.Upper) - shortLastBand.Candle.Close) / shortLastBand.Candle.Close * float32(100)
	if midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
		if result.Position == models.ABOVE_SMA {
			longLastBand := result.Bands[len(result.Bands)-1]
			midLastBand := midInterval.Bands[len(result.Bands)-1]

			percentLong := (float32(longLastBand.Upper) - longLastBand.Candle.Close) / longLastBand.Candle.Close * float32(100)
			percentMid := (float32(midLastBand.Upper) - midLastBand.Candle.Close) / midLastBand.Candle.Close * float32(100)

			if (percentLong < 3.3 && percentMid < 3.3 && percentshort < 3.3) && (shortInterval.AllTrend.SecondTrend != models.TREND_UP || midInterval.AllTrend.SecondTrend != models.TREND_UP || isLongIntervalTrendNotUp) {
				ignoredReason = "all band bellow 3.1 from upper or not up trend"
				return true
			}

			if isHasCrossSMA(result.Bands[len(result.Bands)-1:len(result.Bands)], true) {
				ignoredReason = "all interval above upper and long interval cross sma"
				return true
			}
		}

		if isLastBandOrPreviousBandCrossSMA(result.Bands) && midInterval.AllTrend.Trend != models.TREND_UP {
			ignoredReason = "short and mid above sma but long interval cross sma"
			return true
		}

		if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "short and mid above sma but long interval second wave down trend"
			return true
		}
	}

	if isLastBandCrossUpperAndPreviousBandNot(shortInterval.Bands) {
		if isLastBandCrossUpperAndPreviousBandNot(midInterval.Bands) {
			if isLastBandCrossUpperAndPreviousBandNot(result.Bands) {
				ignoredReason = "band above upper and just one in all interval, skip"
				return true
			}
		}
	}

	if lastBand.Candle.Open > float32(lastBand.Upper) {
		ignoredReason = "open close above upper"
		return true
	}

	allTrendUp := midInterval.AllTrend.Trend == models.TREND_UP && result.AllTrend.Trend == models.TREND_UP
	if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-4:], false) || allTrendUp {
		if result.Position == models.ABOVE_UPPER {
			highestHightIndex := getHighestHightIndex(result.Bands)
			highestIndex := getHighestIndex(result.Bands)
			higestHightBand := result.Bands[highestHightIndex]
			percent := (higestHightBand.Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
			if highestHightIndex == len(result.Bands)-1 || (percent <= 3 && highestIndex == len(result.Bands)-1) {
				ignoredReason = "all interval above upper or all trend up and new hight created"
				return true
			}
		}

		if midInterval.Position == models.ABOVE_UPPER && result.AllTrend.Trend == models.TREND_DOWN && lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA) {
			ignoredReason = "short and mid above upper and long down cross sma"
			return true
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && !isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-5:len(shortInterval.Bands)-1], false) && (midInterval.AllTrend.Trend != models.TREND_UP || result.AllTrend.Trend != models.TREND_UP) {
		ignoredReason = "short interval position above uppper but previous band bellow upper and mid/long interval not up trend"
		return true
	}

	percentFromHeight := (lastBand.Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
	if result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && percentFromHeight < 3 && !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
		ignoredReason = "sideway, above sma and percent from upper bellow 3"
		return true
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent < 10 {
		if result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
			ignoredReason = "significan down, below sma, not cross lower yet"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 20 && result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.SecondTrendPercent < 20 {
		index := getHighestIndex(result.Bands[len(result.Bands)/2:])
		if index < len(result.Bands)/2 && CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && !isLastBandOrPreviousBandCrossSMA(result.Bands) {
			if index < (len(result.Bands)/2)-4 {
				index = (len(result.Bands) / 2) - 4
			}
			lastDataFromHight := result.Bands[(len(result.Bands)/2)+index:]
			if CalculateTrendShort(lastDataFromHight) != models.TREND_UP {
				ignoredReason = "down and up significan then down trend"
				return true
			}
		}
	}

	percentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	secondLastBand := result.Bands[len(result.Bands)-2]
	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.SecondTrend < 10 {
		if lastBand.Candle.Hight < float32(lastBand.Upper) && secondLastBand.Candle.Hight < float32(secondLastBand.Upper) && percentFromUpper < 3 {
			ignoredReason = "up significan but last two band not cross upper"
			return true
		}
	}

	if midInterval.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.Trend == models.TREND_DOWN {
		if midInterval.Position == models.BELOW_SMA && result.Position == models.BELOW_SMA && percentFromUpper <= 3 {
			ignoredReason = "mid and long interval down and below sma and mergin form short below 3"
			return true
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && countBelowSMA(result.Bands[len(result.Bands)-6:len(result.Bands)-1], false) == 5 && result.AllTrend.ShortTrend != models.TREND_UP {
		if lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA) && percentshort <= 3 {
			ignoredReason = "long interval down and below sma (5) and cross sma && mergin form short below 3"
			return true
		}
	}

	if result.AllTrend.SecondTrendPercent < 20 {
		secondWaveBands := result.Bands[len(result.Bands)/2:]
		if countCrossUpper(secondWaveBands) > 2 {
			secondLastBand := result.Bands[len(result.Bands)-2]
			if secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Open > float32(secondLastBand.Upper) {
				ignoredReason = "previous band down from upper"
				return true
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_UP && secondLastBand.Candle.Hight < float32(secondLastBand.Upper) {
		if result.AllTrend.FirstTrendPercent < 20 && result.AllTrend.FirstTrendPercent < result.AllTrend.SecondTrendPercent {
			if isHasCrossUpper(result.Bands[len(result.Bands)-1:], true) && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-1:], true) && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-1:], true) {
				ignoredReason = "all interval cross upper"
				return true
			}
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && midInterval.Position == models.ABOVE_UPPER && result.Position == models.ABOVE_SMA {
		secondWaveBands := result.Bands[len(result.Bands)/2:]
		highestIndex := getHighestIndex(secondWaveBands)
		if highestIndex == len(secondWaveBands)-1 {
			secondHightIndex := getHigestIndexSecond(secondWaveBands)
			if secondHightIndex < len(secondWaveBands)-5 {
				resultTrend := CalculateTrendsDetail(secondWaveBands[secondHightIndex:])
				if resultTrend.FirstTrend == models.TREND_DOWN && resultTrend.SecondTrend == models.TREND_UP {
					higestHightIndex := getHighestHightIndex(secondWaveBands)
					percentFromHight := (secondWaveBands[higestHightIndex].Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
					if percentFromHight < 3 && percentFromUpper < 3 {
						ignoredReason = "on up trend up, new hight below sma and below threshold"
						return true
					}
				}
			}
		}
	}

	if lastBand.Candle.Close > float32(lastBand.Upper) && !isHasCrossUpper(result.Bands[:len(result.Bands)-1], false) {
		midSecondLastBand := midInterval.Bands[len(result.Bands)-2]
		midLastBand := midInterval.Bands[len(result.Bands)-1]
		if midSecondLastBand.Candle.Open > float32(midSecondLastBand.Upper) && midSecondLastBand.Candle.Close < float32(midSecondLastBand.Upper) {
			shortHigestIndex := getHighestIndex(shortInterval.Bands)
			if midLastBand.Candle.Close < float32(midLastBand.Upper) && shortHigestIndex == len(shortInterval.Bands)-1 {
				ignoredReason = "long interval above upper just one and mid interval below sma and previous band down from upper"
				return true
			}
		}
	}

	if result.AllTrend.FirstTrendPercent > 20 && result.AllTrend.SecondTrendPercent > 20 && !isHasCrossUpper(result.Bands, false) {
		if result.AllTrend.FirstTrend != result.AllTrend.SecondTrend || (result.AllTrend.FirstTrend == models.TREND_SIDEWAY && result.AllTrend.SecondTrend == models.TREND_SIDEWAY) {
			if result.Position == models.ABOVE_SMA && percentFromUpper < 3 {
				ignoredReason = "long interval above sma sideway and percent from upper below threshold"
				return true
			}

			percentFromSMA := (float32(lastBand.SMA) - lastBand.Candle.Close) / lastBand.Candle.Close * 100
			if result.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.Position == models.ABOVE_UPPER && result.Position == models.BELOW_SMA && percentFromSMA < 3 {
				ignoredReason = "mid interval above upper and long interval cross sma sideway"
				return true
			}
		}
	}

	if isHasCrossLower(result.Bands[len(result.Bands)-4:], false) && percentFromUpper < 3 {
		ignoredReason = "last 4 band cross lower && margin from upper below threshold"
		return true
	}

	if countCrossSMA(result.Bands[len(result.Bands)-3:len(result.Bands)-1]) == 2 && countCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:len(midInterval.Bands)-1]) == 2 {
		ignoredReason = "two band mid interval cross upper and two band long interval cross sma"
		return true
	}

	if result.AllTrend.SecondTrendPercent < 12 && countBelowSMA(result.Bands[len(result.Bands)-6:], true) == 6 {
		if !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) && CalculateTrendShort(result.Bands[len(result.Bands)-6:]) != models.TREND_UP {
			ignoredReason = "significan down, last 6 bellow sma but not cross lower"
			return true
		}
	}

	if midInterval.Position == models.BELOW_SMA && result.Position == models.BELOW_SMA && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if countBelowSMA(midInterval.Bands, false) == len(midInterval.Bands) && countBelowSMA(result.Bands[len(result.Bands)-5:], false) > 2 {
			if !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) {
				ignoredReason = "mid interval all band below sma, long interval below sma but not cross lower"
				return true
			}
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && shortInterval.AllTrend.Trend == models.TREND_UP {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && !isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], false) {
			midLastBand := midInterval.Bands[len(midInterval.Bands)-1]
			midPercentUpper := (float32(midLastBand.Upper) - midLastBand.Candle.Close) / midLastBand.Candle.Close * 100
			if midPercentUpper < 3 && result.AllTrend.SecondTrend == models.TREND_DOWN {
				ignoredReason = "short interval trend up above upper, mid above sma not cross lower and margin < 3"
				return true
			}
		}
	}

	if countDownBand(result.Bands[len(result.Bands)-6:len(result.Bands)-1]) == 5 && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) == models.TREND_DOWN {
		ignoredReason = "squential down 5 band, and short trend not up"
		return true
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.Position == models.BELOW_SMA {
		if isHasDoji(result.Bands[len(result.Bands)-3:]) || isHasHammer(result.Bands[len(result.Bands)-3:]) {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) && isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-3:], true) {
				ignoredReason = "trend down, below sma, doji or hammer pattern"
				return true
			}
		}
	}

	midLastBand := midInterval.Bands[len(midInterval.Bands)-1]
	if result.Position == models.ABOVE_SMA {
		if isHasCrossUpper(result.Bands[:len(result.Bands)/2], true) {
			if countCrossUpper(shortInterval.Bands) >= 9 && midInterval.Position == models.ABOVE_UPPER && midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.PriceChanges > 5 {
				ignoredReason = "up down, mid and short above upper"
				return true
			}
		}

		if isHasCrossLower(result.Bands, false) && !isHasCrossUpper(result.Bands, true) && percentFromUpper < 3 {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:], true) {
				ignoredReason = "below sma, margin from upper below 3"
				return true
			}
		}

		if result.AllTrend.ShortTrend == models.TREND_DOWN && isHasOpenCloseAboveUpper(result.Bands[len(result.Bands)-7:]) {
			midPercentFromSMA := (midLastBand.SMA - float64(midLastBand.Candle.Close)) / float64(midLastBand.Candle.Close) * 100
			if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				if midInterval.AllTrend.FirstTrendPercent > 20 && midInterval.AllTrend.SecondTrendPercent > 20 && midInterval.Position == models.BELOW_SMA {
					if isHasCrossLower(midInterval.Bands[len(result.Bands)/2:], false) && midPercentFromSMA < 3 {
						ignoredReason = "mid interval from upper down to lower but not significan"
						return true
					}
				}
			}
		}

		if lastBand.Candle.Open < float32(lastBand.SMA) && midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Close > float32(midLastBand.SMA) {
			if shortLastBand.Candle.Open < float32(shortLastBand.Upper) && shortLastBand.Candle.Hight > float32(shortLastBand.Upper) {
				ignoredReason = "long interval cross sma, mid interval cross sma and short interval cross upper"
				return true
			}
		}

		if lastBand.Candle.Low < float32(lastBand.SMA) && result.AllTrend.SecondTrend == models.TREND_DOWN && headMoreThan30PrecentToBody(lastBand) {
			if midLastBand.Candle.Hight > float32(midLastBand.Upper) && headMoreThan30PrecentToBody(midLastBand) {
				if shortInterval.Position == models.ABOVE_UPPER && getHighestIndex(shortInterval.Bands) == len(shortInterval.Bands)-1 && headMoreThanBody(shortLastBand) {
					ignoredReason = "head more than 30 on long and mid, short interval head more than body"
					return true
				}
			}
		}
	}

	midPercentFromUpper := (midLastBand.Upper - float64(midLastBand.Candle.Close)) / float64(midLastBand.Candle.Close) * 100
	if lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA) {
		if shortInterval.Position == models.ABOVE_UPPER && countCrossUpper(shortInterval.Bands[len(shortInterval.Bands)/2:]) > 5 {
			if (midInterval.Position == models.ABOVE_SMA && midPercentFromUpper < 3) || (midInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands[len(midInterval.Bands)-5:]) == 1) {
				ignoredReason = "cross sma and short cross upper and mid margin below 3"
				return true
			}
		}
	}

	percentFromSMA := (lastBand.SMA - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if result.Position == models.BELOW_SMA && result.AllTrend.Trend == models.TREND_DOWN && countBelowSMA(result.Bands[len(result.Bands)-7:], false) == 7 {
		if shortLastBand.Candle.Close > float32(shortLastBand.SMA) && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-7:], true) {
			if percentFromSMA < 3 {
				ignoredReason = "down trend last 7 band below SMA and percent from sma < 3"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && percentFromUpper < 3 {
		if (result.AllTrend.FirstTrend != models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP) || (getHighestIndex(result.Bands) == len(result.Bands)-1) {
			if shortLastBand.Candle.Open < float32(shortLastBand.Upper) && shortLastBand.Candle.Close > float32(shortLastBand.Upper) {
				if midLastBand.Candle.Open < float32(midLastBand.Upper) && midLastBand.Candle.Close > float32(midLastBand.Upper) {
					ignoredReason = "above sma and percent from upper below 3"
					return true
				}
			}
		}
	}

	return false
}

func IsIgnoredMasterDown(result, midInterval, longInterval, masterCoin *models.BandResult, checkingTime time.Time) bool {
	minPercentChanges := 2
	midLowestIndex := getLowestIndexSecond(midInterval.Bands)
	if !isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-4:], false) {
		if midLowestIndex < len(midInterval.Bands)-4 {
			ignoredReason = "mid interval is not in lower"
			return true
		}

		if isNotInLower(result.Bands, false) {
			ignoredReason = "is not in lower"
			return true
		}

		if checkingTime.Minute() < 18 {
			ignoredReason = "mid not cross lower and on mid interval time change"
			return true
		}

		if (result.AllTrend.FirstTrend != models.TREND_DOWN || result.AllTrend.FirstTrendPercent > 10) && (result.AllTrend.SecondTrend != models.TREND_DOWN || result.AllTrend.SecondTrendPercent > 10) {
			ignoredReason = "mid not cross lower and short interval not significan down"
			return true
		}

		minPercentChanges = 3
	}

	if midInterval.Direction == BAND_DOWN {
		ignoredReason = "mid interval band down"
		return true
	}

	if CalculateTrendShort(masterCoin.Bands[len(masterCoin.Bands)-4:]) != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		ignoredReason = "short trend not up"
		return true
	}

	lastBand := result.Bands[len(result.Bands)-1]
	marginFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if marginFromUpper < float64(minPercentChanges) && !isReversal(midInterval.Bands) && midLowestIndex < len(midInterval.Bands)-2 {
		ignoredReason = fmt.Sprintf("margin from upper is bellow %d", minPercentChanges)
		return true
	}

	midLastBand := midInterval.Bands[len(midInterval.Bands)-1]
	if result.Position == models.ABOVE_SMA && (midInterval.Position == models.ABOVE_SMA || (midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Hight > float32(midLastBand.SMA))) {
		if marginFromUpper < float64(minPercentChanges) {
			ignoredReason = fmt.Sprintf("2 margin from upper is bellow %d", minPercentChanges)
			return true
		}
	}

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-5:]) < 4 {
		ignoredReason = "count up bellow 2"
		return true
	}

	if CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-3:]) != models.TREND_UP && CountSquentialUpBand(midInterval.Bands[len(midInterval.Bands)-3:]) < 2 && midLowestIndex != len(midInterval.Bands)-1 {
		ignoredReason = "last three band not up"
		return true
	}

	if midInterval.AllTrend.Trend != models.TREND_DOWN || midInterval.AllTrend.FirstTrend == models.TREND_UP {
		ignoredReason = "trend not down or first trend is up"
		return true
	}

	// if isCrossLowerWhenSignificanDown(midInterval) {
	// 	ignoredReason = "mid interval cross lower on significan down"
	// 	return true
	// }

	if isGetBearishEngulfingAfterLowest(result.Bands) {
		ignoredReason = "contain bearish engulfing"
		return true
	}

	secondLastBand := result.Bands[len(result.Bands)-2]
	thirdLastBand := result.Bands[len(result.Bands)-3]
	if result.Position == models.BELOW_SMA && secondLastBand.Candle.Close < secondLastBand.Candle.Open && thirdLastBand.Candle.Hight > float32(thirdLastBand.SMA) {
		ignoredReason = "sideway after hit sma"
		return true
	}

	if result.Position == models.BELOW_SMA && secondLastBand.Candle.Hight > float32(secondLastBand.SMA) {
		ignoredReason = "short interval previous band hit SMA and current band bellow SMA"
		return true
	}

	if lastBand.Candle.Close < secondLastBand.Candle.Close {
		ignoredReason = "close below previous band"
		return true
	}

	oneDownTrend := midInterval.AllTrend.FirstTrendPercent > 10 && midInterval.AllTrend.SecondTrendPercent > 10
	bothDownTrend := midInterval.AllTrend.FirstTrend != models.TREND_DOWN || midInterval.AllTrend.SecondTrend != models.TREND_DOWN || (midInterval.AllTrend.FirstTrendPercent > 20 && midInterval.AllTrend.SecondTrendPercent > 20)
	if oneDownTrend && bothDownTrend {
		ignoredReason = "mid interval not significan down"
		return true
	}

	if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && isHasCrossUpper(midInterval.Bands[:len(midInterval.Bands)/2], true) {
		ignoredReason = "mid interval down but cross upper"
		return true
	}

	if isInAboveUpperBandAndDownTrend(longInterval) && longInterval.AllTrend.ShortTrend != models.TREND_UP {
		ignoredReason = "long interval down from hight, not reversal"
		return true
	}

	if midLastBand.Candle.Close < float32(midLastBand.Lower) {
		ignoredReason = "mid interval open close bellow lower"
		return true
	}

	longLastBand := longInterval.Bands[len(longInterval.Bands)-1]
	if longInterval.AllTrend.ShortTrend == models.TREND_DOWN && longLastBand.Candle.Low < float32(longLastBand.SMA) && longLastBand.Candle.Hight > float32(longLastBand.SMA) {
		ignoredReason = "long interval down and cross sma"
		return true
	}

	if previousBandDownAndLastBandHammerOrDoji(midInterval.Bands) {
		ignoredReason = "mid interval previous band down and last band hammer or doji"
		return true
	}

	if longInterval.AllTrend.FirstTrend == models.TREND_UP && longInterval.AllTrend.SecondTrend == models.TREND_DOWN && (longInterval.AllTrend.FirstTrendPercent < 10 || longInterval.AllTrend.SecondTrendPercent < 10) {
		if longInterval.Position == models.BELOW_SMA && !isHasCrossLower(longInterval.Bands[len(longInterval.Bands)/2:], false) && longInterval.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "long interval up down, below sma, not cross lower yet"
			return true
		}
	}

	if longInterval.AllTrend.Trend == models.TREND_DOWN && longLastBand.Candle.Close < float32(longLastBand.SMA) {
		if isHasCrossLower(longInterval.Bands[len(longInterval.Bands)/2:], false) && longInterval.Direction == BAND_DOWN {
			ignoredReason = "long interval down cross upper, but band not up"
			return true
		}
	}

	longSecondBand := longInterval.Bands[len(longInterval.Bands)-2]
	if longInterval.Position == models.BELOW_SMA && longSecondBand.Candle.Hight > float32(longSecondBand.SMA) && longInterval.AllTrend.SecondTrend != models.TREND_UP {
		ignoredReason = "long interval previous band hit SMA and current band bellow SMA"
		return true
	}

	if countCrossLower(longInterval.Bands[len(longInterval.Bands)-4:len(longInterval.Bands)-1]) == 3 {
		ignoredReason = "long interval 3 band cross lower"
		return true
	}

	trendDetail := CalculateTrendsDetail(longInterval.Bands[len(longInterval.Bands)/2:])
	if trendDetail.SecondTrendPercent < 5 && longInterval.Position == models.ABOVE_SMA {
		if countDownBand(longInterval.Bands[len(longInterval.Bands)-5:]) > 2 && !isHasCrossSMA(longInterval.Bands[len(longInterval.Bands)-5:], true) {
			ignoredReason = "above sma and significan down"
			return true
		}
	}

	if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
		longFirstBand := longInterval.Bands[0]
		if countBelowSMA(longInterval.Bands[len(longInterval.Bands)-2:], true) == 2 && marginFromUpper < 3 {
			if countCrossLower(longInterval.Bands[len(longInterval.Bands)-4:]) < 2 || longFirstBand.Candle.Close < lastBand.Candle.Close {
				ignoredReason = "mid down-down, long interval below sma not cross lower"
				return true
			}
		}
	}

	if countBelowSMA(longInterval.Bands[len(longInterval.Bands)-5:], true) >= 4 && CalculateTrendShort(longInterval.Bands[len(longInterval.Bands)-3:]) != models.TREND_UP {
		if result.AllTrend.Trend == models.TREND_DOWN && !isHasCrossLower(longInterval.Bands[len(longInterval.Bands)-4:], true) {
			ignoredReason = "long interval last five band bellow sma and not cross lower"
			return true
		}
	}

	if countDownBand(longInterval.Bands[len(longInterval.Bands)-5:]) == 4 && countCrossLower(longInterval.Bands[len(longInterval.Bands)-5:]) > 1 {
		if !isHasHammer(longInterval.Bands[len(longInterval.Bands)-3:len(longInterval.Bands)-1]) && !isHasDoji(longInterval.Bands[len(longInterval.Bands)-3:len(longInterval.Bands)-1]) {
			ignoredReason = "long intervallast 4 band down and cross lowe, no reversal pattern"
			return true
		}
	}

	if isContaineBearishEngulfing(longInterval) {
		ignoredReason = "long interval isContaineBearishEngulfing"
		return true
	}

	if lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA) {
		if secondLastBand.Candle.Open > secondLastBand.Candle.Close {
			ignoredReason = "short interval cross sma and previous band down"
			return true
		}
	}

	if longLastBand.Candle.Open < float32(longLastBand.Lower) && longLastBand.Candle.Close < float32(longLastBand.Lower) {
		ignoredReason = "long interval open close below lower"
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		ignoredReason = "short interval above upper"
		return true
	}

	return false
}

func headMoreThan30PrecentToBody(band models.Band) bool {
	hightFromOpen := band.Candle.Hight - band.Candle.Open
	hightHead := band.Candle.Hight - band.Candle.Close

	return (hightHead/hightFromOpen)*100 > 30
}

func headMoreThanBody(band models.Band) bool {
	hightBody := band.Candle.Close - band.Candle.Open
	hightHead := band.Candle.Hight - band.Candle.Close

	return hightHead > hightBody
}

func crossSMAAndPreviousBandNotHaveAboveSMA(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	if lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA) {
		for _, band := range bands {
			if band.Candle.Open > float32(band.SMA) && band.Candle.Close > float32(band.SMA) {
				return false
			}
		}
		return true
	}

	return false
}

func previousBandDownAndLastBandHammerOrDoji(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	if secondLastBand.Candle.Close < secondLastBand.Candle.Open && isHasCrossLower(bands[len(bands)-2:], false) {
		if IsHammer([]models.Band{lastBand}) || IsDoji(lastBand, true) {
			return true
		}
	}

	return false
}

func isGetBearishEngulfingAfterLowest(bands []models.Band) bool {
	lowestIndex := getLowestLowIndex(bands)
	if lowestIndex < len(bands)-7 {
		lowestIndex = len(bands) - 7
	}
	return BearishEngulfing(bands[lowestIndex:])
}

func isReversal(bands []models.Band) bool {
	lowestIndex := getLowestIndex(bands[len(bands)-4:])
	trend := CalculateTrendsDetail(bands[:len(bands)-4+lowestIndex])
	shortTrend := CalculateTrendShort(bands[len(bands)-4:])
	return trend.Trend == models.TREND_DOWN && shortTrend == models.TREND_UP
}

func reversalFromLower(result models.BandResult) bool {
	trend := CalculateTrendsDetail(result.Bands[:len(result.Bands)-1])
	if isHasCrossLower(result.Bands[len(result.Bands)-4:], false) && isLowerDifferentValid(result.Bands) {
		return trend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_UP
	}
	return false
}

func isLowerDifferentValid(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	midBand := bands[len(bands)/2]
	var percent float32 = 0
	if lastBand.Lower > midBand.Lower {
		percent = float32(midBand.Lower) / float32(lastBand.Lower) * 100
	} else {
		percent = float32(lastBand.Lower) / float32(midBand.Lower) * 100
	}

	return percent > 98.898
}

// func isCrossLowerWhenSignificanDown(result *models.BandResult) bool {
// 	lastBand := result.Bands[len(result.Bands)-1]
// 	if lastBand.Candle.Open < float32(lastBand.Lower) && lastBand.Candle.Close > float32(lastBand.Lower) {
// 		return result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent < 50
// 	}

// 	return false
// }

func isLastBandCrossUpperAndPreviousBandNot(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	if lastBand.Candle.Open < lastBand.Candle.Close {
		secondLastBand := bands[len(bands)-2]
		return !(secondLastBand.Candle.Open < secondLastBand.Candle.Close)
	}
	return false
}

func isHasCrossUpper(bands []models.Band, withHead bool) bool {
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			if withHead {
				if band.Candle.Open < float32(band.Upper) && band.Candle.Hight > float32(band.Upper) {
					return true
				}
			} else {
				if band.Candle.Open < float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
					return true
				}
			}
		}
	}
	return false
}

func isHasCrossSMA(bands []models.Band, bodyOnly bool) bool {
	for _, band := range bands {
		if bodyOnly {
			if band.Candle.Open < float32(band.SMA) && band.Candle.Close > float32(band.SMA) {
				return true
			}
		} else {
			if band.Candle.Low < float32(band.SMA) && band.Candle.Hight > float32(band.SMA) {
				return true
			}
		}
	}
	return false
}

func isHasOpenCloseAboveUpper(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open > float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
			return true
		}
	}
	return false
}

func isNotInLower(bands []models.Band, skipped bool) bool {
	lowestIndex := getLowestIndex(bands)
	if !isHasCrossLower(bands[len(bands)-10:], false) {
		if isHasCrossLower(bands[len(bands)-20:], false) || skipped {
			return lowestIndex < len(bands)-10
		}
		return true
	}

	return false
}

func isHasCrossLower(bands []models.Band, bodyOnly bool) bool {
	crossLowerBand := false
	for _, data := range bands {
		if bodyOnly {
			if (data.Candle.Close < float32(data.Lower) && data.Candle.Open > float32(data.Lower)) || (data.Candle.Open < float32(data.Lower) && data.Candle.Close > float32(data.Lower)) {
				crossLowerBand = true
				break
			}
		} else {
			if data.Candle.Low < float32(data.Lower) && data.Candle.Close > float32(data.Lower) {
				crossLowerBand = true
				break
			}
		}
	}

	return crossLowerBand
}

func countCrossLower(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Low < float32(data.Lower) && data.Candle.Hight > float32(data.Lower) {
			count++
		}
	}

	return count
}

func countCrossUpper(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Open < float32(data.Upper) && data.Candle.Hight > float32(data.Upper) {
			count++
		}
	}

	return count
}

func countCrossSMA(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Low < float32(data.SMA) && data.Candle.Hight > float32(data.SMA) {
			count++
		}
	}

	return count
}

func countBelowSMA(bands []models.Band, strict bool) int {
	count := 0
	for _, data := range bands {
		if strict {
			if data.Candle.Close < float32(data.SMA) && data.Candle.Open < float32(data.SMA) {
				count++
			}
		} else {
			if data.Candle.Close < float32(data.SMA) {
				count++
			}
		}
	}

	return count
}

func isLastBandOrPreviousBandCrossSMA(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	var isSecondLastBandCrossSMA bool
	if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
		isSecondLastBandCrossSMA = secondLastBand.Candle.Open <= float32(secondLastBand.SMA) && secondLastBand.Candle.Hight >= float32(secondLastBand.SMA)
	} else {
		isSecondLastBandCrossSMA = secondLastBand.Candle.Low <= float32(secondLastBand.SMA) && secondLastBand.Candle.Hight >= float32(secondLastBand.SMA)
	}
	isLastBandCrossSMA := lastBand.Candle.Low <= float32(lastBand.SMA) && lastBand.Candle.Hight >= float32(lastBand.SMA)

	return isLastBandCrossSMA || isSecondLastBandCrossSMA
}

// func isUpThreeOnMidIntervalChange(result *models.BandResult, requestTime time.Time) bool {
// 	lastBand := result.Bands[len(result.Bands)-1]
// 	isCrossSMA := lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA)
// 	isCrossUpper := lastBand.Candle.Low < float32(lastBand.Upper) && lastBand.Candle.Hight > float32(lastBand.Upper)
// 	if isCrossSMA || isCrossUpper {
// 		if CountSquentialUpBand(result.Bands[len(result.Bands)-4:]) >= 3 && requestTime.Minute() < 17 {
// 			return true
// 		}
// 	}

// 	return false
// }

func isPosititionBellowUpperMarginBellowThreshold(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Close < float32(lastBand.SMA) && result.AllTrend.SecondTrend != models.TREND_UP {
		margin := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100

		return margin < 2.1
	}

	return false
}

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	index := getHighestIndex(result.Bands)
	lastBand := result.Bands[len(result.Bands)-1]
	if index == len(result.Bands)-1 || lastBand.Candle.Close < float32(lastBand.SMA) {
		return false
	}

	if index > len(result.Bands)-5 {
		index = len(result.Bands) - 5
	}
	lastDataFromHight := result.Bands[index:]
	if ((result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_UP) || isHeighestOnHalfEndAndAboveUpper(result)) && CalculateTrendShort(lastDataFromHight) != models.TREND_UP {
		return true
	}

	return false
}

func isHeighestOnHalfEndAndAboveUpper(result *models.BandResult) bool {
	hiIndex := getHighestIndex(result.Bands)
	if hiIndex >= len(result.Bands)/2 {
		return isHasCrossUpper(result.Bands[len(result.Bands)-5:], true)
	}

	return false
}

func isContaineBearishEngulfing(result *models.BandResult) bool {
	hiIndex := len(result.Bands) - 4
	return BearishEngulfing(result.Bands[hiIndex:]) && CalculateTrendShort(result.Bands[hiIndex:]) == models.TREND_DOWN
}

func getHighestIndex(bands []models.Band) int {
	hiIndex := 0
	for i, band := range bands {
		if bands[hiIndex].Candle.Close <= band.Candle.Close {
			hiIndex = i
		}
	}

	return hiIndex
}

func getHighestHightIndex(bands []models.Band) int {
	hiIndex := 0
	for i, band := range bands {
		if bands[hiIndex].Candle.Hight <= band.Candle.Hight {
			hiIndex = i
		}
	}

	return hiIndex
}

func isDoubleUp(bands []models.Band) bool {
	secondWaveBands := bands[len(bands)/2:]
	if countCrossUpper(secondWaveBands) == 2 {
		hiIndex := getHighestHightIndex(secondWaveBands)
		secondHiIndex := -1
		for i, band := range secondWaveBands {
			if i != hiIndex {
				if secondHiIndex < 0 {
					secondHiIndex = i
				} else if secondWaveBands[secondHiIndex].Candle.Hight <= band.Candle.Hight {
					secondHiIndex = i
				}
			}
		}

		if hiIndex == len(secondWaveBands)-1 || secondHiIndex == len(secondWaveBands)-1 {
			different := 0
			if hiIndex < secondHiIndex {
				different = secondHiIndex - hiIndex
			} else {
				different = hiIndex - secondHiIndex
			}
			percent := secondWaveBands[secondHiIndex].Candle.Hight / secondWaveBands[hiIndex].Candle.Hight * 100

			return different >= 5 && percent > 97
		}
	}
	return false
}

func getLowestIndex(bands []models.Band) int {
	lowIndex := 0
	for i, band := range bands {
		if lowestFromBand(bands[lowIndex]) >= lowestFromBand(band) {
			lowIndex = i
		}
	}

	return lowIndex
}

func lowestFromBand(band models.Band) float32 {
	if band.Candle.Open > band.Candle.Close {
		return band.Candle.Close
	}

	return band.Candle.Open
}

func getLowestLowIndex(bands []models.Band) int {
	lowIndex := 0
	for i, band := range bands {
		if bands[lowIndex].Candle.Low > band.Candle.Low {
			lowIndex = i
		}
	}

	return lowIndex
}

func getLowestIndexSecond(bands []models.Band) int {
	firstLow := getLowestIndex(bands)

	secondLow := -1
	for i, band := range bands {
		if i != firstLow && lowestFromBand(bands[firstLow]) != lowestFromBand(band) {
			if secondLow < 0 {
				secondLow = i
			} else if lowestFromBand(bands[secondLow]) >= lowestFromBand(band) {
				secondLow = i
			}
		}
	}

	if firstLow > secondLow {
		return firstLow
	}

	return secondLow
}

func getHigestIndexSecond(bands []models.Band) int {
	firstHight := getHighestIndex(bands)

	secondHight := -1
	for i, band := range bands {
		if i != firstHight {
			if secondHight < 0 {
				secondHight = i
			} else if bands[secondHight].Candle.Close < band.Candle.Close {
				secondHight = i
			}
		}
	}

	return secondHight
}

func whenHeightTripleAverage(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open
	var totalHeight float32 = 0
	for _, band := range result.Bands[len(result.Bands)-6:] {
		if band.Candle.Open > band.Candle.Close {
			totalHeight += band.Candle.Open - band.Candle.Close
		} else {
			totalHeight += band.Candle.Close - band.Candle.Open
		}
	}
	average := totalHeight / float32(6)
	percent := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100

	return lastBandHeight > 3*average && percent > 2.5
}

func lastBandHeadDoubleBody(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Close > lastBand.Candle.Open {
		head := lastBand.Candle.Hight - lastBand.Candle.Close
		body := lastBand.Candle.Close - lastBand.Candle.Open
		return head > body*2.99
	}

	return false
}

func ignored(result, masterCoin *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open <= float32(lastBand.SMA) && lastBand.Candle.Hight >= float32(lastBand.Upper) {
		ignoredReason = "up from bellow sma to upper"
		return true
	}

	return false
}

func lastFourCandleNotUpTrend(bands []models.Band) bool {
	return CalculateTrendShort(bands[len(bands)-4:]) != models.TREND_UP
}

func isTrendUpLastThreeBandHasDoji(result *models.BandResult) bool {
	if result.AllTrend.SecondTrend != models.TREND_UP {
		return false
	}

	lastThreeBand := result.Bands[len(result.Bands)-2:]
	var difference float32 = 0
	var percent float32 = 0
	for _, band := range lastThreeBand {
		if band.Candle.Close > band.Candle.Open {
			difference = band.Candle.Close - band.Candle.Open
			percent = difference / band.Candle.Open * 100
		} else {
			difference = band.Candle.Open - band.Candle.Close
			percent = difference / band.Candle.Close * 100
		}

		if percent < 0.09 {
			return true
		}
	}

	return false
}

func isUpSignificanAndNotUp(result *models.BandResult) bool {
	if result.AllTrend.SecondTrendPercent < 40 && result.AllTrend.SecondTrend == models.TREND_UP {
		mid := len(result.Bands) / 2
		indexDoubleBody := getIndexBandDoubleLong(result.Bands[len(result.Bands)-mid:])
		if indexDoubleBody > -1 {
			realIndex := len(result.Bands)%2 + mid + indexDoubleBody
			if len(result.Bands)-realIndex > 4 {
				trend := CalculateTrendsDetail(result.Bands[15:])
				return trend.Trend != models.TREND_UP
			}
		}
	}

	return false
}

func getIndexBandDoubleLong(bands []models.Band) int {
	longestIndex := -1
	var total float32 = 0
	var bandLong float32 = 0
	for i, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			bandLong = band.Candle.Close - band.Candle.Open
		} else {
			bandLong = band.Candle.Open - band.Candle.Close
		}

		if band.Candle.Close > band.Candle.Open {
			if longestIndex != -1 {
				if bands[longestIndex].Candle.Close-bands[longestIndex].Candle.Open < bandLong {
					longestIndex = i
				}
			} else {
				longestIndex = i
			}
		}

		total += bandLong
	}

	if longestIndex >= 0 && longestIndex < len(bands)-4 {
		hightBand := bands[longestIndex]
		hight := hightBand.Candle.Close - hightBand.Candle.Open
		total -= hight
		if (total/float32(len(bands)-1))*2 > hight {
			return -1
		}
	}

	return longestIndex
}

func getLongestCandleIndex(bands []models.Band) int {
	longestIndex := 0
	for i, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			if bands[longestIndex].Candle.Close-bands[longestIndex].Candle.Open < band.Candle.Close-band.Candle.Open {
				longestIndex = i
			}
		}
	}

	return longestIndex
}

func getIndexHigestCrossUpper(bands []models.Band) int {
	higestIndex := -1
	lastBand := bands[len(bands)-1]
	for i, band := range bands {
		if band.Candle.Close > lastBand.Candle.Close {
			if higestIndex != -1 {
				if bands[higestIndex].Candle.Close < band.Candle.Close {
					higestIndex = i
				}
			} else {
				higestIndex = i
			}
		}
	}

	return higestIndex
}

func isLastBandChangeMoreThan5AndHeadMoreThan3(lastBand models.Band) bool {
	percentBody := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100
	percentHead := (lastBand.Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
	return percentBody > 5 && percentHead > 3
}

func countDownBand(bands []models.Band) int {
	counter := 0
	for _, band := range bands {
		if band.Candle.Open > band.Candle.Close {
			counter++
		}
	}

	return counter
}

func isHasHammer(bands []models.Band) bool {
	for _, band := range bands {
		if IsHammer([]models.Band{band}) {
			return true
		}
	}

	return false
}

func isHasDoji(bands []models.Band) bool {
	for _, band := range bands {
		isUp := band.Candle.Open < band.Candle.Close
		if IsDoji(band, isUp) {
			return true
		}
	}

	return false
}

func GetIgnoredReason() string {
	return ignoredReason
}

// tambah kondisi untuk bearish engulfing onsell, ketika kurang dari 3 check mid interval bearish engulfing jg? candle complete? atau udah turun 0.5 %
// adjust sell log bos
