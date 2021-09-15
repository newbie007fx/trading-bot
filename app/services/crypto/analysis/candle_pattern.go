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

func SellPattern(bandResult *models.BandResult) bool {
	if bandResult.AllTrend.SecondTrend == models.TREND_UP && bandResult.Direction == BAND_DOWN {
		if bearishEngulfing(bandResult) {
			return true
		}
	}
	return false
}

func bearishEngulfing(result *models.BandResult) bool {
	if result.Position == models.ABOVE_SMA || result.Position == models.ABOVE_UPPER {
		firstBand := result.Bands[len(result.Bands)-1]
		secondBand := result.Bands[len(result.Bands)-2]
		if secondBand.Candle.Open < secondBand.Candle.Close {
			firstDifferent := firstBand.Candle.Open - firstBand.Candle.Close
			secondDifferent := secondBand.Candle.Close - secondBand.Candle.Open
			return firstDifferent > secondDifferent
		}
	}

	return false
}

func hammer(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	if secondLastBand.Candle.Open > secondLastBand.Candle.Close {

		if lastBand.Candle.Low <= float32(lastBand.Lower) || secondLastBand.Candle.Low <= float32(secondLastBand.Lower) {
			different := lastBand.Candle.Hight - lastBand.Candle.Close
			candleBody := lastBand.Candle.Hight - lastBand.Candle.Low
			percent := different / candleBody * 100
			if percent < 5 {
				different = lastBand.Candle.Open - lastBand.Candle.Low
				percent := different / candleBody * 100
				return percent >= 60
			}
		}
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
			different := lastBand.Candle.Close - lastBand.Candle.Open
			candleBody := lastBand.Candle.Hight - lastBand.Candle.Low
			percent := different / candleBody * 100
			if percent < 20 {
				different = lastBand.Candle.Open - lastBand.Candle.Low
				percent := different / candleBody * 100
				return percent >= 60
			}
		}
	}

	return false
}

func turnPattern(bands []models.Band) bool {
	numberOfData := len(bands) / 4
	for i := numberOfData; i >= 0; i-- {
		currentValue := (bands[i].Candle.Open + bands[i].Candle.Close) / 2
		if !(currentValue > float32(bands[i].SMA)) {
			return false
		}
	}

	return true
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
