package analysis

import "telebot-trading/app/models"

func CalculateTrends(data []models.Band) int8 {
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
			total += (val.Candle.Open + val.Candle.Close) / 2
		}
	}

	average := total / float32(len(data)-2)

	var percent float32 = 0
	if lastCandle.Close > average {
		difference := lastCandle.Close - average
		percent = (difference / average) * 100
	} else {
		difference := average - lastCandle.Close
		percent = (difference / lastCandle.Close) * 100
	}

	if percent <= float32(3) {
		return models.TREND_SIDEWAY
	}

	averageTrend := (data[lowestIndex].Candle.Close + data[highestIndex].Candle.Close) / 2
	if lastCandle.Close > averageTrend {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
