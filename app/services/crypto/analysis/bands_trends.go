package analysis

import (
	"log"
	"math"
	"telebot-trading/app/models"
)

func CalculateTrends(data []models.Band) int8 {
	highestIndex, lowestIndex := 0, 0
	thirtyPercent := float64(len(data)) * float64(30) / float64(100)
	limit := int(math.Floor(thirtyPercent))

	var totalFirstData float32 = 0
	var totalLastData float32 = 0
	var totalMidleData float32 = 0
	var totalBaseLine float32 = 0

	middleIndex := (len(data) / 2) - 1
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

		if i > middleIndex-(limit/2) && i <= middleIndex+(limit/2) {
			totalMidleData += val.Candle.Close
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
	midleAvg := totalMidleData / float32(limit)
	baseLinePoint := totalBaseLine / float32(limit)

	firstToMidleTrend := getTrend(baseLinePoint, firstAvg, midleAvg)
	midleToLastTrend := getTrend(baseLinePoint, midleAvg, lastAvg)

	log.Println(firstAvg)
	log.Println(lastAvg)
	log.Println(midleAvg)
	log.Println(baseLinePoint)

	log.Println(firstToMidleTrend)
	log.Println(midleToLastTrend)

	if firstToMidleTrend == models.TREND_SIDEWAY {
		return midleToLastTrend
	}

	if midleToLastTrend == models.TREND_SIDEWAY {
		return firstToMidleTrend
	}

	if firstToMidleTrend == models.TREND_UP && midleToLastTrend == models.TREND_UP {
		return models.TREND_UP
	}

	if firstToMidleTrend == models.TREND_DOWN && midleToLastTrend == models.TREND_DOWN {
		return models.TREND_DOWN
	}

	if firstToMidleTrend == models.TREND_UP && midleToLastTrend == models.TREND_DOWN {
		if firstAvg < lastAvg {
			fourtyFromMidle := 40 * (midleAvg - firstAvg) / 100
			if lastAvg > (midleAvg - fourtyFromMidle) {
				return models.TREND_UP
			}
		}
	}

	if firstToMidleTrend == models.TREND_DOWN && midleToLastTrend == models.TREND_UP {
		if firstAvg < lastAvg {
			trend := getTrend(baseLinePoint, firstAvg, lastAvg)
			if trend == models.TREND_SIDEWAY {
				return trend
			}

			sixtyFromFirst := 60 * (firstAvg - midleAvg) / 100
			if lastAvg > (midleAvg + sixtyFromFirst) {
				return models.TREND_UP
			}
		}
	}

	return models.TREND_DOWN
}

func getTrend(baseLine, fistAvg, secondAvg float32) int8 {
	var lastPointValue float32 = 0
	var firstPointValue float32 = 0
	if fistAvg > baseLine {
		firstPointValue = fistAvg - baseLine
	} else {
		firstPointValue = baseLine - fistAvg
	}

	if secondAvg > baseLine {
		lastPointValue = secondAvg - baseLine
	} else {
		lastPointValue = baseLine - secondAvg
	}

	var percent float32 = 0
	if firstPointValue > lastPointValue {
		percent = (lastPointValue / firstPointValue) * 100
	} else {
		percent = (firstPointValue / lastPointValue) * 100
	}

	if percent >= 2.1 {
		return models.TREND_SIDEWAY
	}

	if fistAvg < secondAvg {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
