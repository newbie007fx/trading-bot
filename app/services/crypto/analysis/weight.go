package analysis

import (
	"telebot-trading/app/models"
)

var weightLogData map[string]float32
var longIntervalWeightLogData map[string]float32

func CalculateWeight(result *models.BandResult) float32 {
	weightLogData = map[string]float32{}

	weight := priceChangeWeight(result.PriceChanges)
	weightLogData["priceWeight"] = weight

	positionWeight := getPositionWeight(result.Bands, false)
	weightLogData["positionWeight"] = positionWeight
	weight += positionWeight

	priceMarginWeight := getPriceMarginWithUpperBandWeight(result.Bands)
	weightLogData["priceMarginWeight"] = priceMarginWeight
	weight += priceMarginWeight

	patternWeight := getPatternWeight(result)
	weightLogData["patternWeight"] = patternWeight
	weight += patternWeight

	weightReversal := reversalWeight(result)
	weightLogData["reversalWeight"] = weightReversal
	weight += weightReversal

	return weight
}

func CalculateWeightLongInterval(result *models.BandResult) float32 {
	longIntervalWeightLogData = map[string]float32{}

	positionWeight := getPositionWeight(result.Bands, true)
	weight := positionWeight
	longIntervalWeightLogData["PositionWeight"] = positionWeight

	priceMarginWeight := getPriceMarginWithUpperBandWeight(result.Bands)
	weight += priceMarginWeight
	longIntervalWeightLogData["PriceMargin"] = priceMarginWeight

	patternWeight := getPatternWeight(result)
	weight += patternWeight
	longIntervalWeightLogData["PatternWeight"] = patternWeight

	weightReseversal := reversalWeight(result)
	weight += weightReseversal
	longIntervalWeightLogData["weightReversal"] = weightReseversal

	return weight
}

func GetWeightLogData() map[string]float32 {
	return weightLogData
}

func GetLongIntervalWeightLogData() map[string]float32 {
	return longIntervalWeightLogData
}

func priceChangeWeight(priceChange float32) float32 {
	if priceChange >= 1.4 {
		return 0.5
	} else if priceChange >= 1.2 {
		return 0.46
	} else if priceChange >= 1 {
		return 0.42
	} else if priceChange >= 0.75 {
		return 0.38
	} else if priceChange >= 0.5 {
		return 0.34
	} else if priceChange >= 0.3 {
		return 0.28
	}

	return 0.16
}

func reversalWeight(result *models.BandResult) float32 {
	var weight float32 = 0
	trend := CalculateTrendsDetail(result.Bands[:len(result.Bands)-1])

	lastSixData := result.Bands[len(result.Bands)-6:]
	if trend.Trend == models.TREND_UP || CalculateTrendShort(lastSixData[2:]) != models.TREND_UP || ((result.PriceChanges < 0.8 && CountSquentialUpBand(lastSixData) < 3) || result.PriceChanges < 1.1) {
		lastSixDataTrend := CalculateTrendsDetail(lastSixData)
		if lastSixDataTrend.FirstTrend == models.TREND_DOWN && lastSixDataTrend.SecondTrend == models.TREND_UP {
			weight = 0.08
		} else {
			weight = 0.05
		}
	} else {
		weight = 0.15
	}

	lastBand := lastSixData[5]
	highUpNotInterested := CalculateTrendShort(lastSixData[2:]) == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.Upper)
	if CountUpBand(lastSixData[1:]) < 2 && highUpNotInterested {
		return 0.08
	}

	firstBand := lastSixData[0]
	secondBand := lastSixData[1]
	thirdBand := lastSixData[2]
	fourthBand := lastSixData[3]
	isBandCrossWithLower := firstBand.Candle.Low <= float32(firstBand.Lower) || secondBand.Candle.Low <= float32(secondBand.Lower) || thirdBand.Candle.Low <= float32(thirdBand.Lower) || fourthBand.Candle.Low <= float32(fourthBand.Lower)
	isBandCrossWithSMA := firstBand.Candle.Low <= float32(firstBand.SMA) || secondBand.Candle.Low <= float32(secondBand.SMA) || thirdBand.Candle.Low <= float32(thirdBand.SMA) || fourthBand.Candle.Low <= float32(fourthBand.SMA)
	isBandCrossWithUpper := firstBand.Candle.Low <= float32(firstBand.Upper) || secondBand.Candle.Low <= float32(secondBand.Upper) || thirdBand.Candle.Low <= float32(thirdBand.Upper) || fourthBand.Candle.Low <= float32(fourthBand.Upper)
	if isBandCrossWithLower && float64(lastBand.Candle.Open) > lastBand.Lower {
		weight += 0.12
		if result.AllTrend.FirstTrend != models.TREND_DOWN || result.AllTrend.SecondTrend != models.TREND_DOWN {
			weight += 0.3
		}
	} else if isBandCrossWithSMA {
		weight += 0.05
		if float64(lastBand.Candle.Open) > lastBand.SMA {
			weight += 0.05
		}
	} else if isBandCrossWithUpper && float64(lastBand.Candle.Open) > lastBand.Upper {
		weight += 0.08
	}

	return weight
}

