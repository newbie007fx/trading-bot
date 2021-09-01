package analysis

import "telebot-trading/app/models"

func CalculateTrends(data []models.Band) int8 {
	lastCandle := data[len(data)-1].Candle

	highestIndex, lowestIndex := 0, 0
	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}
	}

	average := (data[lowestIndex].Candle.Close + data[highestIndex].Candle.Close) / 2
	highestValueDifference := data[highestIndex].Candle.Close - average
	lowestValueDifference := average - data[lowestIndex].Candle.Close

	var percent float32 = 0
	if highestValueDifference > lowestValueDifference {
		difference := highestValueDifference - lowestValueDifference
		percent = (difference / highestValueDifference) * 100
	} else {
		difference := lowestValueDifference - highestValueDifference
		percent = (difference / lowestValueDifference) * 100
	}

	if percent <= float32(25) {
		return models.TREND_SIDEWAY
	}

	if lastCandle.Close > average {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
