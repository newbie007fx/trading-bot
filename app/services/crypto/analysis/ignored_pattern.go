package analysis

import (
	"log"
	"telebot-trading/app/models"
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
		if data.Candle.Open <= float32(data.Upper) && data.Candle.Hight > float32(data.Upper) {
			count++
		}
	}

	return count
}

func countCrossUpperOnBody(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Open < float32(data.Upper) && data.Candle.Close > float32(data.Upper) {
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
	if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-2:len(midInterval.Bands)-1]) && (countHeadMoreThanBody(midInterval.Bands[bandLen-2:bandLen-1]) == 1 || countDownBand(midInterval.Bands[bandLen-2:]) == 1) {
		ignoredReason = "contain open close above upper"
		return true
	}

	if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-1:]) {
		if isHasOpenCloseAboveUpper(shortInterval.Bands[bandLen-2:]) || isUpperHeadMoreThanUpperBody(shortInterval.Bands[bandLen-1]) {
			ignoredReason = "short and mid contain open close above upper"
			return true
		}
	}

	higestIndex := getHighestIndex(midInterval.Bands[len(midInterval.Bands)-10:])
	if higestIndex != len(midInterval.Bands[len(midInterval.Bands)-10:])-1 {
		ignoredReason = "not in higest"
		return true
	}

	higestHightIndex := getHighestHightIndex(midInterval.Bands[len(midInterval.Bands)-5:])
	if higestHightIndex < len(midInterval.Bands[len(midInterval.Bands)-5:])-2 {
		ignoredReason = "previous band have higest high"
		return true
	}

	if higestHightIndex != len(midInterval.Bands[len(midInterval.Bands)-5:])-1 && countDownBand(midInterval.Bands[len(midInterval.Bands)-2:]) == 1 {
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

	if midInterval.Position == models.ABOVE_SMA {
		if midInterval.AllTrend.SecondTrend != models.TREND_UP && countCrossUpper(midInterval.Bands[bandLen-5:]) == 0 {
			if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen/2:]) == 1 {
				ignoredReason = "above sma, second trend up, short cross upper"
				return true
			}
		}
	}

	return false
}

