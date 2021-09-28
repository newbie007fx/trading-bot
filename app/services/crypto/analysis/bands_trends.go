package analysis

import (
	"log"
	"math"
	"telebot-trading/app/models"
)

func CalculateTrends(data []models.Band) int8 {
	if len(data) == 0 {
		log.Println("invalid data when calculate trends")
		return models.TREND_DOWN
	}

	highestIndex, lowestIndex := 0, 0
	thirtyPercent := float64(len(data)) * float64(15) / float64(100)
	limit := int(math.Floor(thirtyPercent))
	if limit < 1 {
		limit = 1
	}

	var totalFirstData float32 = 0
	var totalLastData float32 = 0
	var totalMidleData float32 = 0

	var midle_counter int = 0

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
		}

		if i > middleIndex-(limit/2) && i <= middleIndex+(limit/2) {
			totalMidleData += val.Candle.Close
			midle_counter++
		}

		if i >= len(data)-limit {
			totalLastData += val.Candle.Close
		}
	}

	firstAvg := totalFirstData / float32(limit)
	lastAvg := totalLastData / float32(limit)
	midleAvg := totalMidleData / float32(midle_counter)
	baseLinePoint := data[lowestIndex].Candle.Close

	firstToMidleTrend := getTrend(baseLinePoint, firstAvg, midleAvg)
	midleToLastTrend := getTrend(baseLinePoint, midleAvg, lastAvg)

	return getConclusionTrend(firstToMidleTrend, midleToLastTrend, firstAvg, midleAvg, lastAvg, baseLinePoint)
}

func CalculateTrendsDetail(data []models.Band) models.TrendDetail {
	if len(data) == 0 {
		log.Println("invalid data when calculate trends")
		return models.TrendDetail{
			FirstTrend:  models.TREND_DOWN,
			SecondTrend: models.TREND_DOWN,
			Trend:       models.TREND_DOWN,
			ShortTrend:  models.TREND_DOWN,
		}
	}

	trend := models.TrendDetail{}
	highestIndex, lowestIndex := 0, 0
	thirtyPercent := float64(len(data)) * float64(19) / float64(100)
	limit := int(math.Floor(thirtyPercent))
	if limit < 1 {
		limit = 1
	}

	var totalFirstData float32 = 0
	var totalLastData float32 = 0
	var totalMidleData float32 = 0

	var midle_counter int = 0

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
		}

		if i > middleIndex-(limit/2) && i <= middleIndex+(limit/2) {
			totalMidleData += val.Candle.Close
			midle_counter++
		}

		if i >= len(data)-limit {
			totalLastData += val.Candle.Close
		}
	}

	firstAvg := totalFirstData / float32(limit)
	lastAvg := totalLastData / float32(limit)
	midleAvg := totalMidleData / float32(midle_counter)
	baseLinePoint := data[lowestIndex].Candle.Close

	firstToMidleTrend := getTrend(baseLinePoint, firstAvg, midleAvg)
	midleToLastTrend := getTrend(baseLinePoint, midleAvg, lastAvg)
	trend.FirstTrend = firstToMidleTrend
	trend.SecondTrend = midleToLastTrend
	trend.ShortTrend = CalculateTrendShort(data[len(data)-4:])

	if trend.Trend == 0 {
		trend.Trend = getConclusionTrend(firstToMidleTrend, midleToLastTrend, firstAvg, midleAvg, lastAvg, baseLinePoint)
	}

	return trend
}

func getConclusionTrend(firstToMidleTrend, midleToLastTrend int8, firstAvg, midleAvg, lastAvg, baseLinePointFirst float32) int8 {
	if firstToMidleTrend == models.TREND_SIDEWAY {
		if midleToLastTrend == models.TREND_SIDEWAY {
			trend := getTrend(baseLinePointFirst, firstAvg, lastAvg)
			return trend
		}
		return midleToLastTrend
	}

	if midleToLastTrend == models.TREND_SIDEWAY {
		return models.TREND_SIDEWAY
	}

	if firstToMidleTrend == models.TREND_UP && midleToLastTrend == models.TREND_UP {
		return models.TREND_UP
	}

	if firstToMidleTrend == models.TREND_DOWN && midleToLastTrend == models.TREND_DOWN {
		return models.TREND_DOWN
	}

	if firstToMidleTrend == models.TREND_UP && midleToLastTrend == models.TREND_DOWN {
		return models.TREND_DOWN
	}

	if firstToMidleTrend == models.TREND_DOWN && midleToLastTrend == models.TREND_UP {
		if firstAvg < lastAvg {
			trend := getTrend(baseLinePointFirst, firstAvg, lastAvg)
			if trend == models.TREND_SIDEWAY {
				return trend
			}
			return models.TREND_UP
		}

		sixtyFromFirst := 65 * (firstAvg - midleAvg) / 100
		if lastAvg > (midleAvg + sixtyFromFirst) {
			return models.TREND_UP
		}
	}

	return models.TREND_DOWN
}

func getTrend(baseLine, fistAvg, secondAvg float32) int8 {
	firstPointValue := fistAvg - baseLine
	lastPointValue := secondAvg - baseLine

	var percent float32 = 0
	if firstPointValue > lastPointValue {
		percent = (lastPointValue / firstPointValue) * 100
	} else {
		percent = (firstPointValue / lastPointValue) * 100
	}

	if percent >= 79 {
		return models.TREND_SIDEWAY
	}

	if fistAvg < secondAvg {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}

func CalculateTrendShort(data []models.Band) int8 {
	if len(data) == 0 {
		log.Println("invalid data when calculate trends")
		return models.TREND_DOWN
	}

	highestIndex, lowestIndex := 0, 0
	thirtyPercent := float64(len(data)) * float64(15) / float64(100)
	limit := int(math.Floor(thirtyPercent))
	if limit < 1 {
		limit = 1
	}

	var totalFirstData float32 = 0
	var totalLastData float32 = 0

	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}

		if i < limit {
			totalFirstData += val.Candle.Close
		}

		if i >= len(data)-limit {
			totalLastData += val.Candle.Close
		}
	}

	firstAvg := totalFirstData / float32(limit)
	lastAvg := totalLastData / float32(limit)
	baseLinePoint := data[lowestIndex].Candle.Close

	return getTrendShort(baseLinePoint, firstAvg, lastAvg)
}

func getTrendShort(baseLine, fistAvg, secondAvg float32) int8 {
	firstPointValue := fistAvg - baseLine
	lastPointValue := secondAvg - baseLine

	var percent float32 = 0
	if firstPointValue > lastPointValue {
		percent = (lastPointValue / firstPointValue) * 100
	} else {
		percent = (firstPointValue / lastPointValue) * 100
	}

	if percent >= 81 {
		return models.TREND_SIDEWAY
	}

	if fistAvg < secondAvg {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}
