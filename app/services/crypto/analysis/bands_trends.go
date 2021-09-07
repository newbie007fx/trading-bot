package analysis

import (
	"math"
	"telebot-trading/app/models"
)

func CalculateTrends(data []models.Band) int8 {
	highestIndex, lowestIndex := 0, 0
	thirtyPercent := float64(len(data)) * float64(30) / float64(100)
	limit := int(math.Floor(thirtyPercent))

	var totalFirstData float32 = 0
	var totalLastData float32 = 0
	var totalBaseLine float32 = 0
	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}

		if i < limit {
			totalFirstData += val.Candle.Close
			totalBaseLine += (val.Candle.Open + val.Candle.Close) / 2
		}

		if i >= len(data)-limit {
			totalLastData += val.Candle.Close
		}
	}

	if highestIndex == len(data)-1 {
		return models.TREND_UP
	}

	if lowestIndex == len(data)-1 {
		return models.TREND_DOWN
	}

	firstAvg := totalFirstData / float32(limit)
	lastAvg := totalLastData / float32(limit)
	baseLinePoint := totalBaseLine / float32(limit)

	var lastPointValue float32 = 0
	var firstPointValue float32 = 0
	if firstAvg > baseLinePoint {
		firstPointValue = firstAvg - baseLinePoint
	} else {
		firstPointValue = baseLinePoint - firstAvg
	}

	if lastAvg > baseLinePoint {
		lastPointValue = lastAvg - baseLinePoint
	} else {
		lastPointValue = baseLinePoint - lastAvg
	}

	var percent float32 = 0
	if firstPointValue > lastPointValue {
		percent = (lastPointValue / firstPointValue) * 100
	} else {
		percent = (firstPointValue / lastPointValue) * 100
	}

	if percent >= 2.2 {
		return models.TREND_SIDEWAY
	}

	if firstAvg < lastAvg {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
