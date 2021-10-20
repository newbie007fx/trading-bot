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

	if lastFourCandleNotUpTrend(result.Bands) {
		ignoredReason = "lastFourCandleNotUpTrend"
		return true
	}

	if result.Trend == models.TREND_UP {
		secondDown := result.Bands[len(result.Bands)-2].Candle.Close < result.Bands[len(result.Bands)-2].Candle.Open
		thirdDown := result.Bands[len(result.Bands)-3].Candle.Close < result.Bands[len(result.Bands)-3].Candle.Open
		if CountUpBand(result.Bands[len(result.Bands)-5:]) < 3 || (secondDown && thirdDown) {
			ignoredReason = fmt.Sprintf("count up when trend up: %d", CountUpBand(result.Bands[len(result.Bands)-5:]))
			return true
		}
	} else {
		if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-4:]) < 3 {
			ignoredReason = "count up when trend not up"
			return true
		}
	}

	if isUpThreeOnMidIntervalChange(result, requestTime) {
		ignoredReason = "isUpThreeOnMidIntervalChange"
		return true
	}

	secondLastBand := result.Bands[len(result.Bands)-2]
	if result.Position == models.ABOVE_UPPER && !isHasCrossUpper(result.Bands[len(result.Bands)-4:len(result.Bands)-1]) {
		ignoredReason = "position above uppper but previous band bellow upper"
		return true
	}

	if secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Open > float32(secondLastBand.Upper) {
		ignoredReason = "previous band down from upper"
		return true
	}

	if IsHammer(result.Bands) && result.AllTrend.SecondTrend != models.TREND_DOWN {
		ignoredReason = "hammer pattern"
		return true
	}

	if ThreeBlackCrowds((result.Bands[len(result.Bands)-5:])) {
		ignoredReason = "three black crowds pattern"
		return true
	}

	lastBand := result.Bands[len(result.Bands)-1]
	buffer := lastBand.Candle.Close * 0.1 / 100
	if lastBand.Candle.Open >= float32(lastBand.Upper)-buffer {
		ignoredReason = "open close above upper"
		return true
	}

	if isUpSignificanAndNotUp(result) {
		ignoredReason = "after up significan and trend not up"
		return true
	}

	if afterUpThenDown(result) {
		ignoredReason = "after up then down"
		return true
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrendPercent > 70 && !isReversal(result.Bands) {
		ignoredReason = "above sma and just minor up"
		return true
	}

	if isLastBandChangeMoreThan5AndHeadMoreThan3(lastBand) {
		ignoredReason = "last band change more than 5 and head more than 3"
		return true
	}

	return ignored(result, masterCoin)
}

