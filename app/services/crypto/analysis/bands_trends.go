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

	if highestIndex == len(data)-1 {
		return models.TREND_UP
	}

	if lowestIndex == len(data)-1 {
		return models.TREND_DOWN
	}

	firstSMA := data[0].SMA
	lastSMA := data[len(data)-1].SMA

	var percent float64 = 0
	if firstSMA > lastSMA {
		percent = (lastSMA / firstSMA) * 100
	} else {
		percent = (firstSMA / lastSMA) * 100
	}

	if percent >= float64(93) {
		return models.TREND_SIDEWAY
	}

	averageTrend := (data[lowestIndex].Candle.Close + data[highestIndex].Candle.Close) / 2
	if lastCandle.Close > averageTrend {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
