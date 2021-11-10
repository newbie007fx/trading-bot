package analysis

import (
	"fmt"
	"telebot-trading/app/models"
)

const PATTERN_HAMMER int8 = 1
const PATTERN_INVERTED_HAMMER int8 = 2
const PATTERN_BULLISH_HARAMI int8 = 3
const PATTERN_DRAGONFLY_DOJI int8 = 4
const PATTERN_THREE_WHITE_SOLDIERS int8 = 5
const PATTERN_TURN int8 = 6

func GetCandlePattern(bandResult *models.BandResult) []int8 {
	bands := bandResult.Bands
	result := []int8{}

	if bandResult.Trend == models.TREND_DOWN {
		if hammer(bands) {
			result = append(result, PATTERN_HAMMER)
		}

		if invertedHammer(bands) {
			result = append(result, PATTERN_INVERTED_HAMMER)
		}

		if bullishHarami(bands) {
			result = append(result, PATTERN_BULLISH_HARAMI)
		}

		if dragonflyDoji(bands) {
			result = append(result, PATTERN_DRAGONFLY_DOJI)
		}

		if threeWhiteSoldiers(bands) {
			result = append(result, PATTERN_THREE_WHITE_SOLDIERS)
		}
	}

	if bandResult.AllTrend.SecondTrend == models.TREND_UP && turnPattern(bands) {
		result = append(result, PATTERN_TURN)
	}

	if len(result) > 0 {
		fmt.Println(result)
	}

	return result
}

func BearishEngulfing(bands []models.Band) bool {
	isContainBearishEngulfing := false
	for i := 0; i < len(bands)-1; i++ {
		firstBand := bands[i]
		secondBand := bands[i+1]
		if firstBand.Candle.Open < firstBand.Candle.Close && secondBand.Candle.Open > secondBand.Candle.Close {
			firstDifferent := firstBand.Candle.Close - firstBand.Candle.Open
			if firstDifferent/firstBand.Candle.Open*100 > 0.1 {
				secondDifferent := secondBand.Candle.Open - secondBand.Candle.Close
				if secondDifferent > firstDifferent {
					percent := firstDifferent / secondDifferent * 100
					isContainBearishEngulfing = percent > 40
				}
			}
		}
		if isContainBearishEngulfing {
			return true
		}
	}

	return isContainBearishEngulfing
}

func hammer(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	if secondLastBand.Candle.Open > secondLastBand.Candle.Close {
		if lastBand.Candle.Low <= float32(lastBand.Lower) || secondLastBand.Candle.Low <= float32(secondLastBand.Lower) {
			return IsHammer(bands)
		}
	}

	return false
}

func IsHammer(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	different := lastBand.Candle.Hight - lastBand.Candle.Close
	candleBody := lastBand.Candle.Hight - lastBand.Candle.Low
	percent := different / candleBody * 100
	if percent < 10 {
		different = lastBand.Candle.Open - lastBand.Candle.Low
		percent := different / candleBody * 100
		return percent >= 65
	}

	return false
}

func invertedHammer(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	if secondLastBand.Candle.Open > secondLastBand.Candle.Close {

		if lastBand.Candle.Low <= float32(lastBand.Lower) || secondLastBand.Candle.Low <= float32(secondLastBand.Lower) {
			different := lastBand.Candle.Open - lastBand.Candle.Low
			candleBody := lastBand.Candle.Hight - lastBand.Candle.Low
			percent := different / candleBody * 100
			if percent < 5 {
				different = lastBand.Candle.Close - lastBand.Candle.Low
				percent := different / candleBody * 100
				return percent <= 40
			}
		}
	}

	return false
}

func bullishHarami(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	if secondLastBand.Candle.Open > secondLastBand.Candle.Close {
		differenceSecondLast := secondLastBand.Candle.Open - secondLastBand.Candle.Close
		differenceLast := lastBand.Candle.Close - lastBand.Candle.Open
		if differenceLast > differenceSecondLast {
			return false
		}
		return lastBand.Candle.Low > secondLastBand.Candle.Close || secondLastBand.Candle.Hight < secondLastBand.Candle.Open
	}

	return false
}

func dragonflyDoji(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	if secondLastBand.Candle.Open > secondLastBand.Candle.Close {

		if lastBand.Candle.Low <= float32(lastBand.Lower) || secondLastBand.Candle.Low <= float32(secondLastBand.Lower) {
			return IsDoji(lastBand, true)
		}
	}

	return false
}

func IsDoji(band models.Band, isUp bool) bool {
	if isUp && band.Candle.Close > band.Candle.Open {
		different := band.Candle.Close - band.Candle.Open
		candleBody := band.Candle.Hight - band.Candle.Low
		percent := different / candleBody * 100
		if percent < 15 {
			different = band.Candle.Open - band.Candle.Low
			percent := different / candleBody * 100
			return percent >= 60
		}
	}

	if !isUp && band.Candle.Close < band.Candle.Open {
		different := band.Candle.Open - band.Candle.Close
		candleBody := band.Candle.Hight - band.Candle.Low
		percent := different / candleBody * 100
		if percent < 15 {
			different = band.Candle.Close - band.Candle.Low
			percent := different / candleBody * 100
			return percent >= 60
		}
	}

	return false
}

func turnPattern(bands []models.Band) bool {
	numberOfData := len(bands) / 3
	countDown := 0
	for i := len(bands) - 1; i >= len(bands)-numberOfData; i-- {
		currentValue := (bands[i].Candle.Open + bands[i].Candle.Close) / 2
		if currentValue < float32(bands[i].SMA) {
			return false
		}
		if bands[i].Candle.Open > bands[i].Candle.Close {
			countDown++
		}
	}

	return countDown < numberOfData/2
}

func threeWhiteSoldiers(bands []models.Band) bool {
	lastFive := bands[len(bands)-5:]

	if hasAnyBandCrossWithLower(lastFive) {
		threeBand := lastFive[1:4]
		return checkWhiteSoldiers(threeBand)
	}

	return false
}

func hasAnyBandCrossWithLower(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Low <= float32(band.Lower) {
			return true
		}
	}

	return false
}

func checkWhiteSoldiers(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open > band.Candle.Close {
			return false
		}

		differentUp := band.Candle.Hight - band.Candle.Close
		differentDown := band.Candle.Open - band.Candle.Low
		candleBody := band.Candle.Hight - band.Candle.Low
		percentUp := differentUp / candleBody * 100
		percentDown := differentDown / candleBody * 100
		if percentDown > 15 || percentUp > 15 {
			return false
		}
	}

	return true
}

func ThreeBlackCrowds(bands []models.Band) bool {
	var bandDown []models.Band = []models.Band{}
	for i, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			if len(bandDown)+len(bands)-i-1 < 3 {
				break
			}
			continue
		}

		bandDown = append(bandDown, band)
	}

	if len(bandDown) >= 3 {
		return bandDown[0].Candle.Hight > float32(bandDown[0].Upper)
	}

	return false
}
