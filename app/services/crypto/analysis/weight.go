package analysis

import (
	"telebot-trading/app/models"
)

var weightLogData map[string]float32
var longIntervalWeightLogData map[string]float32

func CalculateWeight(result *models.BandResult, masterCoin models.BandResult) float32 {
	weightLogData = map[string]float32{}

	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.Upper) {
		return 0
	}

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent < 2 {
		return 0
	}

	if masterCoin.Trend == models.TREND_DOWN {
		lastFourData := result.Bands[len(result.Bands)-4:]
		if CalculateTrends(lastFourData) != models.TREND_UP {
			return 0
		}
	}

	weight := priceChangeWeight(result.PriceChanges)
	if weight == 0 {
		return 0
	}
	weightLogData["priceWeight"] = weight

	isMasterCoinReversal := isMasterReversal(&masterCoin)
	positionWeight := getPositionWeight(result.Bands, result.Trend, masterCoin.Trend, false, isMasterCoinReversal, masterCoin.Direction)
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

	crossBandWeight := crossBandWeight(result)
	weightLogData["crossBandWeight"] = crossBandWeight
	weight += crossBandWeight

	if masterCoin.Trend == models.TREND_UP {
		weight += 0.1
		weightLogData["masterTrenWeight"] = 0.1
	}

	if result.Trend == models.TREND_UP {
		if isMasterCoinReversal {
			weight += 0.15
			weightLogData["TrenWeight"] = 0.15
		} else {
			weight += 0.1
			weightLogData["TrenWeight"] = 0.1
		}
	}

	return weight
}

func CalculateWeightLongInterval(result *models.BandResult, masterTrend int8) float32 {
	longIntervalWeightLogData = map[string]float32{}

	positionWeight := getPositionWeight(result.Bands, result.Trend, masterTrend, true, false, BAND_DOWN)
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

	weightCrossBand := crossBandWeight(result)
	weight += weightCrossBand
	longIntervalWeightLogData["weightCrossBand"] = weightCrossBand

	if result.Trend == models.TREND_UP {
		weight += 0.15
		longIntervalWeightLogData["trendWeight"] = 0.15
	}

	return weight
}

func CalculateWeightOnDown(result *models.BandResult) float32 {
	lastFour := result.Bands[len(result.Bands)-4:]

	crossLowerBand := false
	for _, data := range lastFour {
		if data.Candle.Low < float32(data.Lower) {
			crossLowerBand = true
			break
		}
	}

	if !crossLowerBand {
		return 0
	}

	marginUpper := getPriceMarginWithUpperBandWeight(result.Bands)

	return marginUpper + result.PriceChanges
}

func GetWeightLogData() map[string]float32 {
	return weightLogData
}

func GetLongIntervalWeightLogData() map[string]float32 {
	return longIntervalWeightLogData
}

func isMasterReversal(master *models.BandResult) bool {
	trend := CalculateTrends(master.Bands[:len(master.Bands)-1])

	lastFiveData := master.Bands[len(master.Bands)-4:]
	if trend == models.TREND_UP || CalculateTrends(lastFiveData[1:]) != models.TREND_UP || master.PriceChanges < 0.3 {
		return false
	}

	return true
}

func priceChangeWeight(priceChange float32) float32 {
	if priceChange >= 1.4 {
		return 0.5
	} else if priceChange >= 1.2 {
		return 0.45
	} else if priceChange >= 1 {
		return 0.4
	} else if priceChange >= 0.75 {
		return 0.35
	} else if priceChange >= 0.5 {
		return 0.3
	} else if priceChange >= 0.3 {
		return 0.2
	}

	return 0
}

func reversalWeight(result *models.BandResult) float32 {
	var weight float32 = 0
	trend := CalculateTrends(result.Bands[:len(result.Bands)-1])

	lastSixData := result.Bands[len(result.Bands)-6:]
	if trend == models.TREND_UP || CalculateTrends(lastSixData[1:]) != models.TREND_UP || ((result.PriceChanges < 1 && countSquentialUpBand(lastSixData) < 3) || result.PriceChanges < 1.3) {
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
	highUpNotInterested := CalculateTrends(lastSixData[:4]) == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.Upper)
	if countUpBand(lastSixData[1:]) < 2 && highUpNotInterested {
		return 0.08
	}

	firstBand := lastSixData[0]
	secondBand := lastSixData[1]
	thirdBand := lastSixData[2]
	fourthBand := lastSixData[3]
	isBandCrossWithLower := firstBand.Candle.Low <= float32(firstBand.Lower) || secondBand.Candle.Low <= float32(secondBand.Lower) || thirdBand.Candle.Low <= float32(thirdBand.Lower) || fourthBand.Candle.Low <= float32(fourthBand.Lower)
	isBandCrossWithSMA := firstBand.Candle.Low <= float32(firstBand.SMA) || secondBand.Candle.Low <= float32(secondBand.SMA) || thirdBand.Candle.Low <= float32(thirdBand.SMA) || fourthBand.Candle.Low <= float32(fourthBand.SMA)
	isBandCrossWithUpper := firstBand.Candle.Low <= float32(firstBand.Upper) || secondBand.Candle.Low <= float32(secondBand.Upper) || thirdBand.Candle.Low <= float32(thirdBand.Upper) || fourthBand.Candle.Low <= float32(fourthBand.Upper)
	if isBandCrossWithLower && float64(lastBand.Candle.Low) > lastBand.Lower {
		weight += 0.12
	} else if isBandCrossWithSMA && float64(lastBand.Candle.Low) > lastBand.SMA {
		weight += 0.1
	} else if isBandCrossWithUpper && float64(lastBand.Candle.Low) > lastBand.Upper {
		weight += 0.08
	}

	return weight
}