func IgnoredOnUpTrendLong(longInterval, midInterval, shortInterval models.BandResult) bool {
	bandLen := len(longInterval.Bands)
	longLastBand := longInterval.Bands[bandLen-1]

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

	higestIndex := getHighestIndex(longInterval.Bands[len(longInterval.Bands)-5:])
	if longInterval.AllTrend.SecondTrend != models.TREND_DOWN && higestIndex != len(longInterval.Bands[len(longInterval.Bands)-5:])-1 {
		ignoredReason = "not in higest"
		return true
	}

	if longInterval.Position == models.ABOVE_UPPER {
		if isUpperHeadMoreThanUpperBody(longInterval.Bands[bandLen-1]) {
			if isBandHeadDoubleBody(midInterval.Bands[bandLen-1:]) || isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
				if isHasCrossUpper(shortInterval.Bands[bandLen-1:], true) && isBandHeadDoubleBody(shortInterval.Bands[bandLen-2:]) {
					if getBodyPercentage(shortInterval.Bands[bandLen-1]) < 95 || !isLastBandPercentMoreThan10AndJustOnce(shortInterval.Bands) {
						ignoredReason = "cross upper, mid head more than body upper, short head double body"
						return true
					}
				}
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

		if isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
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
		}

		if countCrossUpperOnBody(longInterval.Bands[bandLen-4:]) == 1 {
			if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-1:]) {
				ignoredReason = "mid contain open close above upper"
				return true
			}

			if isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) && shortInterval.Position == models.ABOVE_UPPER && !longSignificanUpAndJustOne(longInterval.Bands) {
				ignoredReason = "mid and short has head more than body upper"
				return true
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

		if countCrossUpperOnBody(longInterval.Bands[bandLen-4:]) > 2 && longInterval.PriceChanges > 15 {
			if midInterval.Position == models.ABOVE_UPPER && countDownBand(midInterval.Bands[bandLen-4:]) > 0 {
				if shortInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(shortInterval.Bands[bandLen-5:]) == 1 {
					ignoredReason = "cross upper, price change alredy more than 15"
					return true
				}
			}
		}

		if countCrossUpper(longInterval.Bands[bandLen-3:]) > 1 {
			if countCrossUpper(midInterval.Bands[bandLen-4:]) > 0 {
				if countCrossUpper(shortInterval.Bands[bandLen-4:]) > 0 {
					if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) != 1 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) != 1 && countBandPercentChangesMoreThan(longInterval.Bands[bandLen-4:], 5) != 1 && !longSignificanUpAndJustOne(longInterval.Bands) {
						ignoredReason = "all band cross upper more than one but percent change is minor"
						return true
					}
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

		if isHasOpenCloseAboveUpper(midInterval.Bands[bandLen-1:]) || isUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1]) {
			if shortInterval.Position == models.ABOVE_UPPER {
				ignoredReason = "open close above upper or head upper more than body"
				return true
			}
		}

		longPercentFromUpper := (float32(longLastBand.Upper) - longLastBand.Candle.Close) / longLastBand.Candle.Close * 100
		if isHasCrossSMA(longInterval.Bands[bandLen-5:], false) && isHasCrossUpper(longInterval.Bands[bandLen/2:bandLen-5], true) {
			if CountUpBand(longInterval.Bands[bandLen-3:]) == 1 || isHasCrossUpper(longInterval.Bands[bandLen-3:], false) || !isHasCrossUpper(longInterval.Bands[bandLen-6:], false) || longInterval.PriceChanges < 10 {
				if longPercentFromUpper < 5 && !(isHasCrossSMA(longInterval.Bands[bandLen-2:bandLen-1], false) && countCrossSMA(longInterval.Bands[bandLen-4:]) == 1) {
					ignoredReason = "up down and below upper"
					return true
				}
			}
		}

		if !isHasCrossUpper(longInterval.Bands[bandLen-5:], true) && longInterval.PriceChanges > 10 {
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

		if isHasUpperHeadMoreThanUpperBody(midInterval.Bands[bandLen-1:]) {
			if isBandHeadDoubleBody(shortInterval.Bands[bandLen-1:]) {
				ignoredReason = "above sma, mid upper head more than body, short head double body"
				return true
			}
		}

		if isHasOpenCloseAboveUpper(midInterval.Bands[len(midInterval.Bands)-1:]) {
			ignoredReason = "above sma and mid contain open close above upper"
			return true
		}

		if midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
			if countBandPercentChangesMoreThan(midInterval.Bands[bandLen-5:], midCloseBandAverage) != 1 || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-5:], shortCloseBandAverage) != 1 {
				ignoredReason = "all interval above sma"
				return true
			}
		}

		if isHasCrossSMA(longInterval.Bands[bandLen-1:], true) && longPercentFromUpper < 9 {
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
				if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 5) == 0 && countBandPercentChangesMoreThan(midInterval.Bands[bandLen-4:], 5) == 0 {
					ignoredReason = "cross sma, mid and short cross upper more than one but have no significan up band"
					return true
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
		}

		if countCrossUpper(longInterval.Bands[bandLen-4:]) == 0 {
			if isHasCrossUpper(midInterval.Bands[bandLen-4:], true) && midInterval.PriceChanges < 10 {
				if shortInterval.PriceChanges > 5 && countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) == 0 && !(isHasCrossSMA(longInterval.Bands[bandLen-2:bandLen-1], false) && countCrossSMA(longInterval.Bands[bandLen-4:]) == 1) {
					ignoredReason = "mid cross upper, and short interval do not have significan band"
					return true
				}

				if countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 && countCrossUpper(shortInterval.Bands[bandLen-4:]) > 1 {
					if countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 5) == 0 || countBandPercentChangesMoreThan(shortInterval.Bands[bandLen-4:], 3) == 0 {
						ignoredReason = "mid cross upper, and short interval do not have significan band 2"
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
	}

	if longInterval.Position == models.BELOW_SMA {
		if countCrossSMA(longInterval.Bands[bandLen-7:]) == 0 || (countCrossSMA(longInterval.Bands[bandLen-7:]) == 1 && isHasCrossSMA(longInterval.Bands[bandLen-1:], false)) && longInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if midInterval.Position == models.ABOVE_UPPER && countCrossUpperOnBody(midInterval.Bands[bandLen-4:]) == 1 {
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
		}
	}

	higestHightIndex := getHighestHightIndex(longInterval.Bands[len(longInterval.Bands)-5:])
	if higestHightIndex < len(longInterval.Bands[len(longInterval.Bands)-5:])-2 {
		if midInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands[bandLen-4:]) == 1 {
			ignoredReason = "previous band have higest high and mid cross upper just one"
			return true
		}
	}

	midLastBand := midInterval.Bands[bandLen-1]
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
			ignoredReason = "second trend down"
			return true
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

	return false
}

func longSignificanUpAndJustOne(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	lastPercent := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100
	log.Println(lastPercent)
	if isHasCrossUpper(bands[len(bands)-1:], false) && lastPercent > 5 && lastPercent < 17 && countBandPercentChangesMoreThan(bands[len(bands)-4:], 5) == 1 {
		return true
	}

	return false
}

func isUpperHeadMoreThanUpperBody(band models.Band) bool {
	body := float32(band.Upper) - band.Candle.Open
	head := band.Candle.Close - float32(band.Upper)

	return body < head
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