func IsIgnoredMidInterval(result *models.BandResult, shortInterval *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) {
		ignoredReason = "isInAboveUpperBandAndDownTrend"
		return true
	}

	if lastBandHeadDoubleBody(result) {
		ignoredReason = "lastBandHeadDoubleBody"
		return true
	}

	if result.Trend == models.TREND_DOWN && shortInterval.Trend == models.TREND_DOWN && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		ignoredReason = "first trend down and second not up"
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && CalculateTrendShort(result.Bands[len(result.Bands)-5:]) != models.TREND_UP {
		ignoredReason = "first trend up and second sideway"
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		ignoredReason = "first trend up and second down"
		return true
	}

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-4:]) < 3 && !isReversal(result.Bands) {
		ignoredReason = "count up"
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		if result.Trend != models.TREND_UP || result.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "position above upper trend not up"
			return true
		}

		secondLastBand := result.Bands[len(result.Bands)-2]
		if CountUpBand(result.Bands[len(result.Bands)-3:]) < 2 || !(secondLastBand.Candle.Close > secondLastBand.Candle.Open && secondLastBand.Candle.Close > float32(secondLastBand.Upper)) {
			ignoredReason = "position above upper but previous band not upper or count up bellow 3"
			return true
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER {
		if CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP || shortInterval.Trend != models.TREND_UP {
			ignoredReason = "short interval above upper and mid trend down or trend not up"
			return true
		}

		if result.Trend != models.TREND_UP {
			ignoredReason = "short interval above upper and mid inteval trend not up"
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	if isContaineBearishEngulfing(result) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if isLastBandOrPreviousBandCrossSMA(result.Bands) {
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
		if shortInterval.Trend == models.TREND_DOWN || result.Trend == models.TREND_DOWN {
			ignoredReason = "short interval cross  sma and mid interval or short interval trend down"
			return true
		}
	}

	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open > float32(lastBand.Upper) {
		ignoredReason = fmt.Sprintf("open close above upper, %.2f, %.2f", lastBand.Candle.Open, lastBand.Upper)
		return true
	}

	if shortInterval.Position == models.ABOVE_SMA && result.Trend != models.TREND_UP && !isReversal(result.Bands) {
		ignoredReason = "above sma and trend not up"
		return true
	}

	if afterUpThenDown(result) {
		ignoredReason = "after up then down"
		return true
	}

	if shortInterval.AllTrend.FirstTrend != models.TREND_UP && shortInterval.AllTrend.SecondTrend == models.TREND_UP {
		shortIntervalHalfBands := shortInterval.Bands[:len(shortInterval.Bands)/2]
		if (shortInterval.AllTrend.FirstTrendPercent <= 50 || isHasCrossLower(shortIntervalHalfBands)) && shortInterval.AllTrend.SecondTrendPercent <= 50 {
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

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent <= 3 {
		ignoredReason = "hight and low bellow 3"
		return true
	}

	return false
}

func IsIgnoredLongInterval(result *models.BandResult, shortInterval *models.BandResult, midInterval *models.BandResult, masterTrend, masterMidTrend int8) bool {
	if isInAboveUpperBandAndDownTrend(result) && CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 {
		ignoredReason = "isInAboveUpperBandAndDownTrend"
		return true
	}

	if result.Trend == models.TREND_DOWN && shortInterval.Trend == models.TREND_DOWN && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) == models.TREND_DOWN {
		ignoredReason = "first trend down and seconddown"
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP {
		ignoredReason = "first trend up and second not up"
		return true
	}

	if isContaineBearishEngulfing(result) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if result.Position == models.ABOVE_UPPER && (result.Trend != models.TREND_UP || result.AllTrend.ShortTrend != models.TREND_UP) {
		ignoredReason = "position above upper and not trend up"
		return true
	}

	isMidIntervalTrendNotUp := CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-4:]) != models.TREND_UP
	isLongIntervalTrendNotUp := CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP
	if shortInterval.Position == models.ABOVE_UPPER || midInterval.Position == models.ABOVE_UPPER || result.Position == models.ABOVE_UPPER || masterTrend == models.TREND_DOWN || masterMidTrend == models.TREND_DOWN {
		if shortInterval.Trend != models.TREND_UP || isMidIntervalTrendNotUp || isLongIntervalTrendNotUp {
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

	if midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
		if result.Position == models.ABOVE_SMA {
			longLastBand := result.Bands[len(result.Bands)-1]
			midLastBand := midInterval.Bands[len(result.Bands)-1]
			shortLastBand := shortInterval.Bands[len(result.Bands)-1]

			percentLong := (float32(longLastBand.Upper) - longLastBand.Candle.Close) / longLastBand.Candle.Close * float32(100)
			percentMid := (float32(midLastBand.Upper) - midLastBand.Candle.Close) / midLastBand.Candle.Close * float32(100)
			percentshort := (float32(shortLastBand.Upper) - shortLastBand.Candle.Close) / shortLastBand.Candle.Close * float32(100)

			if (percentLong < 3.3 && percentMid < 3.3 && percentshort < 3.3) && (shortInterval.Trend != models.TREND_UP || isMidIntervalTrendNotUp || isLongIntervalTrendNotUp) {
				ignoredReason = "all band bellow 3.1 from upper or not up trend"
				return true
			}
		}

		if isLastBandOrPreviousBandCrossSMA(result.Bands) {
			ignoredReason = "short and mid above sma but long interval cross sma"
			return true
		}

		if result.AllTrend.SecondTrend == models.TREND_DOWN {
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

	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open > float32(lastBand.Upper) {
		ignoredReason = "open close above upper"
		return true
	}

	if shortInterval.Position == models.ABOVE_UPPER && midInterval.Position == models.ABOVE_UPPER && result.Position == models.ABOVE_UPPER {
		highestIndex := getHighestHightIndex(result.Bands)
		if highestIndex == len(result.Bands)-1 {
			ignoredReason = "all interval above upper and new hight created"
			return true
		}
	}

	return false
}

func IsIgnoredMasterDown(result, midInterval, masterCoin *models.BandResult, checkingTime time.Time) bool {
	minPercentChanges := 2
	midLowestIndex := getLowestIndexSecond(midInterval.Bands)
	if !isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-4:]) {
		if midLowestIndex < len(midInterval.Bands)-4 {
			ignoredReason = "mid interval is not in lower"
			return true
		}

		if isNotInLower(result.Bands, false) {
			ignoredReason = "is not in lower"
			return true
		}

		minPercentChanges = 3
	} else {
		if isNotInLower(result.Bands, true) {
			ignoredReason = "is not in lower"
			return true
		}
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

	isHammer := IsHammer(result.Bands[len(midInterval.Bands)-3:]) || IsHammer(result.Bands[len(midInterval.Bands)-3:len(midInterval.Bands)-1])
	if CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-3:]) != models.TREND_UP && (CountSquentialUpBand(midInterval.Bands[len(midInterval.Bands)-3:]) < 2 || isHammer) && midLowestIndex != len(midInterval.Bands)-1 {
		ignoredReason = "last three band not up"
		return true
	}

	if midInterval.Trend != models.TREND_DOWN || midInterval.AllTrend.FirstTrend == models.TREND_UP {
		ignoredReason = "trend not down or first trend is up"
		return true
	}

	if isCrossLowerWhenSignificanDown(midInterval) {
		ignoredReason = "mid interval cross lower on significan down"
		return true
	}

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
		ignoredReason = "previous band hit SMA and current band bellow SMA"
		return true
	}

	if lastBand.Candle.Close < secondLastBand.Candle.Close {
		ignoredReason = "close below previous band"
		return true
	}

	// kalo nemu 1 lg case baru dienable
	// if checkingTime.Minute() < 18 {
	// 	ignoredReason = "skip on mid interval time change"
	// 	return true
	// }

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
	trend := CalculateTrends(bands[:len(bands)-1])
	shortTrend := CalculateTrendShort(bands[len(bands)-4:])
	return trend == models.TREND_DOWN && shortTrend == models.TREND_UP
}

func isCrossLowerWhenSignificanDown(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open < float32(lastBand.Lower) && lastBand.Candle.Close > float32(lastBand.Lower) {
		return result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent < 50
	}

	return false
}

func isLastBandCrossUpperAndPreviousBandNot(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	if lastBand.Candle.Open < lastBand.Candle.Close {
		secondLastBand := bands[len(bands)-2]
		return !(secondLastBand.Candle.Open < secondLastBand.Candle.Close)
	}
	return false
}

func isHasCrossUpper(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open < float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
			return true
		}
	}
	return false
}

