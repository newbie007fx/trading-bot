package analysis

import "telebot-trading/app/models"

func CalculateTrends(data []models.Band) int8 {
	highestIndex, lowestIndex := 0, 0
	var totalSMA float64 = 0
	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}

		if i < len(data)-1 {
			totalSMA += val.SMA
		}
	}

	if highestIndex == len(data)-1 {
		return models.TREND_UP
	}

	if lowestIndex == len(data)-1 {
		return models.TREND_DOWN
	}

	avgSMA := totalSMA / float64(len(data)-1)
	lastSMA := data[len(data)-1].SMA

	firstAvg := (data[0].Candle.Open + data[0].Candle.Close) / 2
	lastAvg := (data[len(data)-1].Candle.Open + data[len(data)-1].Candle.Close) / 2

	var percent float64 = 0
	if avgSMA > lastSMA {
		percent = (lastSMA / avgSMA) * 100
	} else {
		percent = (avgSMA / lastSMA) * 100
	}

	if percent >= float64(99.55) {
		return models.TREND_SIDEWAY
	}

	if firstAvg < lastAvg {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
