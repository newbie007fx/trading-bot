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
			count++
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

	bandLen := len(longInterval.Bands)
	longLastBand := longInterval.Bands[bandLen-1]
	midLastBand := midInterval.Bands[bandLen-1]

	midPercentFromUpper := (midLastBand.Upper - float64(midLastBand.Candle.Close)) / float64(midLastBand.Candle.Close) * 100

	longPercentFromSMA := (float32(longLastBand.SMA) - longLastBand.Candle.Close) / longLastBand.Candle.Close * 100

	shortLastBand := shortInterval.Bands[bandLen-1]
	shortPercentFromUpper := (float32(shortLastBand.Upper) - shortLastBand.Candle.Close) / shortLastBand.Candle.Close * 100

	longPercentFromUpper := (float32(longLastBand.Upper) - longLastBand.Candle.Close) / longLastBand.Candle.Close * 100

	midCloseBandAverage := closeBandAverage(midInterval.Bands[bandLen/2:]) * 2
	if midCloseBandAverage > 5 {
		midCloseBandAverage = 5
	} else if midCloseBandAverage < 4 {
		midCloseBandAverage = 4
	}

	shortCloseBandAverage := closeBandAverage(shortInterval.Bands[bandLen/2:]) * 2
	if shortCloseBandAverage > 5 {
		shortCloseBandAverage = 5
	} else if shortCloseBandAverage < 4 {
		shortCloseBandAverage = 4
	}

	if isHasOpenCloseAboveUpper(longInterval.Bands[len(longInterval.Bands)-2:]) && countBandPercentChangesMoreThan(longInterval.Bands[bandLen-4:], 6) != 1 {
		ignoredReason = "contain open close above upper"
		return true
	}

	if longInterval.Position == models.ABOVE_UPPER {
		if isUpperHeadMoreThanUpperBody(longInterval.Bands[bandLen-1]) || isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-1:]) {
			if isBandHeadDoubleBody(midInterval.Bands[bandLen-1:]) || isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
				if isHasCrossUpper(shortInterval.Bands[bandLen-1:], true) && isBandHeadDoubleBody(shortInterval.Bands[bandLen-2:]) {
					if getBodyPercentage(shortInterval.Bands[bandLen-1]) < 95 || !isLastBandPercentMoreThan10AndJustOnce(shortInterval.Bands) {
						ignoredReason = "cross upper, mid head more than body upper, short head double body"
						return true
					}
				}
			}

			if CalculateTrendShort(shortInterval.Bands[bandLen-5:bandLen-1]) == models.TREND_DOWN && isHasCrossUpper(shortInterval.Bands[bandLen-5:bandLen-1], false) {
				ignoredReason = "upper head more than upper body"
				return true
			}

			if countAboveUpper(midInterval.Bands[bandLen-4:]) > 0 && isHasBadBand(midInterval.Bands[bandLen-4:]) {
				ignoredReason = "upper head more than upper body and mid has open close above upper"
				return true
			}
		}

		if countHeadMoreThanBody(longInterval.Bands[bandLen-4:]) > 1 {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen/2:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen/2:]) == 1 {
					ignoredReason = "contain head more than body more than one, mid and short cross upper and just one"
					return true
				}
			}
		}

		if countCrossUpper(longInterval.Bands[bandLen-2:]) == 1 {
			if midInterval.AllTrend.FirstTrend == models.TREND_UP && midInterval.AllTrend.SecondTrend == models.TREND_UP {
				if countCrossUpper(longInterval.Bands[bandLen-4:]) == 1 || (countDownBand(longInterval.Bands[bandLen-4:]) > 0 && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1) {
					if shortInterval.Position == models.ABOVE_UPPER {
						if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
							ignoredReason = "cross upper and just one, short inter val head more than body upper"
							return true
						}
					}
				}
			}
		}

		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && countCrossUpper(longInterval.Bands[bandLen-4:]) > 1 {
			ignoredReason = "mid first trend down"
			return true
		}

		if midInterval.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
			if countDownBand(shortInterval.Bands[bandLen-4:]) > 2 {
				ignoredReason = "cross upper, short have many band down"
				return true
			}

			if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-2:]) || isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				if getBodyPercentage(shortInterval.Bands[bandLen-1]) < 95 || !isLastBandPercentMoreThan10AndJustOnce(shortInterval.Bands) {
					if shortInterval.PriceChanges < 10 || isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-4:]) || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-4:]) || isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-4:]) {
						ignoredReason = "cross upper, mid head more than body upper"
						return true
					}
				}
			}

			shortPercent := (shortInterval.Bands[bandLen-1].Candle.Close - shortInterval.Bands[bandLen-1].Candle.Open) / shortInterval.Bands[bandLen-1].Candle.Open * 100
			if countBadBands(shortInterval.Bands[bandLen-4:]) > 1 && shortPercent > 6 && (countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 2) == 1 || isHasBadBand(shortInterval.Bands[bandLen-2:])) {
				ignoredReason = "long and mid cross upper, short significan up and just one"
				return true
			}

			if shortInterval.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				ignoredReason = "all cross upper, mid and short has upper head more than upper body"
				return true
			}
		}

		if countCrossUpperOnBody(longInterval.Bands[bandLen-4:]) == 1 {
			if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-1:]) {
				ignoredReason = "mid contain open close above upper"
				return true
			}

			if isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) || countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER {
					if !longSignificanUpAndJustOne(longInterval.Bands) {
						ignoredReason = "mid and short has head more than body upper"
						return true
					}

					if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-4:]) || countBadBands(shortInterval.Bands[bandLen-3:]) > 1 {
						ignoredReason = "long and mid cros upper and just one, short contain open close above upper"
						return true
					}

					if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
						ignoredReason = "long and mid cros upper and just one, upper head more than upper bodyr"
						return true
					}

					if (isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) && midInterval.Position == models.ABOVE_UPPER) || (isHasBandDownFromUpper(midInterval.Bands[bandLen-2:]) && midInterval.Position != models.ABOVE_UPPER) {
						if countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) == 1 {
							ignoredReason = "long and short cros upper and just one, mid head more than upper bodyr"
							return true
						}
					}
				}

			}

			if isContainHeadMoreThanBody(midInterval.Bands[bandLen-1:]) && isContainHeadMoreThanBody(shortInterval.Bands[bandLen-2:]) {
				ignoredReason = "mid and short has head more than body"
				return true
			}

			if countCrossUpperOnBody(longInterval.Bands[bandLen-8:]) > 1 && (CalculateTrendShort(longInterval.Bands[bandLen-5:bandLen-1]) != models.TREND_UP || CalculateTrendShort(longInterval.Bands[bandLen-4:bandLen-1]) != models.TREND_UP) {
				ignoredReason = "after up"
				return true
			}

			if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 && !longSignificanUpAndJustOne(longInterval.Bands) {
				ignoredReason = "cross upper and just one"
				return true
			}

			if CountUpBand(midInterval.Bands[bandLen-2:]) == 2 && isHasHammer(midInterval.Bands[bandLen-2:]) {
				ignoredReason = "cross upper and just one, mid has hammer"
				return true
			}

			if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 || midInterval.AllTrend.SecondTrend != models.TREND_UP || (midInterval.Position == models.ABOVE_SMA && isHasCrossUpper(midInterval.Bands[bandLen-4:], true)) {
				if midInterval.PriceChanges > 10 && shortInterval.Position == models.ABOVE_UPPER && shortInterval.PriceChanges > 5 && !longSignificanUpAndJustOne(longInterval.Bands) {
					if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen/2:], 3) > 1 {
						ignoredReason = "cross upper and just one, mid trend down"
						return true
					}
				}
			}

			if isHasCrossUpper(midInterval.Bands[bandLen-1:], true) && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen/2:]) == 1 {
					ignoredReason = "cross upper and just one, mid and short cross upper just one too"
					return true
				}
			}

			if longInterval.PriceChanges > 10 {
				if midInterval.Position == models.ABOVE_UPPER && isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-4:]) {
					if shortInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(shortInterval.Bands[bandLen/2:], true) {
						ignoredReason = "cross upper and just one, mid contain open close above opper and short not cross upper"
						return true
					}
				}

				if countCrossUpper(midInterval.Bands[bandLen-4:]) > 1 {
					if countCrossUpper(shortInterval.Bands[bandLen-4:]) > 1 {
						if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) != 1 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) != 1 && !longSignificanUpAndJustOne(longInterval.Bands) {
							ignoredReason = "all band cross upper more than one but already reach higest"
							return true
						}

						if countBandPercentChangesMoreThan(shortInterval.Bands[len(shortInterval.Bands)-4:], 3) == 1 && getHigestPercentChangesIndex(shortInterval.Bands[len(shortInterval.Bands)-4:]) == 0 && countDownBand(shortInterval.Bands[len(shortInterval.Bands)-4:]) > 1 {
							ignoredReason = "all band cross upper more than one short starting down"
							return true
						}
					}
				}
			}

			if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-1:]) {
				if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-2:]) || isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
					ignoredReason = "short and mid contain open close above upper"
					return true
				}
			}

			if midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_UPPER && midPercentFromUpper < 5 {
				ignoredReason = "cross upper and just one and mid above sma"
				return true
			}

			if isHasBandDownFromUpper(midInterval.Bands[bandLen-2:]) && !isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) {
				ignoredReason = "cross upper just one. and mid has down from upper"
				return true
			}

			if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) > 2 {
				if countBadBands(shortInterval.Bands[bandLen-3:]) >= 2 {
					ignoredReason = "long cross upper just one, mid cross upper on body more than 2"
					return true
				}
			}

			if midInterval.Position == models.ABOVE_UPPER && (isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-4:]) || isHasBandDownFromUpper(midInterval.Bands[bandLen-4:])) {
				if shortInterval.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
					ignoredReason = "long cross upper just one, mid contain open close above upper and short upper head more than upper body"
					return true
				}
			}

			if midInterval.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
				if shortPercentFromUpper < 3 {
					ignoredReason = "long cross upper just one, mid upper head more than upper body, short percent from upper below 3"
					return true
				}
			}
		}

		if countAboveUpper(longInterval.Bands[bandLen-4:]) > 0 {
			if midInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
				if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-2:]) || isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-2:]) || countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) == 1 {
					ignoredReason = "above upper and mid above sma not cross upper, short above upper"
					return true
				}
			}

			if countDownBand(longInterval.Bands[bandLen-3:]) >= 1 {
				if midInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
					if isBandHeadDoubleBody(shortInterval.Bands[bandLen-2 : bandLen-1]) {
						ignoredReason = "above upper, mid not cross upper, short head double body"
						return true
					}

					if isTailMoreThan(midInterval.Bands[bandLen-1], 40) {
						ignoredReason = "above upper, mid not cross upper tail more than threshold"
						return true
					}
				}
			}
		}

		if midInterval.AllTrend.SecondTrend != models.TREND_UP && midInterval.Position == models.ABOVE_SMA {
			if !isHasCrossUpper(midInterval.Bands[bandLen-5:], true) {
				if shortInterval.Position == models.ABOVE_UPPER && (countCrossUpper(shortInterval.Bands[bandLen-4:]) > 1 || isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-2:])) {
					ignoredReason = "cross upper, mid second trend sideway, short cross upper more than one "
					return true
				}
			}
		}

		if longInterval.Bands[bandLen-2].Candle.Hight > float32(longInterval.Bands[bandLen-2].Upper) && (isBandHeadDoubleBody(longInterval.Bands[bandLen-2:bandLen-1]) || longInterval.Bands[bandLen-2].Candle.Open > longInterval.Bands[bandLen-2].Candle.Close) && !longSignificanUpAndJustOne(longInterval.Bands) {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) {
					ignoredReason = "previous band head double body, mid cross upper, and short upper head more than upper body1"
					return true
				}
			}
		}

		if countCrossUpperOnBody(longInterval.Bands[bandLen-4:]) > 2 {
			if longInterval.PriceChanges > 15 {
				if midInterval.Position == models.ABOVE_UPPER && countDownBand(midInterval.Bands[bandLen-4:]) > 0 {
					if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-5:]) == 1 {
						ignoredReason = "cross upper, price change alredy more than 15"
						return true
					}
				}
			}

			if midInterval.Position == models.ABOVE_SMA && countCrossUpper(midInterval.Bands[bandLen-4:]) == 0 {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) > 2 {
					if midPercentFromUpper < 5 {
						ignoredReason = "long and short above upper, mid above sma"
						return true
					}
				}
			}
		}

		if countCrossUpper(longInterval.Bands[bandLen-3:]) > 0 {
			if countCrossUpper(midInterval.Bands[bandLen-4:]) > 0 {
				if countCrossUpper(shortInterval.Bands[bandLen-4:]) > 0 {
					if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) != 1 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) != 1 && countBandPercentChangesMoreThan(longInterval.Bands[bandLen-4:], 5) != 1 && !longSignificanUpAndJustOne(longInterval.Bands) {
						ignoredReason = "all band cross upper more than one but percent change is minor"
						return true
					}

					if isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) && countBandPercentChangesMoreThan(shortInterval.Bands[len(shortInterval.Bands)-1:], 6) == 1 && isHasUpperHeadMoreThanUpperBody(midInterval.Bands[len(midInterval.Bands)-1:]) && countBandPercentChangesMoreThan(midInterval.Bands[len(midInterval.Bands)-1:], 10) == 1 {
						ignoredReason = "all band cross upper more than one but percent double threshold"
						return true
					}

					if !(countBandPercentChangesMoreThan(shortInterval.Bands[len(shortInterval.Bands)-4:], 3) >= 1 && countBandPercentChangesMoreThan(shortInterval.Bands[len(shortInterval.Bands)-4:], 2) > 1) && countBandPercentChangesMoreThan(midInterval.Bands[len(midInterval.Bands)-4:], 5) != 1 {
						ignoredReason = "all band cross upper and threshold literally just one"
						return true
					}
				}
			}

			if midInterval.PriceChanges > 15 && countAboveUpper(midInterval.Bands[bandLen-4:]) > 0 && countAboveUpper(shortInterval.Bands[bandLen-4:]) > 0 {
				ignoredReason = "mid and short Open close above upper"
				return true
			}
		}

		if countCrossUpper(longInterval.Bands[bandLen-4:]) > 1 {
			if midInterval.Position == models.ABOVE_UPPER && shortInterval.Position == models.ABOVE_UPPER {
				if isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-4:]) && isHasBadBand(midInterval.Bands[bandLen-4:]) {
					ignoredReason = "above upper and mid has open close above upper and badbands"
					return true
				}

				if isUpperHeadMoreThanUpperBody(longInterval.Bands[bandLen-1]) && isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
					ignoredReason = "above upper and upper head more than body, mid too"
					return true
				}
			}
		}

		if midInterval.Position == models.ABOVE_UPPER && shortInterval.Position == models.ABOVE_UPPER {
			if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-2:]) && getHigestPercentChangesIndex(shortInterval.Bands[len(shortInterval.Bands)-4:]) == 3 {
				if countBandPercentChangesMoreThan(shortInterval.Bands[len(shortInterval.Bands)-4:], 5) == 1 && countBandPercentChangesMoreThan(shortInterval.Bands[len(shortInterval.Bands)-4:], 2) == 1 {
					ignoredReason = "mid contain above upper, short significan up and just one"
					return true
				}
			}
		}

		if countBadBands(longInterval.Bands[bandLen-2:]) > 0 {
			if isHasBandDownFromUpper(midInterval.Bands[bandLen-4:]) {
				if !isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) {
					ignoredReason = "above upper, mid has down from upper, short not cross upper"
					return true
				}
			}
		}

		if isHasBandDownFromUpper(longInterval.Bands[bandLen-2:]) {
			if !isHasCrossUpper(midInterval.Bands[bandLen-4:], false) && !isHasCrossUpper(shortInterval.Bands[bandLen-4:], false) {
				if midPercentFromUpper < 5 && shortPercentFromUpper < 5 {
					ignoredReason = "above upper, mid and short not cross upper"
					return true
				}
			}

			if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 && countBadBands(midInterval.Bands[bandLen-4:]) > 2 {
				ignoredReason = "contain down from upper, mid cross upper just one"
				return true
			}
		}

		if isHasBandDownFromUpper(longInterval.Bands[bandLen-3:]) && countBadBands(longInterval.Bands[bandLen-3:]) > 1 {
			if !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) && midPercentFromUpper < 5 {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) == 1 {
					ignoredReason = "above upper but contain down from upper, mid not cross upper and short cross upper just one"
					return true
				}
			}
		}

		if countCrossUpper(midInterval.Bands[bandLen-4:]) <= 1 && countCrossUpper(longInterval.Bands[bandLen-4:]) == 1 {
			if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) > 1 && isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				ignoredReason = "upper head more thant upper body"
				return true
			}
		}

		if shortInterval.Position == models.ABOVE_UPPER && midInterval.Position == models.ABOVE_UPPER {
			if isLastBandPercentMoreThan10AndJustOnce(midInterval.Bands) || (isHasBadBand(midInterval.Bands[bandLen-2:]) && isLastBandPercentMoreThan10AndJustOnce(midInterval.Bands)) {
				ignoredReason = "band percent more than 10 and just one"
				return true
			}
		}

		if isHasBandDownFromUpper(longInterval.Bands[bandLen-4:]) || isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-4:]) || isHasUpperHeadMoreThanUpperBody(longInterval.Bands[bandLen-4:]) {
			if countBadBands(midInterval.Bands[bandLen-3:]) > 1 {
				if isHasBandDownFromUpper(shortInterval.Bands[bandLen-6:]) || isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-6:]) {
					ignoredReason = "updown, contain band down from upper"
					return true
				}
			}
		}
	}

	if longInterval.Position == models.ABOVE_SMA {
		if !isHasCrossUpper(longInterval.Bands[bandLen-2:], true) {
			if !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-8:]) == 1 {
					ignoredReason = "short above upper but just one"
					return true
				}

				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpper(shortInterval.Bands[bandLen-8:]) > 5 && countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) != 1 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) != 1 {
					ignoredReason = "short above upper many, no percent above threshold"
					return true
				}
			}
		}

		if longInterval.AllTrend.SecondTrend != models.TREND_UP {
			if (countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], 5) != 1 || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-5:], 3) != 1) && !(isHasCrossSMA(longInterval.Bands[bandLen-2:bandLen-1], false) && countCrossSMA(longInterval.Bands[bandLen-4:]) == 1) {
				ignoredReason = "above sma and second trend not up"
				return true
			}

			if midInterval.AllTrend.SecondTrend != models.TREND_UP && isHasCrossUpper(midInterval.Bands[bandLen-1:], true) && countCrossUpper(midInterval.Bands[bandLen-7:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && isHasCrossUpper(shortInterval.Bands[bandLen-1:], true) && countCrossUpper(shortInterval.Bands[bandLen-7:]) == 1 {
					ignoredReason = "above sma and second trend not up, short and mid cross upper"
					return true
				}
			}
		}

		if !isHasCrossUpper(longInterval.Bands[bandLen-5:], true) {
			if longPercentFromUpper < 5 {
				if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) > 0 || countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) > 0 {
					ignoredReason = "above sma and percent from upper bellow 5"
					return true
				}
			}

			if longInterval.PriceChanges > 10 {
				if isHasCrossSMA(longInterval.Bands[bandLen-1:], false) || longPercentFromUpper < 3 {
					if countAboveSMA(longInterval.Bands[bandLen-5:]) < 2 {
						if isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-4:]) || countDownBand(midInterval.Bands[bandLen-4:]) > 1 {
							ignoredReason = "above sma, significan up and mid has open close above upper"
							return true
						}
					}

					if isHasCrossUpper(midInterval.Bands[bandLen-2:], true) {
						if isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-2:]) {
							if countBandPercentChangesMoreThan(shortInterval.Bands, 5) == 0 && countBandPercentChangesMoreThan(shortInterval.Bands, 3) < 2 {
								ignoredReason = "above sma, not significan up"
								return true
							}

							if longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
								ignoredReason = "above sma, not significan up, second trend down"
								return true
							}
						}
					}
				}
			}
		}

		if isHasUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1:]) {
			if isBandHeadDoubleBody(shortInterval.Bands[bandLen-1:]) {
				ignoredReason = "above sma, mid upper head more than body, short head double body"
				return true
			}
		}

		if isHasCrossSMA(longInterval.Bands[bandLen-2:], false) || (longPercentFromUpper < 5 && isHasCrossUpper(longInterval.Bands[bandLen-3:], true)) {
			if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-1:]) {
				ignoredReason = "above sma and mid contain open close above upper"
				return true
			}

			if (midInterval.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(midLastBand)) || (shortInterval.Position == models.ABOVE_UPPER && countCrossUpper(shortInterval.Bands[bandLen-4:]) == 1) {
				ignoredReason = "above sma and cross sma, mid or short cross upper just one or upper head more than body"
				return true
			}
		}

		if midInterval.Position == models.ABOVE_SMA {
			if shortInterval.Position == models.ABOVE_SMA {
				if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], midCloseBandAverage) != 1 || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-5:], shortCloseBandAverage) != 1 {
					ignoredReason = "all interval above sma"
					return true
				}
			}

			if midPercentFromUpper < 5 && longPercentFromUpper < 5 {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) > 1 {
					ignoredReason = "short cross upper but mid and long above sma and percent from upper below 5"
					return true
				}
			}
		}

		if isHasCrossSMA(longInterval.Bands[bandLen-1:], true) {
			if longPercentFromUpper < 9 {
				if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
					if isHasCrossUpper(shortInterval.Bands[bandLen-1:], true) && countCrossUpper(shortInterval.Bands[bandLen-4:]) == 1 {
						ignoredReason = "cross sma, mid and short cross upper"
						return true

					}
				}

				if midInterval.Position == models.ABOVE_SMA && isHasCrossUpper(midInterval.Bands[bandLen-1:], true) && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
					if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-3:]) {
						ignoredReason = "above sma and mid cross upper and short open close above upper"
						return true
					}
				}
			}

			if longInterval.AllTrend.SecondTrend != models.TREND_UP || (longInterval.AllTrend.FirstTrend == models.TREND_DOWN && longInterval.AllTrend.SecondTrend == models.TREND_UP && longInterval.AllTrend.FirstTrendPercent < longInterval.AllTrend.SecondTrendPercent) {
				if countBandPercentChangesMoreThan(longInterval.Bands[bandLen-4:], 5) == 0 || countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) == 0 {
					if shortInterval.Position == models.ABOVE_UPPER {
						ignoredReason = "trend not up, or up but just minor"
						return true
					}
				}
			}
		}

		if countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 || midInterval.AllTrend.SecondTrend != models.TREND_UP || (midInterval.Position == models.ABOVE_SMA && isHasCrossUpper(midInterval.Bands[bandLen-4:], true)) {
			if midInterval.PriceChanges > 10 && shortInterval.Position == models.ABOVE_UPPER && shortInterval.PriceChanges > 5 {
				if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen/2:], 5) > 1 {
					ignoredReason = "above sma, mid trend down"
					return true
				}
			}
		}

		if countCrossUpper(midInterval.Bands[bandLen-4:]) > 1 && longPercentFromUpper < 5 {
			if countCrossUpper(shortInterval.Bands[bandLen-4:]) > 1 {
				if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-1:]) || isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-1:]) || (isHasUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1:]) && checkTime.Minute() > 18 && checkTime.Minute() < 3) {
					if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 5) == 0 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) == 0 {
						ignoredReason = "cross sma, mid and short cross upper more than one but have no significan up band"
						return true
					}
				}
			}
		}

		if isHasCrossUpper(longInterval.Bands[bandLen-1:], true) && countCrossUpper(longInterval.Bands[bandLen-4:]) == 1 {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 5) == 0 {
					ignoredReason = "long & mid cross upper but just one"
					return true
				}
			}

			if isHasBandDownFromUpper(midInterval.Bands[bandLen-2:]) || countAboveUpper(midInterval.Bands[bandLen-2:]) > 0 {
				if countCrossUpper(shortInterval.Bands[bandLen-4:]) == 1 && isHasCrossUpper(shortInterval.Bands[bandLen-1:], true) {
					ignoredReason = "long & short cross upper but just one, mid contain open close above upper"
					return true
				}
			}
		}

		if countCrossUpper(longInterval.Bands[bandLen-3:]) == 0 {
			if isHasCrossUpper(midInterval.Bands[bandLen-4:], true) && midInterval.PriceChanges < 10 {
				if shortInterval.PriceChanges > 5 && countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) == 0 && !(isHasCrossSMA(longInterval.Bands[bandLen-2:bandLen-1], false) && countCrossSMA(longInterval.Bands[bandLen-4:]) == 1) {
					ignoredReason = "mid cross upper, and short interval do not have significan band"
					return true
				}

				if countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 && countCrossUpper(shortInterval.Bands[bandLen-4:]) > 1 {
					if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) == 0 || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) == 0 {
						ignoredReason = "mid cross upper, and short interval do not have significan band 2"
						return true
					}
				}
			}

			if longPercentFromUpper < 5 {
				if midInterval.Position == models.ABOVE_UPPER {
					if shortInterval.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
						ignoredReason = "below upper, percent from upper below 5 and short upper head more thant body"
						return true
					}
				}
			}
		}

		if longInterval.Bands[bandLen-2].Candle.Hight > float32(longInterval.Bands[bandLen-2].Upper) && (isBandHeadDoubleBody(longInterval.Bands[bandLen-2:bandLen-1]) || longInterval.Bands[bandLen-2].Candle.Open > longInterval.Bands[bandLen-2].Candle.Close) {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) {
					ignoredReason = "previous band head double body, mid cross upper, and short upper head more than upper body"
					return true
				}
			}
		}

		if isHasCrossUpper(longInterval.Bands[bandLen-4:bandLen-2], true) && !isHasCrossUpper(longInterval.Bands[bandLen-2:], true) {
			if longInterval.PriceChanges > 10 && midInterval.PriceChanges > 10 && longPercentFromUpper < 7 {
				if shortInterval.Position == models.ABOVE_UPPER && shortInterval.Bands[bandLen-2].Candle.Open > shortInterval.Bands[bandLen-2].Candle.Close {
					ignoredReason = "after up down, and then start up again but seems just minor up"
					return true
				}
			}
		}

		if !isHasCrossUpper(longInterval.Bands[bandLen-4:], true) {
			if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-1:], 3) == 0 && countDownBand(shortInterval.Bands[bandLen-3:]) > 0 && isContainHeadMoreThanBody(shortInterval.Bands[bandLen-3:]) {
				ignoredReason = "above sma and short percent below 3"
				return true
			}

			if countBadBands(longInterval.Bands[bandLen-4:]) > 2 {
				if midInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
					if (midInterval.AllTrend.FirstTrend != models.TREND_UP || midInterval.AllTrend.SecondTrend != models.TREND_UP) && midPercentFromUpper < 5 {
						ignoredReason = "above sma, min above sma not cross upper"
						return true
					}
				}
			}
		}

		if countCrossUpper(longInterval.Bands[bandLen-4:]) > 0 && countCrossUpperOnBody(longInterval.Bands[bandLen-4:]) == 0 {
			if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-1:]) || isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-1:]) || (isHasUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1:]) && checkTime.Minute() > 18 && checkTime.Minute() < 3) {
				if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-1:], 3) == 0 && countCrossUpper(shortInterval.Bands[bandLen-3:]) > 1 && shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "above sma and short percent below 3 2nd"
					return true
				}
			}
		}

		if isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) && longPercentFromUpper < 5 {
			if countBadBands(midInterval.Bands[bandLen-4:]) > 1 {
				ignoredReason = "above sma, mid contain bad band"
				return true
			}
		}

		if countCrossUpperOnBody(longInterval.Bands[bandLen:]) == 0 && longPercentFromUpper < 5 {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) < 2 {
				if shortInterval.Position == models.ABOVE_SMA && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) == 0 && shortPercentFromUpper < 3 {
					ignoredReason = "above sma, mid cross upper just one, short above sma"
					return true
				}
			}
		}

		if longInterval.AllTrend.ShortTrend == models.TREND_DOWN && countDownBand(longInterval.Bands[bandLen-4:]) > 2 {
			ignoredReason = "above sma, short trend down"
			return true
		}
	}

	if longInterval.Position == models.BELOW_SMA {
		if (!isHasCrossSMA(longInterval.Bands[bandLen-4:], false) && longPercentFromSMA < 5) || (countCrossSMA(longInterval.Bands[bandLen-7:]) == 0 || (countCrossSMA(longInterval.Bands[bandLen-7:]) == 1 && isHasCrossSMA(longInterval.Bands[bandLen-1:], false))) && longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if (isHasCrossUpper(midInterval.Bands[bandLen-3:], true) && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1) || (midInterval.Position == models.ABOVE_SMA && midPercentFromUpper < 5) {
				if isHasCrossUpper(shortInterval.Bands[bandLen-1:], true) && (countCrossUpper(shortInterval.Bands[bandLen-4:]) == 1 || isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:])) {
					ignoredReason = "below sma, mid and short cross upper"
					return true
				}
			}
		}

		if longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if midInterval.AllTrend.SecondTrend != models.TREND_UP && isHasCrossUpper(midInterval.Bands[bandLen-1:], true) && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
				if shortInterval.Position == models.ABOVE_UPPER && isHasCrossUpper(shortInterval.Bands[bandLen-1:], false) && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) == 1 {
					ignoredReason = "below sma, second trend down, mid and short cross upper"
					return true
				}
			}

			if isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) || isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				ignoredReason = "below sma, second trend down, mid and short has upper head more than body"
				return true
			}
		}

		if (longInterval.AllTrend.FirstTrend != models.TREND_UP && longInterval.AllTrend.SecondTrend != models.TREND_UP) && (longInterval.AllTrend.FirstTrend == models.TREND_DOWN || longInterval.AllTrend.SecondTrend == models.TREND_DOWN) {
			if countCrossSMA(longInterval.Bands[bandLen-4:]) == 1 && countAboveSMA(longInterval.Bands[bandLen-4:]) == 0 {
				if midInterval.Bands[bandLen-2].Candle.Open > midInterval.Bands[bandLen-2].Candle.Close && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-1:], 5) == 0 {
					ignoredReason = "below sma and mid previous band down and minor up"
					return true
				}
			}
		}

		if isTailMoreThan(longInterval.Bands[bandLen-2], 40) {
			if countCrossUpper(midInterval.Bands[bandLen-4:]) > 1 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 3) == 0 {
				ignoredReason = "below sma, tail more than body and mid cross upper more than one but minor up"
				return true
			}
		}

		if midInterval.Position == models.BELOW_SMA || isHasBadBand(midInterval.Bands[bandLen-2:]) {
			if isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) || isHasBadBand(shortInterval.Bands[bandLen-4:]) || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) == 0 {
				ignoredReason = "long bellow sma, short bad band"
				return true
			}
		}
	}

	if isHasCrossUpper(longInterval.Bands[bandLen-4:], true) && countCrossUpper(midInterval.Bands[bandLen-4:]) <= 1 {
		if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-5:]) == 1 {
			if isHasBadBand(shortInterval.Bands[bandLen-2:]) || isHasBadBand(midInterval.Bands[bandLen-2:]) {
				if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 2) == 1 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-2:], 5) == 0 {
					ignoredReason = "minor up"
					return true
				}
			}
		}
	}

	higestHightIndex := getHighestHightIndex(longInterval.Bands[len(longInterval.Bands)-5:])
	if higestHightIndex < len(longInterval.Bands[len(longInterval.Bands)-5:])-2 {
		if midInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
			ignoredReason = "previous band have higest high and mid cross upper just one"
			return true
		}
	}

	midLastBandPriceChange := (midLastBand.Candle.Close - midLastBand.Candle.Open) / midLastBand.Candle.Open * 100
	if isUpperHeadMoreThanUpperBody(longInterval.Bands[bandLen-1]) {
		if isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) && !(countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], 5) == 1 && midLastBandPriceChange > 5) {
			if getBodyPercentage(shortInterval.Bands[bandLen-1]) < 95 || !isLastBandPercentMoreThan10AndJustOnce(shortInterval.Bands) {
				if shortInterval.PriceChanges < 10 || isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-4:]) || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-4:]) || isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-4:]) {
					ignoredReason = "long and mid head more than body upper"
					return true
				}
			}
		}

		if countDownBand(shortInterval.Bands[len(shortInterval.Bands)-4:]) > 1 && countDownBand(shortInterval.Bands[len(shortInterval.Bands)-2:]) == 1 {
			ignoredReason = "more than body upper && short up down"
			return true
		}

		if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
			if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen/2:]) == 1 {
				ignoredReason = "more than body upper && short cross upper but just one"
				return true
			}
		}
	}

	if longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
		if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], midCloseBandAverage) != 1 || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-5:], shortCloseBandAverage) != 1 {
			if longInterval.Position != models.ABOVE_UPPER || midInterval.Position != models.ABOVE_UPPER || shortInterval.Position != models.ABOVE_UPPER {
				ignoredReason = "second trend down"
				return true
			}
		}
	}

	if isContainHeadMoreThanBody(midInterval.Bands[len(midInterval.Bands)-2:len(midInterval.Bands)-1]) && countCrossUpper(midInterval.Bands[bandLen-4:]) > 1 && !longSignificanUpAndJustOne(longInterval.Bands) {
		if isBandHeadDoubleBody(midInterval.Bands[len(midInterval.Bands)-2 : len(midInterval.Bands)-1]) {
			ignoredReason = "head double body"
			return true
		}

		if !isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-2:len(midInterval.Bands)-1], false) && midInterval.Bands[len(midInterval.Bands)-2].Candle.Close > midInterval.Bands[len(midInterval.Bands)-2].Candle.Open {
			if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], 5) != 1 {
				ignoredReason = "more than body and cross upper head only"
				return true
			}
		}
	}

	if countAboveUpper(longInterval.Bands[bandLen-2:]) > 0 && countAboveUpper(shortInterval.Bands[bandLen-4:]) > 0 {
		ignoredReason = "contain open close above upper"
		return true
	}

	if countAboveUpper(longInterval.Bands[bandLen-1:]) > 0 && countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-1:], 3) == 0 {
		ignoredReason = "above upper"
		return true
	}

	if isHasBadBand(longInterval.Bands[bandLen-2:]) && isHasBadBand(midInterval.Bands[bandLen-2:]) && isHasBadBand(shortInterval.Bands[bandLen-2:]) {
		ignoredReason = "contain bad bands"
		return true
	}

	if longInterval.AllTrend.SecondTrend == models.TREND_UP && longInterval.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-2:], 2) == 0 {
				ignoredReason = "on down, no significan up"
				return true
			}
		}
	}

	if countBadBands(longInterval.Bands[bandLen-4:]) > 2 {
		if longInterval.Position == models.ABOVE_UPPER {
			if midInterval.Position == models.ABOVE_UPPER && isHasBadBand(midInterval.Bands[bandLen-2:bandLen-1]) {
				if shortInterval.Position == models.ABOVE_UPPER || shortPercentFromUpper < 4 {
					ignoredReason = "above upper, bad bands more than 2, mid has bad band"
					return true
				}
			}
		}

		if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) > 1 {
			if isHasBandDownFromUpper(midInterval.Bands[bandLen-3:]) || countBadBands(shortInterval.Bands[bandLen-4:]) > 2 {
				if shortInterval.Position == models.ABOVE_UPPER || shortPercentFromUpper < 4 {
					ignoredReason = " bad bands more than 2, mid has down from upper"
					return true
				}
			}
		}

		if longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) > 1 {
				if countDownBand(shortInterval.Bands[bandLen-2:]) > 0 {
					if shortInterval.Position == models.ABOVE_UPPER || shortPercentFromUpper < 4 {
						ignoredReason = "trend down, bad bands more than 2"
						return true
					}
				}
			}
		}

		if isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-2:]) {
			if midInterval.Position == models.ABOVE_SMA && midPercentFromUpper < 5 {
				if shortInterval.Position == models.ABOVE_UPPER && isHasBandDownFromUpper(shortInterval.Bands[bandLen-3:]) {
					ignoredReason = "has open close above upper and  short has ban down from upper"
					return true
				}
			}
		}
	}

	if longInterval.AllTrend.ShortTrend != models.TREND_UP {
		if countDownBand(longInterval.Bands[bandLen-3:]) > 1 {
			if midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				if shortInterval.Position == models.ABOVE_UPPER && countBadBands(shortInterval.Bands[bandLen-2:]) == 1 {
					ignoredReason = "short trend not up, mid second trend down"
					return true
				}
			}
		}

		if longInterval.Position == models.ABOVE_SMA && isHasCrossUpper(longInterval.Bands[bandLen-4:], true) {
			if midInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
				if midPercentFromUpper < 4 && longPercentFromUpper < 4 {
					if shortInterval.Position == models.ABOVE_UPPER {
						ignoredReason = "above sma and percent from upper below 4, short above upper"
						return true
					}
				}
			}
		}
	}

	if isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-4:]) || isHasBandDownFromUpper(longInterval.Bands[bandLen-4:]) {
		if !isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
			if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				ignoredReason = "short trend not up, contain open close above upper, short interval upper head more than body"
				return true
			}

			if !isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) && shortPercentFromUpper < 3 {
				ignoredReason = "short trend not up, contain open close above upper, short interval percent from upper below 3"
				return true
			}
		}
	}

	if longInterval.AllTrend.FirstTrend == models.TREND_UP && longInterval.AllTrend.SecondTrend == models.TREND_UP {
		if !isHasCrossUpper(longInterval.Bands[bandLen-4:], true) && longPercentFromUpper < 6 && longInterval.PriceChanges > 10 {
			if midInterval.Position == models.ABOVE_UPPER && shortInterval.Position == models.ABOVE_UPPER {
				ignoredReason = "trend up up, not cross upper and percent from upper below 6"
				return true
			}
		}
	}

	if (isHasCrossSMA(longInterval.Bands[bandLen-1:], false) || (longInterval.Position == models.BELOW_SMA && longPercentFromSMA < 5)) && countBadBands(longInterval.Bands[bandLen-4:]) > 2 {
		if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
			if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) == 1 {
				ignoredReason = "trend down, cross sma, mid and short cross upper just one"
				return true
			}
		}

		if midInterval.Position == models.ABOVE_SMA && midPercentFromUpper < 5 {
			if isHasCrossUpper(shortInterval.Bands[bandLen-4:], false) || shortPercentFromUpper < 4 {
				ignoredReason = "trend down, below sma, mid and short percent below 5"
				return true
			}
		}
	}

	if isHasOpenCloseAboveUpper(longInterval.Bands[bandLen-1:]) && isHasBadBand(longInterval.Bands[bandLen-2:]) {
		ignoredReason = "open close above upper and bad bands"
		return true
	}

	if isUpperHeadMoreThanUpperBody(longInterval.Bands[bandLen-1]) && countBadBands(longInterval.Bands[bandLen-4:]) > 1 {
		if midInterval.Position == models.ABOVE_UPPER && countBadBands(midInterval.Bands[bandLen-3:]) > 1 {
			ignoredReason = "upper head more than upper body, contain band bands"
			return true
		}
	}

	if (countCrossUpper(longInterval.Bands[bandLen-4:]) == 1 && isHasCrossUpper(longInterval.Bands[bandLen-1:], true)) || countBadBands(longInterval.Bands[bandLen-4:]) > 2 {
		if midInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
			if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				ignoredReason = "mid cross upper just one and short upper head more than upper body"
				return true
			}
		}
	}

	if isHasCrossUpper(longInterval.Bands[bandLen-2:], true) || isHasCrossSMA(longInterval.Bands[bandLen-2:], false) {
		if countBadBands(midInterval.Bands[bandLen-4:]) > 2 {
			if shortInterval.Position == models.ABOVE_UPPER || shortPercentFromUpper < 3.5 {
				ignoredReason = "mid interval contain 3 bad bands"
				return true
			}
		}
	}

	if countBadBands(longInterval.Bands[bandLen-4:]) > 2 {
		if isHasBandDownFromUpper(longInterval.Bands[bandLen-4:]) {
			if isHasCrossSMA(midInterval.Bands[bandLen-4:], false) {
				if isHasCrossUpper(shortInterval.Bands[bandLen-2:], false) {
					ignoredReason = "contain 3 bad bands, mid cross sma"
					return true
				}
			}
		}

		if countDownBand(midInterval.Bands[bandLen-4:]) > 1 && !(isHasCrossSMA(midInterval.Bands[bandLen-4:], false) || isHasCrossLower(midInterval.Bands[bandLen-4:], false)) {
			if countBadBands(shortInterval.Bands[bandLen-4:]) > 2 {
				ignoredReason = "long and short contain 3 bad bands"
				return true
			}
		}
	}

	if isHasCrossSMA(longInterval.Bands[bandLen-1:], false) {
		if isHasCrossSMA(midInterval.Bands[bandLen-1:], false) {
			if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
				ignoredReason = "cross sma and upper head more than upper body"
				return true
			}
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-4:]) > 1 && isHasBandDownFromUpper(shortInterval.Bands[bandLen-4:]) {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && isHasCrossSMA(midInterval.Bands[bandLen-1:], true) && longInterval.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "short trend not up, mid second trend down2"
			return true
		}
	}

	if shortInterval.AllTrend.Trend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
		if longInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if getHigestPercentChangesIndex(shortInterval.Bands[bandLen-4:]) != 3 {
				ignoredReason = "trend down and last band not higest"
				return true
			}
		}
	}

	return false
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
	if (bandDoublePreviousHigh(short.Bands, 2.3) || bandDoublePreviousHigh(mid.Bands, 2.5)) && mid.AllTrend.ShortTrend == models.TREND_UP {
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
								if !isHasBadBand(short.Bands[bandLen-1:]) && shortLastBandPercent > 3.2 && countBandPercentChangesMoreThan(short.Bands[:bandLen-1], 1) == 0 {
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

func closeBandAverage(bands []models.Band) float32 {
	var sumarize float32 = 0
	for _, band := range bands {
		sumarize += band.Candle.Close
	}

	return sumarize / float32(len(bands))
}

func getBodyPercentage(band models.Band) float32 {
	body := band.Candle.Close - band.Candle.Open
	all := band.Candle.Hight - band.Candle.Low

	return body / all * 100
}

func isLastBandPercentMoreThan10AndJustOnce(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	percent := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Close * 100

	return percent > 10 && countBandPercentChangesMoreThan(bands[len(bands)-4:], 5) == 1
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
