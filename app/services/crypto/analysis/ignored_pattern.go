package analysis

import (
	"log"
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""
var skipped bool = false

func upperLowerReversal(result models.BandResult) bool {
	hightIndex := getHighestIndex(result.Bands)
	if hightIndex > len(result.Bands)/2 && result.Bands[hightIndex].Candle.Hight > float32(result.Bands[hightIndex].Upper) {
		lowIndex := getLowestIndex(result.Bands)
		if lowIndex > hightIndex && lowIndex < len(result.Bands)-1 {
			if !isHasCrossSMA(result.Bands[lowIndex:len(result.Bands)-1], false) && CalculateTrendShort(result.Bands[lowIndex:]) == models.TREND_UP {
				higestPrice := result.Bands[hightIndex].Candle.Hight
				percent := (higestPrice - result.CurrentPrice) / result.CurrentPrice * 100
				return percent > 3
			}
		}
	}

	return false
}

func upperLowerMarginBelow3(result models.BandResult) bool {
	hightIndex := getHighestIndex(result.Bands)
	if hightIndex >= len(result.Bands)/2-3 && result.Bands[hightIndex].Candle.Hight > float32(result.Bands[hightIndex].Upper) {
		lowIndex := getLowestIndex(result.Bands)
		if lowIndex > hightIndex && lowIndex < len(result.Bands)-1 {
			if !isHasCrossSMA(result.Bands[lowIndex:len(result.Bands)-1], false) && CalculateTrendShort(result.Bands[lowIndex:]) == models.TREND_UP {
				higestPrice := result.Bands[hightIndex].Candle.Hight
				percent := (higestPrice - result.CurrentPrice) / result.CurrentPrice * 100
				return percent < 3
			}
		}
	}

	return false
}

func upperSmaMarginBelow3(result models.BandResult) bool {
	hightIndex := getHighestIndex(result.Bands)
	if hightIndex >= len(result.Bands)/2-3 && result.Bands[hightIndex].Candle.Hight > float32(result.Bands[hightIndex].Upper) {
		if isHasCrossSMA(result.Bands[hightIndex:], false) {
			higestPrice := result.Bands[hightIndex].Candle.Hight
			percent := (higestPrice - result.CurrentPrice) / result.CurrentPrice * 100
			return percent < 3
		}
	}

	return false
}

func percentFromHigest(bands []models.Band) float32 {
	higestIndex := getHighestHightIndex(bands)
	lastBand := bands[len(bands)-1]
	if lastBand.Candle.Close < float32(lastBand.SMA) {
		if bands[higestIndex].Candle.Hight < float32(bands[higestIndex].SMA) && countAboveSMA(bands) == 0 {
			higestPrice := bands[higestIndex].Candle.Hight
			return (higestPrice - lastBand.Candle.Close) / lastBand.Candle.Close * 100
		}
	}

	return 0
}

func downFromUpper(result models.BandResult) bool {
	hightIndex := getHighestIndex(result.Bands[len(result.Bands)/2:])
	index := hightIndex + len(result.Bands)/2
	lastBand := result.Bands[len(result.Bands)-1]
	percentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if result.Position == models.ABOVE_SMA && result.Bands[index].Candle.Hight > float32(result.Bands[index].Upper) {
		if countBelowSMA(result.Bands[index:], false) == 0 || percentFromUpper < 3 {
			return CalculateTrendShort(result.Bands[index:]) == models.TREND_DOWN
		}
	}

	return false
}

func downFromUpperAboveSMA(result models.BandResult) bool {
	hightIndex := getHighestIndex(result.Bands)
	if hightIndex > len(result.Bands)/2 && result.Position == models.ABOVE_SMA {
		if result.Bands[hightIndex].Candle.Close > float32(result.Bands[hightIndex].Upper) {
			return CalculateTrendShort(result.Bands[hightIndex:]) == models.TREND_DOWN && !isHasCrossSMA(result.Bands[hightIndex:], false)
		}
	}

	return false
}

func downFromUpperBelowSMA(result models.BandResult) bool {
	hightIndex := getHighestIndex(result.Bands[len(result.Bands)/2:])
	index := hightIndex + len(result.Bands)/2
	if result.Position == models.BELOW_SMA && result.Bands[index].Candle.Hight > float32(result.Bands[index].Upper) {
		return CalculateTrendShort(result.Bands[index:]) == models.TREND_DOWN
	}

	return false
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

func isLastBandCrossUpperAndPreviousBandNot(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	if lastBand.Candle.Open < lastBand.Candle.Close {
		secondLastBand := bands[len(bands)-2]
		return secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Low > float32(secondLastBand.Lower)
	}
	return false
}

func isHasCrossUpper(bands []models.Band, withHead bool) bool {
	for _, band := range bands {
		if band.Candle.Close < band.Candle.Open {
			continue
		}

		if withHead {
			if band.Candle.Open < float32(band.Upper) && band.Candle.Hight >= float32(band.Upper) {
				return true
			}
			if band.Candle.Open > band.Candle.Close {
				if band.Candle.Close < float32(band.Upper) && band.Candle.Hight >= float32(band.Upper) {
					return true
				}
			}
		} else {
			if band.Candle.Open <= float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
				return true
			}
		}
	}
	return false
}

func isHasBandDownFromUpper(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open > float32(band.Upper) && band.Candle.Close < band.Candle.Open {
			return true
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

func isHasCrossLower(bands []models.Band, bodyOnly bool) bool {
	crossLowerBand := false
	for _, data := range bands {
		if bodyOnly {
			if (data.Candle.Close < float32(data.Lower) && data.Candle.Open > float32(data.Lower)) || (data.Candle.Open < float32(data.Lower) && data.Candle.Close > float32(data.Lower)) {
				crossLowerBand = true
				break
			}
		} else {
			if data.Candle.Low < float32(data.Lower) && (data.Candle.Close > float32(data.Lower) || data.Candle.Open > float32(data.Lower)) {
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
		if data.Candle.Open <= float32(data.Upper) && data.Candle.Hight > float32(data.Upper) && (data.Candle.Close > data.Candle.Open) {
			if !isHasBadBand([]models.Band{data}) {
				count++
			}
		}
	}

	return count
}

func countHighAboveUpper(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Hight > float32(data.Upper) && (data.Candle.Open < float32(data.Upper) || data.Candle.Close < float32(data.Upper)) {
			count++
		}
	}

	return count
}

func countCrossUpperOnBody(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Open <= float32(data.Upper) && data.Candle.Close >= float32(data.Upper) && (data.Candle.Close > data.Candle.Open) {
			if !isHasBadBand([]models.Band{data}) {
				count++
			}
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

func countBelowLower(bands []models.Band, strict bool) int {
	count := 0
	for _, data := range bands {
		if strict {
			if data.Candle.Hight < float32(data.Lower) && data.Candle.Low < float32(data.Lower) {
				count++
			}
		} else {
			if data.Candle.Close < float32(data.Lower) && data.Candle.Open < float32(data.Lower) {
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

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	index := getHighestIndex(result.Bands)
	lastBand := result.Bands[len(result.Bands)-1]
	if index == len(result.Bands)-1 || lastBand.Candle.Close < float32(lastBand.SMA) {
		return false
	}

	if isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
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

func ignored(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open <= float32(lastBand.SMA) && lastBand.Candle.Hight >= float32(lastBand.Upper) {
		ignoredReason = "up from bellow sma to upper"
		return true
	}

	return false
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
		if IsHammer(band) {
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

func getLastIndexCrossLower(bands []models.Band) int {
	for i := len(bands) - 1; i > 0; i-- {
		if bands[i].Candle.Low < float32(bands[i].Lower) {
			return i
		}
	}

	return -1
}

func GetIgnoredReason() string {
	return ignoredReason
}

func isBandHeadDoubleBody(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			head := band.Candle.Hight - band.Candle.Close
			body := band.Candle.Close - band.Candle.Open
			if head > body*1.8 {
				return true
			}
		} else {
			head := band.Candle.Hight - band.Candle.Open
			body := band.Candle.Open - band.Candle.Close
			if head > body*1.5 {
				return true
			}
		}
	}

	return false
}

func CountUpBand(bands []models.Band) int {
	counter := 0
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			difference := band.Candle.Close - band.Candle.Open
			if (difference / band.Candle.Open * 100) > 0.1 {
				counter++
			}
		}
	}

	return counter
}

func isContainHeadMoreThanBody(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			head := band.Candle.Hight - band.Candle.Close
			body := band.Candle.Close - band.Candle.Open

			if head > body {
				return true
			}
		}
	}

	return false
}

func countHeadMoreThanBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		head := band.Candle.Hight - band.Candle.Close
		body := band.Candle.Close - band.Candle.Open
		if band.Candle.Open > band.Candle.Close {
			head = band.Candle.Hight - band.Candle.Open
			body = band.Candle.Open - band.Candle.Close
		}

		if head > body {
			count++
		}
	}

	return count
}

func IgnoredOnUpTrendShort(shortInterval models.BandResult) bool {
	if skipped {
		return false
	}

	if isHasOpenCloseAboveUpper(shortInterval.Bands[len(shortInterval.Bands)-2:]) {
		ignoredReason = "contain open close above upper"
		return true
	}

	if shortInterval.Position == models.ABOVE_UPPER && isContainHeadMoreThanBody(shortInterval.Bands[len(shortInterval.Bands)-1:]) {
		ignoredReason = "above upper and head more than body"
		return true
	}

	if shortInterval.PriceChanges < 3 {
		ignoredReason = "price change below 3"
		return true
	}

	if isHasBadBand(shortInterval.Bands[len(shortInterval.Bands)-1:]) {
		ignoredReason = "bad bands"
		return true
	}

	return false
}

func IgnoredOnUpTrendMid(midInterval, shortInterval models.BandResult) bool {
	if skipped {
		return false
	}

	bandLen := len(midInterval.Bands)
	lastBand := midInterval.Bands[bandLen-1]
	midPercentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-2:len(midInterval.Bands)-1]) && (countHeadMoreThanBody(midInterval.Bands[bandLen-2:bandLen-1]) == 1 || countDownBand(midInterval.Bands[bandLen-2:]) == 1) {
		ignoredReason = "contain open close above upper"
		return true
	}

	higestIndex := getHighestIndex(midInterval.Bands[len(midInterval.Bands)-4:])
	if higestIndex != len(midInterval.Bands[len(midInterval.Bands)-4:])-1 {
		ignoredReason = "not in higest"
		return true
	}

	higestHightIndex := getHighestHightIndex(midInterval.Bands[len(midInterval.Bands)-4:])
	if higestHightIndex < len(midInterval.Bands[len(midInterval.Bands)-4:])-2 {
		ignoredReason = "previous band have higest high"
		return true
	}

	if higestHightIndex != len(midInterval.Bands[len(midInterval.Bands)-4:])-1 && countDownBand(midInterval.Bands[len(midInterval.Bands)-2:]) == 1 {
		ignoredReason = "previous band have higest high and previous band down"
		return true
	}

	if midInterval.PriceChanges < 3 {
		ignoredReason = "price changes below 5"
		return true
	}

	if countDownBand(midInterval.Bands[len(midInterval.Bands)-4:]) > 1 && countDownBand(midInterval.Bands[len(midInterval.Bands)-2:]) == 1 {
		if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], 5) != 1 && !(isHasCrossSMA(midInterval.Bands[bandLen-1:], false) && countCrossSMA(midInterval.Bands[bandLen-4:]) == 1) {
			ignoredReason = "up down"
			return true
		}
	}

	if countDownBand(midInterval.Bands[len(midInterval.Bands)-2:]) == 1 && CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-5:len(midInterval.Bands)-1]) != models.TREND_UP {
		if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], 5) != 1 && !(isHasCrossSMA(midInterval.Bands[bandLen-1:], false) && countCrossSMA(midInterval.Bands[bandLen-4:]) == 1) {
			ignoredReason = "contain head more thank body and previous band down"
			return true
		}
	}

	if midInterval.Position == models.ABOVE_SMA || countCrossUpper(midInterval.Bands[bandLen-5:]) == 1 {
		if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], 3) == 0 && countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-5:], 3) == 0 {
			ignoredReason = "above sma and no significant band change"
			return true
		}
	}

	if midInterval.Position == models.ABOVE_SMA && midPercentFromUpper < 5 {
		if midInterval.AllTrend.SecondTrend != models.TREND_UP && countCrossUpper(midInterval.Bands[bandLen-5:]) == 0 {
			if isHasBadBand(shortInterval.Bands[bandLen-2:]) || isHasBadBand(midInterval.Bands[bandLen-2:]) {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen/2:]) == 1 {
					ignoredReason = "above sma, second trend up, short cross upper"
					return true
				}
			}
		}
	}

	if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) > 2 || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-1:]) {
		if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-3:]) {
			ignoredReason = "short contain open close above upper"
			return true
		}
	}

	if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
		if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
			ignoredReason = "above upper and just one and short upper head more than upper body"
			return true
		}
	}

	return false
}