func getPatternWeight(result *models.BandResult) float32 {
	listMatchPattern := GetCandlePattern(result)

	var weight float32 = 0
	if len(listMatchPattern) > 0 {
		weight = 0.2
	}

	return weight
}

func getPriceMarginWithUpperBandWeight(bands []models.Band) float32 {
	lastBand := bands[len(bands)-1]
	var percent float32 = 0

	if lastBand.Candle.Close < float32(lastBand.SMA) {
		different := float32(lastBand.SMA) - lastBand.Candle.Close
		percent = different / lastBand.Candle.Close * 100
	} else if lastBand.Candle.Close < float32(lastBand.Upper) {
		different := float32(lastBand.Upper) - lastBand.Candle.Close
		percent = different / lastBand.Candle.Close * 100
	}

	return getPriceMarginWithUpperBandPercentWeight(percent)
}

func getPriceMarginWithUpperBandPercentWeight(percent float32) float32 {
	if percent >= 3.5 {
		return 0.5
	} else if percent >= 3 {
		return 0.46
	} else if percent >= 2.5 {
		return 0.42
	} else if percent >= 2 {
		return 0.38
	} else if percent >= 1 {
		return 0.34
	}

	return 0.21
}

func getPositionWeight(bands []models.Band, isLongInterval bool) float32 {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]

	weightUp := weightUpSquential(bands)

	// low hight dibawah lower
	if lastBand.Candle.Hight < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.46
		}
	}

	// hight menyentuh lower tp close dibaawh lower
	if lastBand.Candle.Hight >= float32(lastBand.Lower) && lastBand.Candle.Close < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.44
		}
	}

	// close menyentuh lower tp open dibaawh lower
	if lastBand.Candle.Close >= float32(lastBand.Lower) && lastBand.Candle.Open < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.48
		}
	}

	// open menyentuh lower tp low dibaawh lower
	if lastBand.Candle.Open >= float32(lastBand.Lower) && lastBand.Candle.Low < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.5
		}
	}

	// low hight dibawah SMA
	if lastBand.Candle.Hight < float32(lastBand.SMA) {
		return 0.40 + weightUp
	}

	// hight menyentuh SMA tp close dibaawh SMA
	if lastBand.Candle.Hight >= float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		return 0.38 + weightUp
	}

	// close menyentuh SMA tp open dibaawh SMA
	if lastBand.Candle.Close >= float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) {
		return 0.42 + weightUp
	}

	// open menyentuh SMA tp low dibaawh SMA
	if lastBand.Candle.Open >= float32(lastBand.SMA) && lastBand.Candle.Low < float32(lastBand.SMA) {
		return 0.38 + weightUp
	}

	// low hight dibawah Upper
	if lastBand.Candle.Hight < float32(lastBand.Upper) {
		return 0.36 + weightUp
	}

	// hight menyentuh Upper tp close dibaawh Upper
	if lastBand.Candle.Hight >= float32(lastBand.Upper) && lastBand.Candle.Close < float32(lastBand.Upper) {
		return 0.32 + weightUp
	}

	if !isLongInterval {

		// close menyentuh Upper tp open dibaawh Upper
		if lastBand.Candle.Close >= float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) {
			return 0.36
		}

		// close diatas upper dan band sebelumya juga diatas upper
		if lastBand.Candle.Close > float32(lastBand.Upper) {
			var val float32 = 0.26
			if secondLastBand.Candle.Close > float32(secondLastBand.Upper) {
				val += 0.12
			}

			return val
		}

	}

	return 0.26
}

func CountUpBand(bands []models.Band) int {
	counter := 0
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			difference := band.Candle.Close - band.Candle.Open
			if (difference / band.Candle.Open * 100) > 0.1 {
				counter++
			}
		}
	}

	return counter
}

func CountSquentialUpBand(bands []models.Band) int {
	counter := 0
	for i := len(bands) - 1; i >= 0; i-- {
		if bands[1].Candle.Open < bands[i].Candle.Close {
			difference := bands[1].Candle.Close - bands[1].Candle.Open
			if (difference / bands[1].Candle.Open * 100) > 0.1 {
				counter++
			} else {
				return counter
			}
		} else {
			return counter
		}
	}

	return counter
}

func weightUpSquential(bands []models.Band) float32 {
	counter := CountSquentialUpBand(bands)
	var weight float32 = 0.05
	if counter >= 5 {
		weight = 0.111
	} else if counter == 4 {
		weight = 0.125
	} else if counter == 3 {
		weight = 0.1
	} else if counter == 2 {
		weight = 0.05
	}

	return weight
}
