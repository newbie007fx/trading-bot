package analysis

import (
	"log"
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""

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
		if band.Candle.Open > float32(band.Upper) && band.Candle.Close < float32(band.Upper) {
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
		if band.Candle.Close > band.Candle.Open && band.Candle.Open > float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
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

func countCrossUpperOnBody(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Open < float32(data.Upper) && data.Candle.Close > float32(data.Upper) && (data.Candle.Close > data.Candle.Open) {
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

func countHeadDoubleBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			head := band.Candle.Hight - band.Candle.Close
			body := band.Candle.Close - band.Candle.Open
			if head > body*1.8 {
				count++
			}
		}
	}

	return count
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
	if isHasOpenCloseAboveUpper(shortInterval.Bands[len(shortInterval.Bands)-1:]) {
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

	return false
}

func IgnoredOnUpTrendMid(midInterval, shortInterval models.BandResult) bool {
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

	return false
}

func IgnoredOnUpTrendLong(longInterval, midInterval, shortInterval models.BandResult, checkTime time.Time) bool {

	bandLen := len(longInterval.Bands)
	longLastBand := longInterval.Bands[bandLen-1]
	midLastBand := midInterval.Bands[bandLen-1]

	midPercentFromUpper := (midLastBand.Upper - float64(midLastBand.Candle.Close)) / float64(midLastBand.Candle.Close) * 100

	longPercentFromSMA := (float32(longLastBand.SMA) - longLastBand.Candle.Close) / longLastBand.Candle.Close * 100

	shortLastBand := shortInterval.Bands[bandLen-1]
	shortPercentFromUpper := (float32(shortLastBand.Upper) - shortLastBand.Candle.Close) / shortLastBand.Candle.Close * 100

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

	if approvedPattern(shortInterval, midInterval, longInterval) {
		return false
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
				if isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
					ignoredReason = "above upper and short upper head more than body"
					return true
				}

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

		if countCrossUpper(midInterval.Bands[bandLen-4:]) <= 1 {
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

		longPercentFromUpper := (float32(longLastBand.Upper) - longLastBand.Candle.Close) / longLastBand.Candle.Close * 100
		if isHasCrossSMA(longInterval.Bands[bandLen-1:], false) || longPercentFromUpper < 3 {
			if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-1:]) || isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-1:]) || (isHasUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1:]) && checkTime.Minute() > 18 && checkTime.Minute() < 3) {
				if shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "open close above upper or head upper more than body"
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

		if isHasBadBand(longInterval.Bands[bandLen-2:]) {
			if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-1:]) || isHasUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1:]) || isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-1:]) || (isHasUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1:]) && checkTime.Minute() > 18 && checkTime.Minute() < 3) {
				ignoredReason = "above sma, contain bad band and mid or short above upper or uppper head more than body"
				return true
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
				log.Println("coba")
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
			// possibly check if short/mid interval cross upper
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

	return false
}