func IgnoredOnUpTrendLong(longInterval, midInterval, shortInterval models.BandResult, checkTime time.Time) bool {
	if skipped {
		return false
	} else {
		ignoredReason = "skip reguler ignore checking"
		return true
	}
}

func isHasBadBand(bands []models.Band) bool {
	if countDownBand(bands) > 0 || countAboveUpper(bands) > 0 || isTailMoreThan(bands[len(bands)-1], 50) || countHeadMoreThanBody(bands) > 0 {
		return true
	}

	return false
}

func countBadBands(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isHasBadBand([]models.Band{band}) {
			count++
		}
	}

	return count
}

func longSignificanUpAndJustOne(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	lastPercent := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100
	var threshold float32 = 2.5
	if lastPercent >= 5.5 {
		threshold = lastPercent / 2
	}

	if lastPercent > threshold && countBandPercentChangesMoreThan(bands[len(bands)-4:], threshold-(threshold/8)) == 1 {
		return true
	}

	return false
}

func bandDoublePreviousHigh(bands []models.Band, multiplier float32) bool {
	lastBand := bands[len(bands)-1]
	lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open
	if lastBandHeight < 0 {
		return false
	}

	var higestHeigh float32 = 0
	for _, band := range bands[len(bands)-5 : len(bands)-1] {
		high := band.Candle.Close - band.Candle.Open
		if high > higestHeigh {
			higestHeigh = high
		}
	}

	return higestHeigh*multiplier < lastBandHeight
}

