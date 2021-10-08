package analysis

import (
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
			ignoredReason = "count up when trend up"
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
	if result.Position == models.ABOVE_UPPER && secondLastBand.Candle.Close < float32(secondLastBand.Upper) {
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
	if lastBand.Candle.Open > float32(lastBand.Upper) {
		ignoredReason = "open close above upper"
		return true
	}

	if isUpSignificanAndNotUp(result) {
		ignoredReason = "after up significan and trend not up"
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

	if result.Trend == models.TREND_DOWN && shortInterval.Trend != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
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

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-4:]) < 3 {
		ignoredReason = "count up"
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		if result.Trend != models.TREND_UP || result.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "position above upper trend not up"
			return true
		}

		secondLastBand := result.Bands[len(result.Bands)-2]
		if CountUpBand(result.Bands[len(result.Bands)-4:]) < 3 || (secondLastBand.Candle.Close > secondLastBand.Candle.Open && secondLastBand.Candle.Close > float32(secondLastBand.Upper)) {
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

	if isLastBandOrPreviousBandCrossSMA(shortInterval.Bands) {
		if shortInterval.Trend == models.TREND_DOWN || result.Trend == models.TREND_DOWN {
			ignoredReason = "short interval cross  sma and mid interval or short interval trend down"
			return true
		}
	}

	return false
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

func IsIgnoredLongInterval(result *models.BandResult, shortInterval *models.BandResult, midInterval *models.BandResult, requestTime time.Time) bool {
	if isInAboveUpperBandAndDownTrend(result) {
		ignoredReason = "isInAboveUpperBandAndDownTrend"
		return true
	}

	if result.Trend == models.TREND_DOWN && shortInterval.Trend == models.TREND_DOWN && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP {
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
	if shortInterval.Position == models.ABOVE_UPPER || midInterval.Position == models.ABOVE_UPPER || result.Position == models.ABOVE_UPPER {
		if shortInterval.Trend != models.TREND_UP || isMidIntervalTrendNotUp || isLongIntervalTrendNotUp {
			ignoredReason = "when above upper and trend not up"
			return true
		}

		lenData := len(result.Bands)
		hight := getHighestIndex(result.Bands[lenData-lenData/3:])
		low := getLowestIndex(result.Bands[lenData-lenData/3:])
		difference := hight - low
		percent := float32(difference) / float32(low) * 100
		if percent > 15 {
			ignoredReason = "up more than 15"
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	if result.Position == models.ABOVE_SMA && midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
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

	if isLastBandCrossUpperAndPreviousBandNot(shortInterval.Bands) {
		if isLastBandCrossUpperAndPreviousBandNot(midInterval.Bands) {
			if isLastBandCrossUpperAndPreviousBandNot(result.Bands) {
				ignoredReason = "band above upper and just one in all interval, skip"
				return true
			}
		}
	}

	return false
}

func IsIgnoredMasterDown(result, midInterval, masterCoin *models.BandResult, checkingTime time.Time) bool {
	if isNotInLower(result.Bands, false) {
		ignoredReason = "is not in lower"
		return true
	}

	if isNotInLower(midInterval.Bands, true) {
		ignoredReason = "min interval is not in lower"
		return true
	}

	if midInterval.Direction == BAND_DOWN {
		ignoredReason = "mid interval band down"
		return true
	}

	if CalculateTrendShort(masterCoin.Bands[len(masterCoin.Bands)-4:]) != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		ignoredReason = "short trend not up"
		return true
	}

	if CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-3:]) != models.TREND_UP && CountSquentialUpBand(midInterval.Bands[len(midInterval.Bands)-3:]) < 2 {
		ignoredReason = "last three band not up"
		return true
	}

	lastBand := result.Bands[len(result.Bands)-1]
	marginFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if marginFromUpper < 3 {
		ignoredReason = "margin from upper is bellow 3"
		return true
	}

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-5:]) < 4 {
		ignoredReason = "count up bellow 2"
		return true
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

func isNotInLower(bands []models.Band, skipCrossLowCheck bool) bool {
	lowestIndex := getLowestIndex(bands)
	if IsLastCandleNotCrossLower(bands, 5) || skipCrossLowCheck {
		if !IsLastCandleNotCrossLower(bands, 15) || skipCrossLowCheck {
			return lowestIndex < len(bands)-6
		}
		return true
	}

	return false
}

func IsLastCandleNotCrossLower(bands []models.Band, number int) bool {
	lastFour := bands[len(bands)-number:]

	crossLowerBand := false
	for _, data := range lastFour {
		if data.Candle.Low < float32(data.Lower) {
			crossLowerBand = true
			break
		}
	}

	return !crossLowerBand
}

func isUpThreeOnMidIntervalChange(result *models.BandResult, requestTime time.Time) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	isCrossSMA := lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA)
	isCrossUpper := lastBand.Candle.Low < float32(lastBand.Upper) && lastBand.Candle.Hight > float32(lastBand.Upper)
	if result.Position == models.BELOW_SMA || isCrossSMA || isCrossUpper {
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
	if index > 5 {
		index = 5
	}
	lastDataFromHight := result.Bands[len(result.Bands)-index:]
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
		if bands[hiIndex].Candle.Close < band.Candle.Close {
			hiIndex = i
		}
	}

	return hiIndex
}

func getLowestIndex(bands []models.Band) int {
	hiIndex := 0
	for i, band := range bands {
		if bands[hiIndex].Candle.Close < band.Candle.Close {
			hiIndex = i
		}
	}

	return hiIndex
}

func whenHeightTripleAverage(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open
	var totalHeight float32 = 0
	for _, band := range result.Bands[len(result.Bands)-6 : len(result.Bands)-1] {
		if band.Candle.Open > band.Candle.Close {
			totalHeight += band.Candle.Open - band.Candle.Close
		} else {
			totalHeight += band.Candle.Close - band.Candle.Open
		}
	}
	average := totalHeight / float32(len(result.Bands)-1)
	percent := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100

	return lastBandHeight/average >= 3.5 && percent > 2.5
}

func lastBandHeadDoubleBody(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Close > lastBand.Candle.Open {
		head := lastBand.Candle.Hight - lastBand.Candle.Close
		body := lastBand.Candle.Close - lastBand.Candle.Open
		return head > body*2.5
	}

	return false
}

func ignored(result, masterCoin *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Low <= float32(lastBand.SMA) && lastBand.Candle.Hight >= float32(lastBand.Upper) {
		ignoredReason = "up from bellow sma to upper"
		return true
	}

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent <= 2.5 {
		ignoredReason = "hight and low bellow 2.5"
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
	if result.AllTrend.SecondTrendPercent < 40 {
		mid := len(result.Bands) / 2
		indexDoubleBody := getIndexBandDoubleLong(result.Bands[len(result.Bands)-mid:])
		if indexDoubleBody > -1 {
			realIndex := len(result.Bands)%2 + mid + indexDoubleBody
			if len(result.Bands)-(mid+indexDoubleBody) > 4 {
				trend := CalculateTrends(result.Bands[realIndex:])
				return trend != models.TREND_UP
			}
		}
	}

	return false
}

func getIndexBandDoubleLong(bands []models.Band) int {
	var total float32
	for _, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			total += band.Candle.Close - band.Candle.Open
		} else {
			total += band.Candle.Open - band.Candle.Close
		}
	}
	average := total / float32(len(bands))
	longestIndex := -1
	for i, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			bandLong := band.Candle.Close - band.Candle.Open
			if bandLong > average*2 {
				longestIndex = i
			}
		}
	}

	return longestIndex
}

func GetIgnoredReason() string {
	return ignoredReason
}