func isNotInLower(bands []models.Band, skipped bool) bool {
	lowestIndex := getLowestIndex(bands)
	if !isHasCrossLower(bands[len(bands)-10:]) {
		if isHasCrossLower(bands[len(bands)-20:]) || skipped {
			return lowestIndex < len(bands)-10
		}
		return true
	}

	return false
}

func isHasCrossLower(bands []models.Band) bool {
	crossLowerBand := false
	for _, data := range bands {
		if data.Candle.Low < float32(data.Lower) {
			crossLowerBand = true
			break
		}
	}

	return crossLowerBand
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
	isLastBandCrossSMA := lastBand.Candle.Open <= float32(lastBand.SMA) && lastBand.Candle.Hight >= float32(lastBand.SMA)

	return isLastBandCrossSMA || isSecondLastBandCrossSMA
}

func isUpThreeOnMidIntervalChange(result *models.BandResult, requestTime time.Time) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	isCrossSMA := lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA)
	isCrossUpper := lastBand.Candle.Low < float32(lastBand.Upper) && lastBand.Candle.Hight > float32(lastBand.Upper)
	if isCrossSMA || isCrossUpper {
		if CountSquentialUpBand(result.Bands[len(result.Bands)-4:]) >= 3 && requestTime.Minute() < 17 {
			return true
		}
	}

	return false
}

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
	if index > len(result.Bands)-5 {
		index = len(result.Bands) - 5
	}
	lastDataFromHight := result.Bands[index:]
	if isHeighestOnHalfEndAndAboveUpper(result) && CalculateTrendShort(lastDataFromHight) != models.TREND_UP {
		return true
	}

	return false
}

func isHeighestOnHalfEndAndAboveUpper(result *models.BandResult) bool {
	hiIndex := getHighestIndex(result.Bands)
	if hiIndex >= len(result.Bands)/2 {
		return result.Bands[hiIndex].Candle.Close > float32(result.Bands[hiIndex].Upper)
	}

	return false
}

func isContaineBearishEngulfing(result *models.BandResult) bool {
	hiIndex := len(result.Bands) - (len(result.Bands) / 4)
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

	lastThreeBand := result.Bands[len(result.Bands)-3:]
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
				trend := CalculateTrends(result.Bands[15:])
				return trend != models.TREND_UP
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

func afterUpThenDown(result *models.BandResult) bool {
	sizeData := len(result.Bands)
	bands := result.Bands[sizeData/5:]
	if result.Position == models.ABOVE_SMA {
		higestIndex := getIndexHigestCrossUpper(bands)
		if higestIndex >= 0 && higestIndex < len(bands)-4 {
			trend := CalculateTrendsDetail(bands[higestIndex:])
			return trend.FirstTrend == models.TREND_DOWN
		}
	}

	return false
}

func getIndexHigestCrossUpper(bands []models.Band) int {
	higestIndex := -1
	lastBand := bands[len(bands)-1]
	for i, band := range bands {
		if band.Candle.Close > float32(band.Upper) || band.Candle.Close > lastBand.Candle.Close {
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

func GetIgnoredReason() string {
	return ignoredReason
}