func isLastBandHigestBody(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	var lastHighBody float32 = lastBand.Candle.Close - lastBand.Candle.Open

	isHighestBody := true
	var higestHeigh float32 = 0
	for _, band := range bands[len(bands)-4 : len(bands)-1] {
		high := band.Candle.Close - band.Candle.Open
		if high > lastHighBody {
			isHighestBody = false
		}

		if higestHeigh < band.Candle.Hight {
			higestHeigh = band.Candle.Hight
		}
	}

	if !isHighestBody {
		return higestHeigh < lastBand.Candle.Close
	}

	return true
}

func isHeadPercentMoreThan(band models.Band, percent float32) bool {
	head := float32(band.Candle.Hight) - band.Candle.Close
	body := band.Candle.Close - band.Candle.Open
	log.Println(head / body * 100)
	return head/body*100 > percent
}

func isAllAboveUpperAndJustOne(short, mid, long models.BandResult) bool {
	bandLen := len(short.Bands)
	if short.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER {
		if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countCrossUpperOnBody(long.Bands[bandLen-4:]) == 1 {
			return true
		}
	}

	return true
}

func isAllIntervalHasUpTrend(short, mid, long models.BandResult) bool {
	if short.AllTrend.FirstTrend == models.TREND_UP || short.AllTrend.SecondTrend == models.TREND_UP {
		if mid.AllTrend.FirstTrend == models.TREND_UP || mid.AllTrend.SecondTrend == models.TREND_UP {
			if long.AllTrend.FirstTrend == models.TREND_UP || long.AllTrend.SecondTrend == models.TREND_UP {
				return true
			}
		}
	}

	return false
}

