package analysis

import (
	"log"
	"telebot-trading/app/models"
	"time"
)

func IsIgnored(result, masterCoin *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) {
		return true
	}

	if lastBandHeadDoubleBody(result) {
		return true
	}

	if isContaineBearishEngulfing(result) {
		return true
	}

	// if whenHeadCrossBandAndMasterDown(result, masterCoin) {
	// 	return true
	// }

	if isPosititionBellowUpperMarginBellowThreshold(result) {
		return true
	}

	if whenHeightTripleAverage(result) {
		return true
	}

	if lastFourCandleNotUpTrend(result.Bands) {
		return true
	}

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && CountUpBand(result.Bands[len(result.Bands)-4:]) < 3 {
		return true
	}

	if isUpMoreThanThreeOnDownBellowSMA(result) {
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		return true
	}

	return ignored(result, masterCoin)
}

func IsIgnoredMidInterval(result *models.BandResult, shortInterval *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) && result.Direction == BAND_DOWN {
		return true
	}

	if lastBandHeadDoubleBody(result) {
		return true
	}

	if result.Trend == models.TREND_DOWN && shortInterval.Trend != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_SIDEWAY {
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP {
		return true
	}

	if CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 {
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		return true
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		return true
	}

	return isContaineBearishEngulfing(result)
}

func IsIgnoredLongInterval(result *models.BandResult, shortInterval *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) && result.Direction == BAND_DOWN {
		return true
	}

	if lastBandHeadDoubleBody(result) {
		return true
	}

	if result.Trend == models.TREND_DOWN && shortInterval.Trend == models.TREND_DOWN && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP {
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend != models.TREND_UP && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) != models.TREND_UP {
		return true
	}

	if isContaineBearishEngulfing(result) {
		return true
	}

	if result.Position == models.ABOVE_UPPER {
		return true
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.Position == models.ABOVE_SMA {
		lenData := len(result.Bands)
		hight := getHighestIndex(result.Bands[lenData-lenData/3:])
		low := getLowestIndex(result.Bands[lenData-lenData/3:])
		difference := hight - low
		percent := float32(difference) / float32(low) * 100
		if percent > 15 {
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		return true
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
	if lastBand.Candle.Close < float32(lastBand.Upper) && result.Trend == models.TREND_DOWN {
		margin := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100

		return margin < 1.5
	}

	return false
}

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	lastFiveData := result.Bands[len(result.Bands)-5:]
	if isHeighestOnHalfEndAndAboveUpper(result) && CalculateTrendShort(lastFiveData) == models.TREND_DOWN {
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
		log.Println("reset to 0 with criteria 1")
		return true
	}

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent < 2 {
		log.Println("reset to 0 with criteria 2")
		return true
	}

	return false
}

// func whenHeadCrossBandAndMasterDown(result, masterCoin *models.BandResult) bool {
// 	lastBand := result.Bands[len(result.Bands)-1]
// 	crossSMA := lastBand.Candle.Close < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA)
// 	crossUpper := lastBand.Candle.Close < float32(lastBand.Upper) && lastBand.Candle.Hight > float32(lastBand.Upper)
// 	if (crossSMA || crossUpper) && CalculateTrendShort(masterCoin.Bands[len(masterCoin.Bands)-5:]) == models.TREND_DOWN {
// 		return true
// 	}

// 	return false
// }

func lastFourCandleNotUpTrend(bands []models.Band) bool {
	return CalculateTrendShort(bands[len(bands)-4:]) != models.TREND_UP
}

func isTrendUpLastThreeBandHasDoji(result *models.BandResult) bool {
	if result.AllTrend.SecondTrend != models.TREND_DOWN {
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
