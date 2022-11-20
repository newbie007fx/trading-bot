package analysis

import (
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""

func lowestFromBand(band models.Band) float32 {
	if band.Candle.Open > band.Candle.Close {
		return band.Candle.Close
	}

	return band.Candle.Open
}

func GetIgnoredReason() string {
	return ignoredReason
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

func isLastBandDoublePreviousHeigest(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	lastBandBodyHeight := lastBand.Candle.Close - lastBand.Candle.Open

	var higestBody float32 = 0
	for _, band := range bands[len(bands)-5 : len(bands)-1] {
		bodyHeight := band.Candle.Close - band.Candle.Open
		if bodyHeight > higestBody {
			higestBody = bodyHeight
		}
	}

	return higestBody*2 < lastBandBodyHeight
}

func bandPercent(band models.Band) float32 {
	return (band.Candle.Close - band.Candle.Open) / band.Candle.Open * 100
}

func isLastBandHeigestBand(bands []models.Band, count int) bool {
	lastBand := bands[len(bands)-1]
	lastBandClose := lastBand.Candle.Close

	for _, band := range bands[len(bands)-count : len(bands)-1] {
		if band.Candle.Close > lastBandClose {
			return false
		}
	}

	return true
}

func isOpenCloseAboveUpper(band models.Band) bool {
	return band.Candle.Open > float32(band.Upper) && band.Candle.Close > float32(band.Upper)
}

func isHeadMoreThanBody(band models.Band) bool {
	head := band.Candle.Hight - band.Candle.Close
	body := band.Candle.Close - band.Candle.Open

	return head > body
}

func isUpperHeadMoreThanUpperBody(band models.Band) bool {
	if band.Candle.Open > band.Candle.Close {
		return false
	}

	if !(band.Candle.Open < float32(band.Upper) && band.Candle.Close > float32(band.Upper)) {
		return false
	}

	upperToHead := band.Candle.Close - float32(band.Upper)
	upperToBody := float32(band.Upper) - band.Candle.Open

	return upperToHead > upperToBody
}

func countUpperHeadMoreThanUpperBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isUpperHeadMoreThanUpperBody(band) {
			count++
		}
	}
	return count
}

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time) bool {
	ignoredReason = ""

	bandLen := len(short.Bands)
	shortLastBand := short.Bands[bandLen-1]
	shortSecondLastBand := short.Bands[bandLen-2]
	midLastBand := mid.Bands[bandLen-1]
	midSecondLastBand := mid.Bands[bandLen-2]

	if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
		if !isHeadMoreThanBody(midSecondLastBand) && !isOpenCloseAboveUpper(midLastBand) {
			ignoredReason = "pattern 1"
			return true
		}
	}

	if isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 3 && isLastBandHeigestBand(short.Bands, 4) {
		if !isOpenCloseAboveUpper(midLastBand) && !isHeadMoreThanBody(midSecondLastBand) {
			if !(isOpenCloseAboveUpper(shortLastBand) && countUpperHeadMoreThanUpperBody(short.Bands[bandLen-3:bandLen-1]) > 1) {
				ignoredReason = "pattern 2"
				return true
			}
		}
	}

	if shortSecondLastBand.Candle.Close > shortLastBand.Candle.Open {
		if isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) && bandPercent(shortSecondLastBand) > 1.5 {
			if !isOpenCloseAboveUpper(midLastBand) && !isHeadMoreThanBody(midSecondLastBand) {
				ignoredReason = "pattern 3"
				return true
			}
		}
	}

	if midSecondLastBand.Candle.Close > midSecondLastBand.Candle.Open {
		if isLastBandDoublePreviousHeigest(mid.Bands[:bandLen-1]) && bandPercent(midSecondLastBand) > 3 {
			if !isOpenCloseAboveUpper(midLastBand) && !isHeadMoreThanBody(midSecondLastBand) {
				if isLastBandHeigestBand(short.Bands[:bandLen-1], 4) || isLastBandHeigestBand(short.Bands, 4) {
					ignoredReason = "pattern 4"
					return true
				}
			}
		}
	}

	return false
}