func allIntervalCrossUpperOnBodyMoreThanThresholdAndJustOne(short, mid, long models.BandResult, currentTime time.Time) bool {
	bandLen := len(long.Bands)
	shortLastBandPercent := (short.Bands[bandLen-1].Candle.Close - short.Bands[bandLen-1].Candle.Open) / short.Bands[bandLen-1].Candle.Open * 100
	if (bandDoublePreviousHigh(short.Bands, 2) || bandDoublePreviousHigh(mid.Bands, 2.2)) && mid.AllTrend.ShortTrend == models.TREND_UP {
		if longSignificanUpAndJustOne(mid.Bands) {

			longLastBand := long.Bands[bandLen-1]
			midLastBand := mid.Bands[bandLen-1]
			shortLastBand := short.Bands[bandLen-1]
			if longLastBand.Candle.Open < float32(longLastBand.SMA) && longLastBand.Candle.Close > float32(longLastBand.Upper) {
				if midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Close > float32(midLastBand.Upper) {
					if shortLastBand.Candle.Close > float32(shortLastBand.Upper) && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						log.Println("before gas 1")
						return false
					}
				}
			}

			if !isHeadPercentMoreThan(short.Bands[bandLen-1], 60) {
				if isAllIntervalHasUpTrend(short, mid, long) {
					if !(isHourInChangesLong(currentTime.Hour(), currentTime.Minute()) && isAllAboveUpperAndJustOne(short, mid, long)) {
						if countBandPercentChangesMoreThan(mid.Bands[bandLen/2:bandLen-1], 1) > 0 || countBandPercentChangesMoreThan(long.Bands[bandLen-5:bandLen-1], 1) > 0 {
							if !(long.AllTrend.SecondTrend == models.TREND_DOWN && isHasCrossLower(long.Bands[bandLen-2:bandLen-1], false) && long.Position != models.ABOVE_UPPER) {
								if !isHasBadBand(short.Bands[bandLen-1:]) && shortLastBandPercent > 2.5 && countBandPercentChangesMoreThan(short.Bands[:bandLen-1], 1) == 0 {
									if countCrossUpperOnBody(mid.Bands[:bandLen-1]) == 0 && countCrossUpperOnBody(short.Bands[:bandLen-1]) == 0 {
										if mid.Position == models.ABOVE_UPPER || long.Position == models.ABOVE_UPPER {
											if !(isHasCrossLower(short.Bands[bandLen-3:], false) && isHasCrossLower(mid.Bands[bandLen-2:], false) && long.Position == models.ABOVE_SMA && isHasCrossSMA(long.Bands[bandLen-1:], true) && countCrossSMA(long.Bands[bandLen-5:bandLen-1]) == 0 && countAboveSMA(long.Bands[bandLen-5:bandLen-1]) == 0) {
												if !(countBadBands(long.Bands[bandLen-4:]) > 2 && countBadBands(short.Bands[bandLen-4:]) > 2) {
													if !isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) && !isUpperHeadMoreThanUpperBody(long.Bands[bandLen-1]) {
														log.Println("gas wae 1")
														return true
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
				if countCrossUpperOnBody(mid.Bands[bandLen-4:]) > 1 && isHasBadBand(mid.Bands[bandLen-2:]) {
					log.Println("1")
					return false
				}
			}

			if isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) && currentTime.Minute() < 3 {
				log.Println("2")
				return false
			}

			if isHasOpenCloseAboveUpper(short.Bands[bandLen-1:]) {
				log.Println("3")
				return false
			}

			if isHasOpenCloseAboveUpper(mid.Bands[bandLen-4:]) {
				log.Println("4")
				return false
			}

			if long.Position == models.BELOW_SMA && long.AllTrend.ShortTrend != models.TREND_UP {
				log.Println("5")
				return false
			}

			if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 3) == 1 && getHigestPercentChangesIndex(short.Bands[len(short.Bands)-4:]) == 0 && countDownBand(short.Bands[len(short.Bands)-4:]) > 1 {
				log.Println("6")
				return false
			}

			if long.Bands[len(long.Bands)-1].Candle.Low < float32(long.Bands[len(long.Bands)-1].Lower) && long.Bands[len(long.Bands)-1].Candle.Close > float32(long.Bands[len(long.Bands)-1].Upper) {
				log.Println("7")
				return false
			}

			if isTailMoreThan(short.Bands[len(short.Bands)-1], 40) {
				log.Println("8")
				return false
			}

			if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 3) > 1 && short.Bands[len(short.Bands)-2].Candle.Open > short.Bands[len(short.Bands)-2].Candle.Close {
				log.Println("9")
				return false
			}

			if isHasBadBand(short.Bands[bandLen-1:]) {
				log.Println("10")
				return false
			}

			if countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) > 1 && countCrossUpperOnBody(short.Bands[len(short.Bands)-4:]) > 1 && isUpperHeadMoreThanUpperBody(mid.Bands[len(mid.Bands)-1]) {
				log.Println("11")
				return false
			}

			if isHasOpenCloseAboveUpper(mid.Bands[len(mid.Bands)-2:]) && getHigestPercentChangesIndex(short.Bands[len(short.Bands)-4:]) == 3 {
				if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 5) == 1 && countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 2) == 1 {
					log.Println("12")
					return false
				}
			}

			if isPreviousBandTripleHeigh(short.Bands) {
				if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) && countCrossUpper(short.Bands[bandLen-2:]) > 1 {
					log.Println("13.1")
					return false
				}

				if short.Position == models.ABOVE_UPPER && countCrossUpper(short.Bands[bandLen-2:]) == 1 {
					log.Println("13.2")
					return false
				}
			}

			if isHasOpenCloseAboveUpper(mid.Bands[len(mid.Bands)-2:]) && (isHasUpperHeadMoreThanUpperBody(short.Bands[len(short.Bands)-2:]) || isHasOpenCloseAboveUpper(short.Bands[len(short.Bands)-2:])) {
				log.Println("14")
				return false
			}

			if isHasBandDownFromUpper(long.Bands[len(long.Bands)-2:]) {
				if countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) == 1 {
					if isUpperHeadMoreThanUpperBody(short.Bands[len(short.Bands)-1]) {
						log.Println("15")
						return false
					}
				}
			}

			if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) && countCrossUpper(short.Bands[bandLen-2:]) > 1 {
				if !isHasCrossUpper(mid.Bands[bandLen-2:], false) && long.Position != models.ABOVE_UPPER {
					log.Println("16")
					return false
				}
			}

			if !isLastBandHigestBody(short.Bands) {
				log.Println("17")
				return false
			}

			if isHasCrossUpper(mid.Bands[bandLen-4:bandLen-1], true) && countBadBands(mid.Bands[bandLen-4:bandLen-1]) > 2 {
				log.Println("18")
				return false
			}

			minute := currentTime.Minute()
			hour := currentTime.Hour()
			if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 || (countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 0) {
				if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
					if isHourInChangesLong(hour, minute) {
						log.Println("19")
						return false
					}
				}
			}

			if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) && (countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 || countCrossUpperOnBody(long.Bands[bandLen-4:]) == 1) {
				if (isHourInChangesLong(hour, minute)) && (isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) || isHasOpenCloseAboveUpper(mid.Bands[bandLen-1:])) {
					log.Println("20")
					return false
				}
			}

			if countCrossUpperOnBody(short.Bands[bandLen-4:]) > 1 && isHasBandDownFromUpper(short.Bands[bandLen-4:]) {
				if mid.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend != models.TREND_UP {
					log.Println("21")
					return false
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_DOWN && long.Position == models.BELOW_SMA {
				if countCrossUpperOnBody(mid.Bands[bandLen-4:]) <= 1 {
					if countCrossUpperOnBody(short.Bands[bandLen-4:]) <= 1 || isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
						log.Println("22")
						return false
					}
				}
			}

			if long.Position == models.ABOVE_UPPER && isHasBandDownFromUpper(long.Bands[bandLen-4:]) {
				if (mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1) || !isHasCrossUpper(mid.Bands[bandLen-4:], false) {
					if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
						log.Println("23")
						return false
					}
				}
			}

			if countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countCrossUpperOnBody(long.Bands[bandLen-4:]) == 1 {
				if countCrossUpperOnBody(short.Bands[bandLen-4:]) > 1 && isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
					log.Println("24")
					return false
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_UP {
				if long.AllTrend.ShortTrend == models.TREND_UP {
					if isHasCrossUpper(long.Bands[bandLen-4:], true) && countBadBands(long.Bands[bandLen-4:]) > 2 {
						if mid.AllTrend.SecondTrend == models.TREND_DOWN && countBadBands(mid.Bands[bandLen-4:]) > 2 {
							if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 && countBadBands(short.Bands[bandLen-4:]) > 2 {
								log.Println("25")
								return false
							}
						}
					}
				}

				if countBadBands(long.Bands[bandLen-4:]) > 2 && isHasCrossUpper(long.Bands[bandLen-7:bandLen-4], false) {
					if short.Position == models.ABOVE_UPPER {
						if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
							log.Println("25.1")
							return false
						}

						if countCrossUpperOnBody(short.Bands[bandLen-4:]) > 1 && isHasBandDownFromUpper(short.Bands[bandLen-4:]) {
							log.Println("25.2")
							return false
						}
					}
				}
			}

			if isPreviousBandTripleHeigh(mid.Bands) {
				if (short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1) || (isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) && countCrossUpper(short.Bands[bandLen-2:]) > 1) {
					log.Println("26")
					return false
				}

				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					log.Println("27")
					return false
				}
			}

			if CountUpBand(mid.Bands[bandLen-4:]) == 4 && countBadBands(mid.Bands[bandLen-3:]) > 1 {
				if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
					log.Println("28")
					return false
				}
			}

			if shortLastBandPercent > 20 {
				if mid.Position == models.ABOVE_UPPER && countCrossUpper(mid.Bands[bandLen-4:]) == 1 {
					if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
						log.Println("29")
						return false
					}
				}
			}

			if long.Position == models.ABOVE_UPPER && countCrossUpperOnBody(long.Bands[bandLen-4:]) == 1 {
				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if isHasBandDownFromUpper(short.Bands[bandLen-2:]) {
						log.Println("31")
						return false
					}

					if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						if long.AllTrend.FirstTrend == models.TREND_UP && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
							log.Println("31.1")
							return false
						}
					}
				}

				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) > 1 && isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) {
					if short.Position == models.ABOVE_UPPER {
						log.Println("31.2")
						return false
					}
				}
			}

			if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
				if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-2]) && isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
					log.Println("32")
					return false
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend != models.TREND_UP {
				if long.Position == models.BELOW_SMA && long.PriceChanges > 20 {
					if mid.Position == models.BELOW_SMA {
						log.Println("33")
						return false
					}
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
				if countCrossUpper(long.Bands[bandLen-4:]) > 1 && countBadBands(long.Bands[bandLen-4:]) > 2 {
					if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
						log.Println("34")
						return false
					}

					if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						if countCrossUpper(short.Bands[bandLen-4:]) > 1 && countBadBands(short.Bands[bandLen-4:]) > 1 {
							log.Println("34.1")
							return false
						}
					}
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend != models.TREND_UP {
				if isHasBandDownFromUpper(long.Bands[bandLen-4:]) {
					if midLastBand.Candle.Close < float32(midLastBand.Upper) {
						if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
							log.Println("35")
							return false
						}

						if short.Position == models.ABOVE_SMA && countCrossUpper(short.Bands[bandLen-4:]) == 0 {
							log.Println("36")
							return false
						}
					}
				}

				if countCrossUpper(long.Bands[bandLen-4:]) > 1 && long.Position != models.ABOVE_UPPER {
					if mid.AllTrend.SecondTrend == models.TREND_DOWN && mid.Position != models.ABOVE_UPPER {
						if short.Position != models.ABOVE_UPPER {
							log.Println("37")
							return false
						}
					}
				}
			}

			if isHasOpenCloseAboveUpper(long.Bands[bandLen-3:]) || (isHasOpenCloseAboveUpper(long.Bands[bandLen-4:]) && getHighestHightIndex(long.Bands[bandLen-6:]) != bandLen-1) {
				if isHasBadBand(long.Bands[bandLen-3:]) {
					if countCrossUpperOnBody(mid.Bands[bandLen-4:]) <= 1 && countBadBands(mid.Bands[bandLen-4:]) > 2 {
						log.Println("38")
						return false
					}
				}
				if countAboveUpper(long.Bands[bandLen-3:]) > 1 {
					if mid.AllTrend.FirstTrend == models.TREND_UP && mid.AllTrend.SecondTrend == models.TREND_UP {
						if mid.AllTrend.ShortTrend == models.TREND_UP && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 0 {
							if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
								log.Println("38.1")
								return false
							}
						}
					}
				}
			}

			if byTurns(long.Bands[bandLen-4:]) && isHasBandDownFromUpper(long.Bands[bandLen-4:]) {
				if countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
					log.Println("39")
					return false
				}
			}

			if long.AllTrend.FirstTrend == models.TREND_UP && long.AllTrend.SecondTrend == models.TREND_UP {
				if long.AllTrend.ShortTrend == models.TREND_UP {
					if countCrossUpperOnBody(long.Bands[bandLen-4:]) == 0 {
						if countCrossUpperOnBody(mid.Bands[bandLen-4:]) <= 1 {
							if countCrossUpperOnBody(short.Bands[bandLen-4:]) <= 1 || isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
								log.Println("40")
								return false
							}
						}
					}

					if short.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
						if isHasBandDownFromUpper(short.Bands[bandLen-4:]) {
							log.Println("40.1")
							return false
						}
					}

					if isHasOpenCloseAboveUpper(long.Bands[bandLen-4:]) || isHasBandDownFromUpper(long.Bands[bandLen-4:]) {
						if countBadBands(mid.Bands[bandLen-4:]) > 2 {
							if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 && shortLastBandPercent < 5 {
								log.Println("40.2")
								return false
							}
						}
					}
				}

				if long.AllTrend.ShortTrend != models.TREND_UP && countBadBands(long.Bands[bandLen-4:]) > 2 {
					if countBadBands(mid.Bands[bandLen-4:]) > 2 && countBadBands(short.Bands[bandLen-4:]) > 2 {
						log.Println("41")
						return false
					}
				}
			}

			if long.AllTrend.SecondTrend != models.TREND_UP {
				if mid.AllTrend.Trend == models.TREND_DOWN && mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						log.Println("42")
						return false
					}
				}

				if long.Position == models.ABOVE_SMA {
					if mid.AllTrend.SecondTrend == models.TREND_UP && countBadBands(mid.Bands[bandLen-4:]) > 2 {
						if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
							log.Println("42.1")
							return false
						}
					}

					if countCrossUpperOnBody(mid.Bands[bandLen-4:]) <= 1 {
						if countCrossUpperOnBody(short.Bands[bandLen-4:]) <= 1 {
							log.Println("42.2")
							return false
						}
					}
				}
			}

			if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
				if mid.AllTrend.FirstTrend == models.TREND_UP && mid.AllTrend.SecondTrend == models.TREND_UP {
					if mid.AllTrend.ShortTrend == models.TREND_UP && isHasBandDownFromUpper(mid.Bands[bandLen-4:]) {
						log.Println("43")
						return false
					}
				}

				if long.AllTrend.FirstTrend == models.TREND_UP && long.AllTrend.SecondTrend == models.TREND_UP {
					if long.AllTrend.ShortTrend == models.TREND_UP && isHasBandDownFromUpper(long.Bands[bandLen-2:]) {
						log.Println("44")
						return false
					}
				}

				if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
					if countBadBands(mid.Bands[bandLen-4:]) > 2 && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
						log.Println("45")
						return false
					}
				}

				longPercentFromUpper := (long.Bands[bandLen-1].Upper - float64(long.Bands[bandLen-1].Candle.Close)) / float64(long.Bands[bandLen-1].Candle.Close) * 100
				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if (long.Position == models.ABOVE_UPPER && countCrossUpperOnBody(long.Bands[bandLen-4:]) == 1) || (long.Position != models.ABOVE_UPPER && countCrossUpperOnBody(long.Bands[bandLen-4:]) == 0 && longPercentFromUpper < 2) {
						if countBandPercentChangesMoreThan(short.Bands[bandLen-1:], 5) == 0 {
							if countBandPercentChangesMoreThan(mid.Bands[bandLen-1:], 5) == 0 {
								if countBandPercentChangesMoreThan(long.Bands[bandLen-1:], 5) == 0 {
									log.Println("46")
									return false
								}
							}
						}
					}
				}
			}

			if countCrossUpperOnBody(long.Bands[bandLen-4:]) <= 1 || (isHasBandDownFromUpper(long.Bands[bandLen-4:]) || isHasOpenCloseAboveUpper(long.Bands[bandLen-4:])) {
				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if mid.AllTrend.FirstTrend == models.TREND_UP && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
						if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
							log.Println("47")
							return false
						}
					}
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_DOWN && long.Position == models.ABOVE_SMA {
				if isHasCrossSMA(long.Bands[bandLen-1:], true) {
					if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
						if short.Position == models.ABOVE_UPPER && (countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 || isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1])) {
							log.Println("48")
							return false
						}
					}
				}
			}

			if long.Position == models.ABOVE_UPPER && countCrossUpperOnBody(long.Bands[bandLen-4:]) > 1 && isUpperHeadMoreThanUpperBody(long.Bands[bandLen-1]) {
				if mid.Position == models.ABOVE_UPPER {
					if short.Position == models.ABOVE_UPPER && (countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 || isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1])) {
						log.Println("49")
						return false
					}
				}
			}

			if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) && shortLastBandPercent < 5 {
				log.Println("50")
				return false
			}

			if isHasBadBand(long.Bands[bandLen-2:]) && countBadBands(long.Bands[bandLen-5:bandLen-1]) > 2 {
				if isHasOpenCloseAboveUpper(long.Bands[bandLen-4:]) || isHasBandDownFromUpper(long.Bands[bandLen-4:]) {
					if isHasBadBand(mid.Bands[bandLen-2:]) && countBadBands(mid.Bands[bandLen-5:bandLen-1]) > 2 {
						if isHasBadBand(short.Bands[bandLen-2:]) && countBadBands(short.Bands[bandLen-5:bandLen-1]) > 2 {
							log.Println("51")
							return false
						}
					}
				}

				if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP && countCrossUpper(long.Bands[bandLen-4:]) > 1 {
					if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
						log.Println("51.1")
						return false
					}
				}

				if countBadBands(mid.Bands[bandLen-3:]) >= 2 {
					if countBadBands(short.Bands[bandLen-4:]) > 2 && isHasBandDownFromUpper(short.Bands[bandLen-4:]) {
						log.Println("51.2")
						return false
					}
				}
			}

			if long.Bands[bandLen-1].Candle.Low < float32(long.Bands[bandLen-1].Lower) && long.Bands[bandLen-1].Candle.Close < float32(long.Bands[bandLen-1].Upper) {
				if isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
					log.Println("52")
					return false
				}
			}

			if long.AllTrend.Trend == models.TREND_DOWN && long.Position == models.ABOVE_SMA && isHasCrossSMA(long.Bands[bandLen-1:], true) {
				if countCrossSMA(long.Bands[bandLen-5:bandLen-1]) == 0 && countAboveSMA(long.Bands[bandLen-5:]) == 0 {
					if isUpperHeadMoreThanUpperBody(midLastBand) && isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
						log.Println("53")
						return false
					}
				}
			}

			if isUpperHeadMoreThanUpperBody(long.Bands[bandLen-1]) && isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) {
				if countDownBand(short.Bands[bandLen-2:]) > 0 {
					log.Println("54")
					return false
				}
			}

			if isHasBandDownFromUpper(long.Bands[bandLen-2:]) {
				if countDownBand(mid.Bands[bandLen-4:]) > 2 {
					if short.AllTrend.Trend == models.TREND_DOWN {
						log.Println("55")
						return false
					}
				}
			}

			if isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) && isUpperHeadMoreThanUpperBody(short.Bands[bandLen-1]) {
				if countCrossUpperOnBody(long.Bands[bandLen-4:]) == 1 && long.Position == models.ABOVE_UPPER {
					log.Println("56")
					return false
				}
			}

			if long.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && short.Position == models.ABOVE_UPPER {
				if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
					if short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
						if countBadBands(short.Bands[bandLen-4:]) > 2 {
							log.Println("57")
							return false
						}
					}
				}
			}

			if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
				if isHasCrossLower(mid.Bands[bandLen-2:], false) && isHasCrossLower(long.Bands[bandLen-2:], false) {
					log.Println("58")
					return false
				}
			}

			if isHasCrossUpper(long.Bands[bandLen-4:], false) || isHasCrossUpper(mid.Bands[bandLen-4:], false) || isHasCrossUpper(short.Bands[bandLen-4:], false) {
				if countBadBands(long.Bands[bandLen-4:]) > 2 && countBadBands(mid.Bands[bandLen-4:]) > 2 && countBadBands(short.Bands[bandLen-4:]) > 2 {
					log.Println("59")
					return false
				}
			}

			if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if isHasCrossSMA(mid.Bands[bandLen-1:], false) && isHasCrossSMA(long.Bands[bandLen-1:], false) {
						log.Println("60")
						return false
					}
				}
			}

			if isHasBandDownFromUpper(long.Bands[bandLen-2:]) || (isHasCrossUpper(long.Bands[bandLen-2:bandLen-1], false) && isHasBadBand(long.Bands[bandLen-2:bandLen-1])) {
				if countBadBands(long.Bands[bandLen-3:]) > 1 {
					if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						log.Println("61")
						return false
					}
				}
			}

			if mid.Position == models.ABOVE_UPPER || long.Position == models.ABOVE_UPPER || mid.AllTrend.ShortTrend != models.TREND_UP || long.AllTrend.ShortTrend != models.TREND_UP {
				if short.AllTrend.FirstTrend == models.TREND_DOWN || short.AllTrend.SecondTrend == models.TREND_DOWN {
					if mid.AllTrend.FirstTrend == models.TREND_DOWN || mid.AllTrend.SecondTrend == models.TREND_DOWN {
						if long.AllTrend.FirstTrend == models.TREND_DOWN || long.AllTrend.SecondTrend == models.TREND_DOWN {
							log.Println("62")
							return false
						}
					}
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
				if countBadBands(long.Bands[bandLen-4:]) > 2 {
					if countBadBands(mid.Bands[bandLen-4:]) > 2 || mid.AllTrend.SecondTrend != models.TREND_UP {
						if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
							log.Println("63")
							return false
						}
					}
				}
			}

			if countAboveUpper(long.Bands[bandLen-2:]) > 1 {
				if countCrossUpper(mid.Bands[bandLen-4:]) > 2 && countBadBands(mid.Bands[bandLen-4:]) > 1 {
					if countCrossUpperOnBody(short.Bands[bandLen-4:]) <= 1 {
						log.Println("64")
						return false
					}
				}
			}

			if long.Position == models.ABOVE_UPPER && isHasBadBand(long.Bands[bandLen-2:bandLen-1]) {
				if mid.Position == models.ABOVE_UPPER && isHasBandDownFromUpper(mid.Bands[bandLen-4:]) {
					if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						log.Println("65")
						return false
					}
				}
			}

			if isHasCrossSMA(long.Bands[bandLen-1:], true) {
				if isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) {
					if countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 || (isHasCrossUpper(short.Bands[bandLen-2:bandLen-1], true) && isHasBadBand(short.Bands[bandLen-2:bandLen-1])) {
						log.Println("66")
						return false
					}
				}
			}

			if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if countHighAboveUpper(mid.Bands[bandLen-4:]) > 1 && countBadBands(mid.Bands[bandLen-4:]) > 2 {
					log.Println("67")
					return false
				}
			}

			if short.Position == models.ABOVE_UPPER && isHasBandDownFromUpper(short.Bands[bandLen-4:]) {
				if isHasBadBand(mid.Bands[bandLen-2:bandLen-1]) && isHasCrossUpper(mid.Bands[bandLen-2:bandLen-1], true) {
					if isHasCrossSMA(long.Bands[bandLen-1:], true) {
						log.Println("68")
						return false
					}
				}
			}

			if short.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER {
				if short.AllTrend.Trend == models.TREND_UP && mid.AllTrend.Trend == models.TREND_UP && long.AllTrend.Trend == models.TREND_UP {
					if isHasBadBand(short.Bands[bandLen-2:]) && isHasBadBand(mid.Bands[bandLen-2:]) && isHasBadBand(long.Bands[bandLen-2:]) {
						log.Println("69")
						return false
					}
				}
			}

			if short.AllTrend.Trend == models.TREND_UP && short.Position == models.ABOVE_UPPER {
				if mid.Position == models.ABOVE_SMA && mid.AllTrend.Trend == models.TREND_DOWN {
					if long.Position == models.ABOVE_SMA && long.AllTrend.ShortTrend == models.TREND_DOWN {
						log.Println("70")
						return false
					}
				}
			}

			if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
				if mid.Position == models.ABOVE_SMA && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
					if short.Position == models.ABOVE_SMA && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
						if countBadBands(long.Bands[bandLen-4:]) > 2 {
							log.Println("71")
							return false
						}
					}
				}
			}

			if isHasBandDownFromUpper(long.Bands[bandLen-2:]) {
				if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
					if short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
						if short.Position == models.ABOVE_SMA || countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
							log.Println("72")
							return false
						}
					}
				}
			}

			if isHasOpenCloseAboveUpper(long.Bands[bandLen-1:]) && isHasBadBand(long.Bands[bandLen-2:]) {
				log.Println("73")
				return false
			}

			if long.Position == models.ABOVE_SMA && !isHasCrossUpper(long.Bands[bandLen-4:], false) && percentFromUpper(longLastBand) < 4 {
				if isHasBandDownFromUpper(mid.Bands[bandLen-4:]) {
					if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						log.Println("74")
						return false
					}
				}
			}

			if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_DOWN {
				if !isHasCrossUpper(mid.Bands[bandLen-4:], true) && !isHasCrossUpper(short.Bands[bandLen-4:], true) {
					log.Println("75")
					return false
				}
			}

			if mid.AllTrend.FirstTrend == models.TREND_DOWN && short.AllTrend.FirstTrend == models.TREND_DOWN {
				if long.AllTrend.FirstTrend == models.TREND_DOWN {
					if countCrossUpperOnBody(short.Bands[bandLen-5:]) == 1 && short.Position == models.ABOVE_UPPER {
						if countCrossUpper(mid.Bands[bandLen-5:]) == 0 && countCrossUpper(long.Bands[bandLen-5:]) == 0 {
							log.Println("76")
							return false
						}
					}
				}

				if isUpperHeadMoreThanUpperBody(shortLastBand) && countCrossUpper(long.Bands[bandLen-5:]) == 0 {
					if countCrossUpperOnBody(short.Bands[bandLen-5:]) == 1 && short.Position == models.ABOVE_UPPER {
						if countCrossUpperOnBody(mid.Bands[bandLen-5:]) == 1 {
							log.Println("77")
							return false
						}
					}
				}
			}

			if long.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && short.Position == models.ABOVE_UPPER {
				if countCrossUpper(long.Bands[bandLen-4:]) == 1 && countBelowSMA(long.Bands[bandLen-4:], false) > 0 {
					if countCrossUpper(mid.Bands[bandLen-4:]) == 1 {
						log.Println("78")
						return false
					}
				}
			}

			midPercentFromUpper := (midLastBand.Upper - float64(midLastBand.Candle.Close)) / float64(midLastBand.Candle.Close) * 100
			if isHasCrossSMA(long.Bands[bandLen-1:], true) && countAboveSMAStrict(long.Bands[bandLen-4:]) == 0 {
				if midLastBand.Candle.Close < float32(midLastBand.Upper) && midPercentFromUpper < 3 {
					if short.Position == models.ABOVE_UPPER && countCrossUpperOnBody(short.Bands[bandLen-4:]) == 1 {
						log.Println("79")
						return false
					}
				}
			}

			if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 3) >= 1 {
				if !isHasOpenCloseAboveUpper(short.Bands[len(short.Bands)-4:]) {
					return true
				}
			}

		}
	}

	return false
}

