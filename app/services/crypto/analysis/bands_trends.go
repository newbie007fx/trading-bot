package analysis

import "telebot-trading/app/models"

func CalculateTrends(data []models.Band) int8 {
	lastCandle := data[len(data)-1].Candle

	highestIndex, lowestIndex := 0, 0
	hightAverage, lowAverage := 0, 0
	var total float32 = 0
	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}

		avg := (val.Candle.Open + val.Candle.Close) / 2
		if hightAverage < avg {
			hightAverage = avg
		}

		if lowAverage > avg || lowAverage == 0 {
			lowAverage = avg
		}

		if i < len(data)-1 {
			total += avg
		}
	}

	average := total / float32(len(data)-2)

	diffHight := hightAverage - average
	diffLow := average - lowAverage
	var percent float32 = 0
	if diffHight > diffLow {
		percent = (diffLow / diffHight) * 100
	} else {
		percent = (diffHight / diffLow) * 100
	}

	if percent > float32(75) {
		return models.TREND_SIDEWAY
	}

	averageTrend := (data[lowestIndex].Candle.Close + data[highestIndex].Candle.Close) / 2
	if lastCandle.Close > averageTrend {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
