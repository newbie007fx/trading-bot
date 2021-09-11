package analysis

import (
	"telebot-trading/app/models"
)

var weightLogData map[string]float32
var longIntervalWeightLogData map[string]float32

func CalculateWeight(result *models.BandResult, masterCoin models.BandResult) float32 {
	weightLogData = map[string]float32{}

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent < 2.3 {
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

	if result.VolumeChanges > 0 {
		volumeWight := getVolumeAverageChangesWeight(result.VolumeChanges)
		weightLogData["volumeWeight"] = volumeWight
		weight += volumeWight
	}

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
			weight += 0.2
			weightLogData["TrenWeight"] = 0.2
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
		weight += 0.2
		longIntervalWeightLogData["trendWeight"] = 0.2
	}

	return weight
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
	if trend == models.TREND_UP || CalculateTrends(lastFiveData[1:]) != models.TREND_UP || master.PriceChanges < 0.5 {
		return false
	}

	return true
}

func priceChangeWeight(priceChange float32) float32 {
	if priceChange >= 1.4 {
		return 0.5
	} else if priceChange >= 1.2 {
		return 0.4
	} else if priceChange >= 1 {
		return 0.3
	} else if priceChange >= 0.75 {
		return 0.2
	} else if priceChange >= 0.5 {
		return 0.1
	}

	return 0
}

func reversalWeight(result *models.BandResult) float32 {
	var weight float32 = 0
	trend := CalculateTrends(result.Bands[:len(result.Bands)-1])

	lastFiveData := result.Bands[len(result.Bands)-5:]
	if trend == models.TREND_UP || CalculateTrends(lastFiveData[1:]) != models.TREND_UP || result.PriceChanges < 1.25 {
		weight = 0
	} else {
		weight = 0.15
	}

	lastBand := lastFiveData[4]
	if trend == models.TREND_SIDEWAY {
		if result.Position == models.ABOVE_UPPER && lastBand.Candle.Open >= float32(lastBand.Upper) {
			return 0.1
		}

		if result.Position == models.ABOVE_SMA && lastBand.Candle.Open >= float32(lastBand.SMA) {
			return 0.101
		}

		if result.Position == models.BELOW_SMA && lastBand.Candle.Open >= float32(lastBand.Lower) {
			return 0.102
		}
	}

	highUpNotInterested := CalculateTrends(lastFiveData[:4]) != models.TREND_UP && lastBand.Candle.Close > float32(lastBand.Upper)
	if countUpBand(lastFiveData[1:]) < 2 && highUpNotInterested {
		return 0.102
	}

	firstBand := lastFiveData[0]
	secondBand := lastFiveData[1]
	thirdBand := lastFiveData[2]
	isBandCrossWithLower := firstBand.Candle.Low <= float32(firstBand.Lower) || secondBand.Candle.Low <= float32(secondBand.Lower) || thirdBand.Candle.Low <= float32(thirdBand.Lower)
	isBandCrossWithSMA := firstBand.Candle.Low <= float32(firstBand.SMA) || secondBand.Candle.Low <= float32(secondBand.SMA) || thirdBand.Candle.Low <= float32(thirdBand.SMA)
	isBandCrossWithUpper := firstBand.Candle.Low <= float32(firstBand.Upper) || secondBand.Candle.Low <= float32(secondBand.Upper) || thirdBand.Candle.Low <= float32(thirdBand.Upper)
	if isBandCrossWithLower && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].Lower {
		weight += 0.25
	} else if isBandCrossWithSMA && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].SMA {
		weight += 0.2
	} else if isBandCrossWithUpper && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].Upper {
		weight += 0.15
	}

	return weight
}