func percentFromUpper(band models.Band) float32 {
	return (float32(band.Upper) - band.Candle.Close) / band.Candle.Close * 100
}

func byTurns(bands []models.Band) bool {
	swithcer := false
	if bands[0].Candle.Open < bands[0].Candle.Close {
		swithcer = true
	}

	for _, band := range bands {
		if swithcer {
			if band.Candle.Open > band.Candle.Close {
				return false
			}
			swithcer = false
		} else {
			if band.Candle.Open < band.Candle.Close {
				return false
			}
			swithcer = true
		}
	}

	return true
}

func isHourInChangesLong(currentHour int, currentMinute int) bool {
	hours := []int{23, 0, 3, 4, 7, 8, 11, 12, 15, 16, 19, 20}
	for i, hour := range hours {
		if currentHour == hour {
			if i%2 == 0 && currentMinute >= 45 {
				return true
			} else if i%2 == 1 && currentMinute < 45 {
				return true
			}
		}
	}

	return false
}

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time) bool {
	if allIntervalCrossUpperOnBodyMoreThanThresholdAndJustOne(short, mid, long, currentTime) {
		skipped = true
		log.Println("skipped1: ", short.Symbol)
		return true
	}

	skipped = false

	return false
}

func isPreviousBandTripleHeigh(bands []models.Band) bool {
	index := getIndexPreviousBand(bands)
	lastBand := bands[len(bands)-index]
	lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open
	if lastBandHeight < 0 {
		return false
	}

	var higestHeigh float32 = 0
	for _, band := range bands[len(bands)-(4+index) : len(bands)-index] {
		high := band.Candle.Close - band.Candle.Open
		if high > higestHeigh {
			higestHeigh = high
		}
	}

	return higestHeigh*2.5 < lastBandHeight
}

