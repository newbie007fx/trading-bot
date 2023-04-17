package analysis

import (
	"log"
	"math"
	"telebot-trading/app/models"
)

const POSSIBLY_UP = 1
const POSSIBLY_DOWN = 2

func CalculateTrendsDetail(data []models.Band) models.TrendDetail {
	if len(data) < 4 {
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
	firstEndLimit := limit
	if limit < 1 {
		limit = 1
	}
	if firstEndLimit == 3 {
		firstEndLimit = 2
	}

	var totalFirstData float32 = 0
	var totalLastData float32 = 0
	var totalMidleData float32 = 0

	var midle_counter int = 0
	var first_counter int = 0
	var last_counter int = 0

	middleIndex := (len(data) / 2)
	lowestIndexFirst := 0
	lowestIndexMiddle := (len(data) / 2)
	for i, val := range data {
		if data[highestIndex].Candle.Close < val.Candle.Close {
			highestIndex = i
		}

		if data[lowestIndex].Candle.Close > val.Candle.Close {
			lowestIndex = i
		}

		if i < firstEndLimit {
			totalFirstData += val.Candle.Close
			first_counter++
		}

		if i >= middleIndex-(limit/2) && i <= middleIndex+(limit/2) {
			totalMidleData += val.Candle.Close
			midle_counter++
		}

		if i >= len(data)-firstEndLimit {
			totalLastData += val.Candle.Close
			last_counter++
		}

		if i <= middleIndex+(limit/2) {
			if data[lowestIndexFirst].Candle.Close > val.Candle.Close {
				lowestIndexFirst = i
			}
		}

		if i >= middleIndex-(limit/2) {
			if data[lowestIndexMiddle].Candle.Close > val.Candle.Close {
				lowestIndexMiddle = i
			}
		}
	}

	firstAvg := totalFirstData / float32(first_counter)
	lastAvg := totalLastData / float32(last_counter)
	midleAvg := totalMidleData / float32(midle_counter)
	baseLinePoint := data[lowestIndex].Candle.Close

	firstPercent, firstToMidleTrend := getTrend(data[lowestIndexFirst].Candle.Close, firstAvg, midleAvg)
	secondPercent, midleToLastTrend := getTrend(data[lowestIndexMiddle].Candle.Close, midleAvg, lastAvg)
	trend.FirstTrend = firstToMidleTrend
	trend.FirstTrendPercent = firstPercent
	trend.SecondTrend = midleToLastTrend
	trend.SecondTrendPercent = secondPercent
	trend.ShortTrend = conclusionShortTrend(CalculateTrendShort(data[len(data)-5:], false), CalculateTrendShort(data[len(data)-3:], true))

	if trend.Trend == 0 {
		trend.Trend = getConclusionTrend(firstToMidleTrend, midleToLastTrend, firstAvg, midleAvg, lastAvg, baseLinePoint)
	}

	return trend
}

func getConclusionTrend(firstToMidleTrend, midleToLastTrend int8, firstAvg, midleAvg, lastAvg, baseLinePointFirst float32) int8 {
	if firstToMidleTrend == models.TREND_SIDEWAY {
		if midleToLastTrend == models.TREND_SIDEWAY {
			_, trend := getTrend(baseLinePointFirst, firstAvg, lastAvg)
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
			_, trend := getTrend(baseLinePointFirst, firstAvg, lastAvg)
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

func getTrend(baseLine, fistAvg, secondAvg float32) (percent float32, trend int8) {
	firstPointValue := fistAvg - baseLine
	lastPointValue := secondAvg - baseLine

	if firstPointValue > lastPointValue {
		percent = (lastPointValue / firstPointValue) * 100
	} else {
		percent = (firstPointValue / lastPointValue) * 100
	}

	if percent >= 79 {
		trend = models.TREND_SIDEWAY
		return
	}

	if fistAvg < secondAvg {
		trend = models.TREND_UP
		return
	}

	trend = models.TREND_DOWN
	return
}

func CalculateShortTrendWithConclusion(data []models.Band) int8 {
	return conclusionShortTrend(CalculateTrendShort(data[len(data)-5:], false), CalculateTrendShort(data[len(data)-3:], true))
}

func CalculateTrendShort(data []models.Band, isIgnorePosiblyChecking bool) int8 {
	if len(data) == 0 {
		log.Println("invalid data when calculate trends")
		return models.TREND_DOWN
	}

	lowestIndex := 0
	thirtyPercent := float64(len(data)) * float64(15) / float64(100)
	limit := int(math.Floor(thirtyPercent))
	if limit < 1 {
		limit = 1
	}

	var totalFirstData float32 = 0
	var totalLastData float32 = 0
	possibly := checkPossibly(data[0], data[len(data)-1])

	for i, val := range data {
		if lowestFromBand(data[lowestIndex]) > lowestFromBand(val) {
			lowestIndex = i
		}

		if i < limit {
			totalFirstData += getBandValue(possibly, val, isIgnorePosiblyChecking)
		}

		if i >= len(data)-limit {
			totalLastData += val.Candle.Close
		}
	}

	firstAvg := totalFirstData / float32(limit)
	lastAvg := totalLastData / float32(limit)
	baseLinePoint := lowestFromBand(data[lowestIndex]) - (lowestFromBand(data[lowestIndex]) / 100)

	return getTrendShort(baseLinePoint, firstAvg, lastAvg)
}

func checkPossibly(firstBand models.Band, lastBand models.Band) int {
	if firstBand.Candle.Open > lastBand.Candle.Close && firstBand.Candle.Close > lastBand.Candle.Close {
		return POSSIBLY_DOWN
	} else if firstBand.Candle.Open < lastBand.Candle.Close && firstBand.Candle.Close < lastBand.Candle.Close {
		return POSSIBLY_UP
	}

	var hight, low float32 = 0, 0

	if firstBand.Candle.Open < firstBand.Candle.Close {
		return POSSIBLY_DOWN
	} else {
		low = lastBand.Candle.Close - firstBand.Candle.Close
		hight = firstBand.Candle.Open - lastBand.Candle.Close
	}

	if hight > low {
		return POSSIBLY_DOWN
	}

	return POSSIBLY_UP
}

func getBandValue(possibly int, band models.Band, isIgnorePosiblyChecking bool) float32 {
	if possibly == POSSIBLY_DOWN || isIgnorePosiblyChecking {
		if band.Candle.Close > band.Candle.Open {
			return band.Candle.Close
		}

		return band.Candle.Open
	}

	if band.Candle.Close < band.Candle.Open {
		return band.Candle.Close
	}

	return band.Candle.Open
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

	if percent >= 85 {
		return models.TREND_SIDEWAY
	}

	if fistAvg < secondAvg {
		return models.TREND_UP
	}

	return models.TREND_DOWN
}

func conclusionShortTrend(firstTrend, secondTrend int8) int8 {
	if firstTrend == secondTrend {
		return firstTrend
	}

	if (firstTrend == models.TREND_DOWN && secondTrend == models.TREND_UP) || (firstTrend == models.TREND_UP && secondTrend == models.TREND_DOWN) {
		return models.TREND_SIDEWAY
	}

	if firstTrend == models.TREND_SIDEWAY && secondTrend != models.TREND_UP {
		return secondTrend
	}

	if secondTrend == models.TREND_SIDEWAY && firstTrend != models.TREND_UP {
		return firstTrend
	}

	return models.TREND_SIDEWAY
}

func higestFromBand(band models.Band) float32 {
	if band.Candle.Close > band.Candle.Open {
		return band.Candle.Close
	}

	return band.Candle.Open
}
