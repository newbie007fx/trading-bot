package analysis

import (
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""

func IsIgnored(result, masterCoin *models.BandResult) bool {
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

	if isUpMoreThanThreeOnDownBellowSMA(result) {
		ignoredReason = "isUpMoreThanThreeOnDownBellowSMA"
		return true
	}

	return ignored(result, masterCoin)
}

func IsIgnoredMidInterval(result *models.BandResult, shortInterval *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) && result.Direction == BAND_DOWN {
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

	if shortInterval.Position == models.ABOVE_UPPER && (CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP || shortInterval.Trend != models.TREND_UP) {
		ignoredReason = "short interval above upper and mid trend down or trend not up"
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		ignoredReason = "position above upper"
		return true
	}

	lastBand := result.Bands[len(result.Bands)-1]
	if shortInterval.Position == models.ABOVE_UPPER && lastBand.Candle.Hight > float32(lastBand.Upper) {
		ignoredReason = "above upper and mid interval height above upper"
		return true
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	if isContaineBearishEngulfing(result) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	return false
}

func IsIgnoredLongInterval(result *models.BandResult, shortInterval *models.BandResult, midInterval *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) && result.Direction == BAND_DOWN {
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

	if result.Position == models.ABOVE_UPPER {
		ignoredReason = "position above upper"
		return true
	}

	isMidIntervalTrendNotUp := CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-4:]) != models.TREND_UP
	isLongIntervalTrendNotUp := CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP
	if shortInterval.Position == models.ABOVE_UPPER {
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

		if (percentLong < 3.1 && percentMid < 3.1 && percentshort < 3.1) && (shortInterval.Trend != models.TREND_UP || isMidIntervalTrendNotUp || isLongIntervalTrendNotUp) {
			ignoredReason = "all band bellow 3.1 from upper or not up trend"
			return true
		}
	}

	return false
}

func IsIgnoredMasterDown(result, masterCoin *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	marginFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if marginFromUpper < 3.5 {
		return true
	}

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 {
		return true
	}

	if CalculateTrendShort(masterCoin.Bands[len(masterCoin.Bands)-4:]) != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		return true
	}

	return false
}

func isUpMoreThanThreeOnDownBellowSMA(result *models.BandResult) bool {
	currentTime := time.Now()
	if result.Position == models.BELOW_SMA {
		if CountSquentialUpBand(result.Bands[len(result.Bands)-4:]) > 3 && currentTime.Minute() < 17 {
			return true
		}
	}

	return false
}

func isPosititionBellowUpperMarginBellowThreshold(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Close < float32(lastBand.Upper) && (result.AllTrend.FirstTrend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN) {
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
	for _, band := range result.Bands[:len(result.Bands)-1] {
		if band.Candle.Open > band.Candle.Close {
			totalHeight += band.Candle.Open - band.Candle.Close
		} else {
			totalHeight += band.Candle.Close - band.Candle.Open
		}
	}
	average := totalHeight / float32(len(result.Bands)-1)

	return lastBandHeight/average >= 3
}

func lastBandHeadDoubleBody(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Close > lastBand.Candle.Open {
		head := lastBand.Candle.Hight - lastBand.Candle.Close
		body := lastBand.Candle.Close - lastBand.Candle.Open
		return head > body*2
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
	if percent < 2 {
		ignoredReason = "hight and low bellow 2"
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

func GetIgnoredReason() string {
	return ignoredReason
}
