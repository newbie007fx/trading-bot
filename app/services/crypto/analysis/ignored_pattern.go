package analysis

import "telebot-trading/app/models"

func IsIgnored(result *models.BandResult) bool {
	return isInAboveUpperBandAndDownTrend(result)
}

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	lastFiveData := result.Bands[len(result.Bands)-5:]

	return isHeighestOnHalfEndAndAboveUpper(result) && CalculateTrends(lastFiveData) == models.TREND_DOWN
}

func isHeighestOnHalfEndAndAboveUpper(result *models.BandResult) bool {
	hiIndex := 0
	for i, band := range result.Bands {
		if result.Bands[hiIndex].Candle.Close < band.Candle.Close {
			hiIndex = i
		}
	}
	if hiIndex < len(result.Bands)-5 {
		return false
	}

	return result.Bands[hiIndex].Candle.Close > float32(result.Bands[hiIndex].Upper)
}
