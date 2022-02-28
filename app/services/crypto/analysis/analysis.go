package analysis

import (
	"math"
	"telebot-trading/app/models"
	"time"
)

const Time_type_15m int = 1
const Time_type_1h int = 2
const Time_type_4h int = 3

func GetCurrentBollingerBands(candlesData []models.CandleData) (bands models.Bands) {
	bands = GenerateBollingerBands(candlesData)
	direction := BAND_DOWN
	if CheckLastCandleIsUp(bands.Data) {
		direction = BAND_UP
	}

	bands.PriceChanges = CalculateBandPriceChangesPercent(bands, direction)

	return
}

func GetHighestHightPriceByTime(currentTime time.Time, bands []models.Band, timeType int) float32 {
	var numberBands int = 0
	var utcZone, _ = time.LoadLocation("UTC")
	currentTime = currentTime.In(utcZone)

	if timeType == Time_type_15m {
		numberBands = (currentTime.Minute() + 1) % 15
		if numberBands == 0 {
			numberBands = 15
		}
	} else if timeType == Time_type_1h {
		numberBands = int(math.Ceil(float64(currentTime.Minute()+1) / 15))
	} else {
		numberBands = (currentTime.Hour() + 1) % 4
		if numberBands == 0 {
			numberBands = 4
		}
	}

	return GetHigestHightPrice(bands[len(bands)-numberBands:])
}

func GetLowestLowPriceByTime(currentTime time.Time, bands []models.Band, timeType int) float32 {
	var numberBands int = 0
	var utcZone, _ = time.LoadLocation("UTC")
	currentTime = currentTime.In(utcZone)

	if timeType == Time_type_15m {
		numberBands = (currentTime.Minute() + 1) % 15
		if numberBands == 0 {
			numberBands = 15
		}
	} else if timeType == Time_type_1h {
		numberBands = int(math.Ceil(float64(currentTime.Minute()+1) / 15))
	} else {
		numberBands = (currentTime.Hour() + 1) % 4
		if numberBands == 0 {
			numberBands = 4
		}
	}

	return GetLowestLowPrice(bands[len(bands)-numberBands:])
}
