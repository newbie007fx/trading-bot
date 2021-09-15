package analysis

import "telebot-trading/app/models"

func IsIgnored(result *models.BandResult) bool {
	return isInAboveUpperBandAndDownTrend(result)
}

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	lastFiveData := result.Bands[len(result.Bands)-5:]
	if isHeighestOnHalfEndAndAboveUpper(result) && CalculateTrends(lastFiveData) == models.TREND_DOWN {
		return true
	}

	return isContaineBearishEngulfing(result)
}

func isHeighestOnHalfEndAndAboveUpper(result *models.BandResult) bool {
	hiIndex := getHighestIndex(result)
	if hiIndex >= len(result.Bands)-5 {
		return result.Bands[hiIndex].Candle.Close > float32(result.Bands[hiIndex].Upper)
	}

	return false
}

func isContaineBearishEngulfing(result *models.BandResult) bool {
	hiIndex := getHighestIndex(result)
	if hiIndex > len(result.Bands)/2 {
		return BearishEngulfing(result.Bands[hiIndex:])
	}

	return false
}

func getHighestIndex(result *models.BandResult) int {
	hiIndex := 0
	for i, band := range result.Bands {
		if result.Bands[hiIndex].Candle.Close < band.Candle.Close {
			hiIndex = i
		}
	}

	return hiIndex
}
