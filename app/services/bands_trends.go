package services

import "telebot-trading/app/models"

func CalculateTrends(data []models.Band) int8 {
	if len(data) < 7 {
		return 0
	}

	highestIndex, lowestIndex := 0, 0
	var total float32 = 0
	for i, val := range data {
		total += val.Candle.Close

		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}
	}

	if highestIndex == len(data)-1 {
		return models.TREND_UP
	}
	if lowestIndex == len(data)-1 {
		return models.TREND_DOWN
	}

	average := total / float32(len(data))
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

	lastCandle := data[len(data)-1].Candle
	if lastCandle.Close > average {
		return models.TREND_UP
	} else {
		return models.TREND_DOWN
	}
}
