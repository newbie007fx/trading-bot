package analysis

import "telebot-trading/app/models"

func IsIgnored(result *models.BandResult) bool {
	return isInAboveUpperBandAndDownTrend(result)
}

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	lastFourData := result.Bands[len(result.Bands)-4 : len(result.Bands)]

	return int8(lastBand.Candle.Low) > int8(lastBand.Upper) && CalculateTrends(lastFourData) == models.TREND_DOWN
}