func crossBandWeight(result *models.BandResult) float32 {
	lastBand := result.Bands[len(result.Bands)-1]
	secondLastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Open > float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA) {
		if secondLastBand.Candle.Open < float32(secondLastBand.SMA) && (secondLastBand.Candle.Hight > float32(secondLastBand.SMA) || secondLastBand.Candle.Close > float32(secondLastBand.SMA)) {
			if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
				return 0.1
			}
		}
	}
	return 0
}

func getPatternWeight(result *models.BandResult) float32 {
	listMatchPattern := GetCandlePattern(result)

	var weight float32 = 0
	if len(listMatchPattern) > 0 {
		weight += 0.2 * float32(len(listMatchPattern))
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
		return 0.45
	} else if percent >= 2.5 {
		return 0.4
	} else if percent >= 2 {
		return 0.35
	} else if percent >= 1 {
		return 0.3
	}

	return 0.2
}

func getPositionWeight(bands []models.Band, trend, masterTrend int8, isLongInterval bool, isMasterCoinReversal bool, masterDirection int8) float32 {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]

	weightUpCounter := weightUpSquential(bands)

	lastFiveTrend := CalculateTrends(bands[len(bands)-5:])

	// low hight dibawah lower
	if lastBand.Candle.Hight < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.46
		}
	}

	// hight menyentuh lower tp close dibaawh lower
	if lastBand.Candle.Hight >= float32(lastBand.Lower) && lastBand.Candle.Close < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.42
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
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.44 + weightUpCounter
		}
	}

	// hight menyentuh SMA tp close dibaawh SMA
	if lastBand.Candle.Hight >= float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.36 + weightUpCounter
		}
	}

	// close menyentuh SMA tp open dibaawh SMA
	if lastBand.Candle.Close >= float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.46 + weightUpCounter
		}
	}

	// open menyentuh SMA tp low dibaawh SMA
	if lastBand.Candle.Open >= float32(lastBand.SMA) && lastBand.Candle.Low < float32(lastBand.SMA) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.48 + weightUpCounter
		}
	}

	// low hight dibawah Upper
	if lastBand.Candle.Hight < float32(lastBand.Upper) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.42 + weightUpCounter
		}
	}

	// hight menyentuh Upper tp close dibaawh Upper
	if lastBand.Candle.Hight >= float32(lastBand.Upper) && lastBand.Candle.Close < float32(lastBand.Upper) {
		return 0.34 + weightUpCounter
	}

	if ((masterTrend != models.TREND_DOWN && masterDirection == BAND_UP) || isMasterCoinReversal) && !isLongInterval {

		// close menyentuh Upper tp open dibaawh Upper
		if lastBand.Candle.Close >= float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) {
			return 0.32 + weightUpCounter
		}

		// close diatas upper dan band sebelumya juga diatas upper
		if lastBand.Candle.Close > float32(lastBand.Upper) {
			var val float32 = 0.3
			if secondLastBand.Candle.Close > float32(secondLastBand.Upper) {
				val += 0.17
			}

			return val
		}

	}

	return 0.28
}

func countUpBand(bands []models.Band) int {
	counter := 0
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			counter++
		}
	}

	return counter
}

func countSquentialUpBand(bands []models.Band) int {
	counter := 0
	for i := len(bands) - 1; i >= 0; i-- {
		if bands[1].Candle.Open < bands[i].Candle.Close {
			counter++
		} else {
			return counter
		}
	}

	return counter
}

func weightUpSquential(bands []models.Band) float32 {
	counter := countSquentialUpBand(bands)
	var weight float32 = 0.05
	if counter >= 5 {
		weight = 0.2
	} else if counter == 4 {
		weight = 0.175
	} else if counter == 3 {
		weight = 0.15
	} else if counter == 2 {
		weight = 0.1
	}

	return weight
}