func crossBandWeight(result *models.BandResult) float32 {
	lastFourData := result.Bands[len(result.Bands)-4 : len(result.Bands)]
	if CalculateTrends(lastFourData) != models.TREND_UP {
		return 0
	}

	lastBand := lastFourData[3]
	secondLastBand := lastFourData[2]
	if lastBand.Candle.Close > float32(lastBand.Upper) && secondLastBand.Candle.Open < float32(lastBand.Upper) {
		return 0.2
	} else if lastBand.Candle.Close > float32(lastBand.SMA) && secondLastBand.Candle.Open < float32(lastBand.SMA) {
		return 0.25
	} else if lastBand.Candle.Close > float32(lastBand.Lower) && secondLastBand.Candle.Open < float32(lastBand.Lower) {
		return 0.3
	}

	return 0.15
}

func getPatternWeight(result *models.BandResult) float32 {
	listMatchPattern := GetCandlePattern(result)

	var weight float32 = 0
	if len(listMatchPattern) > 0 {
		weight += 0.35 * float32(len(listMatchPattern))
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
		return 0.375
	} else if percent >= 2 {
		return 0.3
	} else if percent >= 1 {
		return 0.225
	}

	return 0.2
}

func getVolumeAverageChangesWeight(volumeAverageChanges float32) float32 {
	if volumeAverageChanges >= 101 {
		return 0.35
	} else if volumeAverageChanges >= 81 {
		return 0.3
	} else if volumeAverageChanges >= 61 {
		return 0.25
	} else if volumeAverageChanges >= 41 {
		return 0.2
	} else if volumeAverageChanges >= 21 {
		return 0.15
	} else if volumeAverageChanges >= 1 {
		return 0.1
	}

	return 0
}

func getPositionWeight(bands []models.Band, trend, masterTrend int8, isLongInterval bool, isMasterCoinReversal bool, masterDirection int8) float32 {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]

	lastFiveTrend := CalculateTrends(bands[len(bands)-5:])

	// low hight dibawah lower
	if lastBand.Candle.Hight < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.35
		}
	}

	// hight menyentuh lower tp close dibaawh lower
	if lastBand.Candle.Hight >= float32(lastBand.Lower) && lastBand.Candle.Close < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.2
		}
	}

	// close menyentuh lower tp open dibaawh lower
	if lastBand.Candle.Close >= float32(lastBand.Lower) && lastBand.Candle.Open < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.425
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
			return 0.325
		}
	}

	// hight menyentuh SMA tp close dibaawh SMA
	if lastBand.Candle.Hight >= float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.225
		}
	}

	// close menyentuh SMA tp open dibaawh SMA
	if lastBand.Candle.Close >= float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.4
		}
	}

	// open menyentuh SMA tp low dibaawh SMA
	if lastBand.Candle.Open >= float32(lastBand.SMA) && lastBand.Candle.Low < float32(lastBand.SMA) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.475
		}
	}

	// low hight dibawah Upper
	if lastBand.Candle.Hight < float32(lastBand.Upper) {
		if lastFiveTrend == models.TREND_UP || trend != models.TREND_DOWN {
			return 0.275
		}
	}

	if ((masterTrend != models.TREND_DOWN && masterDirection == BAND_UP) || isMasterCoinReversal) && !isLongInterval {

		// hight menyentuh Upper tp close dibaawh Upper
		if lastBand.Candle.Hight >= float32(lastBand.Upper) && lastBand.Candle.Close < float32(lastBand.Upper) {
			return 0.175
		}

		// close menyentuh Upper tp open dibaawh Upper
		if lastBand.Candle.Close >= float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) {
			return 0.35
		}

		// close diatas upper dan band sebelumya juga diatas upper
		if lastBand.Candle.Close > float32(lastBand.Upper) {
			var val float32 = 0.15
			if secondLastBand.Candle.Close > float32(secondLastBand.Upper) && secondLastBand.Candle.Close < float32(secondLastBand.Upper) {
				val += 0.17
			}

			return val
		}

		// open menyentuh Upper tp low dibaawh Upper
		if lastBand.Candle.Open >= float32(lastBand.Upper) && lastBand.Candle.Low < float32(lastBand.Upper) {
			return 0.425
		}
	}

	return 0.15
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