func isHasBadBand(bands []models.Band) bool {
	if countDownBand(bands) > 0 || countAboveUpper(bands) > 0 || isTailMoreThan(bands[0], 40) || countHeadMoreThanBody(bands) > 0 {
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
	if isHasCrossUpper(bands[len(bands)-1:], false) && lastPercent > 5 && lastPercent < 17 && countBandPercentChangesMoreThan(bands[len(bands)-4:], 5) == 1 {
		return true
	}

	return false
}

func allIntervalCrossUpperOnBodyMoreThanThresholdAndJustOne(short, mid, long models.BandResult) bool {
	if short.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER {
		if longSignificanUpAndJustOne(long.Bands) && countBandPercentChangesMoreThan(long.Bands[len(long.Bands)-5:], 3) > 1 {
			if !isHasUpperHeadMoreThanUpperBody(short.Bands[len(short.Bands)-1:]) || countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-1:], 6) != 1 || !isHasUpperHeadMoreThanUpperBody(mid.Bands[len(mid.Bands)-1:]) || countBandPercentChangesMoreThan(mid.Bands[len(mid.Bands)-1:], 10) != 1 {
				if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 3) == 1 && getHigestPercentChangesIndex(short.Bands[len(short.Bands)-4:]) == 0 && countDownBand(short.Bands[len(short.Bands)-4:]) > 1 {
					return false
				}

				if long.Bands[len(long.Bands)-1].Candle.Low < float32(long.Bands[len(long.Bands)-1].Lower) && long.Bands[len(long.Bands)-1].Candle.Close > float32(long.Bands[len(long.Bands)-1].Upper) {
					return false
				}

				if isTailMoreThan(short.Bands[len(short.Bands)-1], 40) {
					return false
				}

				if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 3) > 1 && short.Bands[len(short.Bands)-2].Candle.Open > short.Bands[len(short.Bands)-2].Candle.Close {
					return false
				}

				if countBadBands(short.Bands[len(short.Bands)-4:len(short.Bands)-1]) == 3 {
					return false
				}

				if countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) > 1 && countCrossUpperOnBody(short.Bands[len(short.Bands)-4:]) > 1 && isUpperHeadMoreThanUpperBody(mid.Bands[len(mid.Bands)-1]) {
					return false
				}

				if isHasOpenCloseAboveUpper(mid.Bands[len(mid.Bands)-2:]) && getHigestPercentChangesIndex(short.Bands[len(short.Bands)-4:]) == 3 {
					if countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 5) == 1 && countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 2) == 1 {
						return false
					}
				}

				if countCrossUpperOnBody(long.Bands[len(long.Bands)-4:]) == 1 && countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) == 1 {
					if short.PriceChanges > 6 {
						return false
					}
				}

				if (isHasUpperHeadMoreThanUpperBody(mid.Bands[len(mid.Bands)-2:]) || isHasOpenCloseAboveUpper(mid.Bands[len(mid.Bands)-2:])) && (isHasUpperHeadMoreThanUpperBody(short.Bands[len(short.Bands)-2:]) || isHasOpenCloseAboveUpper(short.Bands[len(short.Bands)-2:])) {
					return false
				}

				if isHasBandDownFromUpper(long.Bands[len(long.Bands)-2:]) {
					if countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) == 1 {
						if isUpperHeadMoreThanUpperBody(short.Bands[len(short.Bands)-1]) {
							return false
						}
					}
				}

				if mid.Position == models.ABOVE_UPPER && countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) == 1 {
					if isUpperHeadMoreThanUpperBody(short.Bands[len(short.Bands)-1]) {
						return false
					}
				}

				if countCrossUpperOnBody(mid.Bands[len(mid.Bands)-4:]) == 1 && countBadBands(mid.Bands[len(mid.Bands)-4:]) > 2 {
					return false
				}

				if ((countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 3) >= 1 && countBandPercentChangesMoreThan(short.Bands[len(short.Bands)-4:], 2) > 1) || (countBandPercentChangesMoreThan(mid.Bands[len(mid.Bands)-4:], 5) == 1 && !isBandHeadDoubleBody(mid.Bands[len(mid.Bands)-2:]))) && countBandPercentChangesMoreThan(short.Bands[len(mid.Bands)-4:], 5) < 2 {
					if !isHasOpenCloseAboveUpper(short.Bands[len(short.Bands)-4:]) && !isBandHeadDoubleBody(long.Bands[len(long.Bands)-2:len(long.Bands)-1]) {
						return true
					}
				}
			}
		}
	}

	return false
}

func approvedPattern(short, mid, long models.BandResult) bool {
	if allIntervalCrossUpperOnBodyMoreThanThresholdAndJustOne(short, mid, long) {
		log.Println("skipped1")
		return true
	}

	return false
}

func isUpperHeadMoreThanUpperBody(band models.Band) bool {
	allBody := band.Candle.Close - band.Candle.Open
	head := band.Candle.Close - float32(band.Upper)
	percent := head / allBody * 100
	return percent > 55
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
