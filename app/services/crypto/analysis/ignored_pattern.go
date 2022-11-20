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

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time) bool {
	ignoredReason = ""

	bandLen := len(short.Bands)
	shortLastBand := short.Bands[bandLen-1]
	shortSecondLastBand := short.Bands[bandLen-2]
	midLastBand := mid.Bands[bandLen-1]
	midSecondLastBand := mid.Bands[bandLen-2]

	if (isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5) || (isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 3 && isLastBandHeigestBand(short.Bands, 4)) {
		ignoredReason = "pattern 1"
		return true
	}

	if (isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) && bandPercent(shortSecondLastBand) > 1.5) || (isLastBandDoublePreviousHeigest(mid.Bands[:bandLen-1]) && bandPercent(midSecondLastBand) > 3 && (isLastBandHeigestBand(short.Bands[:bandLen-1], 4) || isLastBandHeigestBand(short.Bands, 4))) {
		ignoredReason = "pattern 2"
		return true
	}

	return false
}
