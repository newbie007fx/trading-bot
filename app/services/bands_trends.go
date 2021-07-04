package services

import "telebot-trading/app/models"

func CalculateTrends(data []models.Band) int8 {
	if len(data) < 8 {
		return 0
	}

	lastCandle := data[len(data)-1].Candle

	highestIndex, lowestIndex := 0, 0
	var total float32 = 0
	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}

		if i < len(data)-1 {
			total += (val.Candle.Close + val.Candle.Open) / 2
		}
	}

	if highestIndex == len(data)-1 {
		return models.TREND_UP
	}
	if lowestIndex == len(data)-1 {
		return models.TREND_DOWN
	}

	average := total / float32(len(data)-1)
	highestValueDifference := data[highestIndex].Candle.Close - average
	lowestValueDifference := average - data[highestIndex].Candle.Close

	var percent float32 = 0
	if highestValueDifference > lowestValueDifference {
		difference := highestValueDifference - lowestValueDifference
		percent = (difference / highestValueDifference) * 100
	} else {
		difference := lowestValueDifference - highestValueDifference
		percent = (difference / lowestValueDifference) * 100
	}

	if percent <= float32(20) {
		return models.TREND_SIDEWAY
	}

	if lastCandle.Close > average {
		return models.TREND_UP
	} else {
		return models.TREND_DOWN
	}
}