func getIndexPreviousBand(bands []models.Band) int {
	lastBand := bands[len(bands)-2]
	lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open
	if lastBandHeight < 0 {
		lastBand := bands[len(bands)-2]
		lastBandHeight := lastBand.Candle.Close - lastBand.Candle.Open
		if lastBandHeight > 0 {
			return 3
		}
	}
	return 2
}

func isUpperHeadMoreThanUpperBody(band models.Band) bool {
	allBody := band.Candle.Close - band.Candle.Open
	head := band.Candle.Close - float32(band.Upper)
	percent := head / allBody * 100
	return percent > 51
}

func countBandPercentChangesMoreThan(bands []models.Band, percent float32) int {
	count := 0
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			percentChanges := (band.Candle.Close - band.Candle.Open) / band.Candle.Open * 100
			if percentChanges >= percent {
				count++
			}
		}
	}

	return count
}

func isTailMoreThan(band models.Band, percent float32) bool {
	bodyTail := band.Candle.Close - band.Candle.Low
	tail := band.Candle.Open - band.Candle.Low
	tailPercent := tail / bodyTail * 100

	return tailPercent > percent
}

func isHasUpperHeadMoreThanUpperBody(bands []models.Band) bool {
	for _, band := range bands {
		if isUpperHeadMoreThanUpperBody(band) {
			return true
		}
	}

	return false
}

func getHigestPercentChangesIndex(bands []models.Band) int {
	highestIndex := 0
	for i, band := range bands {
		percentHigest := (bands[highestIndex].Candle.Close - bands[highestIndex].Candle.Open) / bands[highestIndex].Candle.Open * 100
		percent := (band.Candle.Close - band.Candle.Open) / band.Candle.Open * 100
		if percentHigest < percent {
			highestIndex = i
		}
	}

	return highestIndex
}

func GetSkipped() bool {
	if skipped {
		skipped = false

		return true
	}

	return false
}
