package analysis

import (
	"telebot-trading/app/models"
)

func GetCurrentBollingerBands(candlesData []models.CandleData) (bands models.Bands) {
	bands = GenerateBollingerBands(candlesData)
	direction := BAND_DOWN
	if CheckLastCandleIsUp(bands.Data) {
		direction = BAND_UP
	}

	bands.PriceChanges = CalculateBandPriceChangesPercent(bands, direction)
	bands.VolumeAverageChanges = CalculateVolumeAverage(bands.Data)

	return
}
