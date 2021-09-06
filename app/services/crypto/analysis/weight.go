package analysis

import (
	"fmt"
	"telebot-trading/app/models"
)

var weightLog string = ""

func CalculateWeight(result *models.BandResult, masterTrend int8) float32 {
	weightLog = ""

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent < 2.5 {
		return 0
	}

	if masterTrend == models.TREND_DOWN {
		lastFourData := result.Bands[len(result.Bands)-4 : len(result.Bands)]
		if CalculateTrends(lastFourData) != models.TREND_UP {
			return 0
		}
	}

	weight := priceChangeWeight(result.PriceChanges)
	weightLog += fmt.Sprintf("priceWeight: %.2f", weight)
	if weight == 0 {
		return 0
	}

	if result.VolumeChanges > 0 {
		volumeWight := getVolumeAverageChangesWeight(result.VolumeChanges)
		weightLog += fmt.Sprintf(", volumeWeight: %.2f", volumeWight)
		weight += volumeWight
	}

	positionWeight := getPositionWeight(result.Bands, result.Trend, masterTrend, false)
	weightLog += fmt.Sprintf(", positionWeight: %.2f", positionWeight)
	weight += positionWeight

	priceMarginWeight := getPriceMarginWithUpperBandWeight(result.Bands)
	weightLog += fmt.Sprintf(", priceMarginWeight: %.2f", priceMarginWeight)
	weight += priceMarginWeight

	patternWeight := getPatternWeight(result)
	weightLog += fmt.Sprintf(", patternWeight: %.2f", patternWeight)
	weight += patternWeight

	reversalWeight := reversalWeight(result)
	weightLog += fmt.Sprintf(", reversalWeight: %.2f", reversalWeight)
	weight += reversalWeight

	crossBandWeight := crossBandWeight(result)
	weightLog += fmt.Sprintf(", crossBandWeight: %.2f", crossBandWeight)
	weight += crossBandWeight

	if masterTrend == models.TREND_UP {
		weight += 0.1
		weightLog += fmt.Sprintf(", masterTrenWeight: %.2f", 0.1)
	}

	if result.Trend == models.TREND_UP {
		weight += 0.1
		weightLog += fmt.Sprintf(", TrenWeight: %.2f", 0.1)
	}

	return weight
}

func CalculateWeightLongInterval(result *models.BandResult, masterTrend int8) float32 {
	var weight float32 = 0

	weight += getPositionWeight(result.Bands, result.Trend, masterTrend, true)

	weight += getPriceMarginWithUpperBandWeight(result.Bands)

	weight += getPatternWeight(result)

	weight += reversalWeight(result)

	weight += crossBandWeight(result)

	if result.Trend == models.TREND_UP {
		weight += 0.2
	}

	return weight
}

func GetWeightLog() string {
	return weightLog
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
	trend := CalculateTrends(result.Bands[:len(result.Bands)-1])

	lastFiveData := result.Bands[len(result.Bands)-4:]
	if trend == models.TREND_UP || CalculateTrends(lastFiveData[1:]) != models.TREND_UP || result.PriceChanges < 1.5 {
		return 0
	}

	lastBand := lastFiveData[0]
	if trend == models.TREND_SIDEWAY {
		if result.Position == models.ABOVE_UPPER && lastBand.Candle.Open >= float32(lastBand.Upper) {
			return 0.1
		}

		if result.Position == models.ABOVE_SMA && lastBand.Candle.Open >= float32(lastBand.SMA) {
			return 0.1
		}

		if result.Position == models.BELOW_SMA && lastBand.Candle.Open >= float32(lastBand.Lower) {
			return 0.1
		}
	}

	firstBand := lastFiveData[0]
	secondBand := lastFiveData[1]
	thirdBand := lastFiveData[2]
	isBandCrossWithLower := firstBand.Candle.Low <= float32(firstBand.Lower) || secondBand.Candle.Low <= float32(secondBand.Lower) || thirdBand.Candle.Low <= float32(thirdBand.Lower)
	isBandCrossWithSMA := firstBand.Candle.Low <= float32(firstBand.SMA) || secondBand.Candle.Low <= float32(secondBand.SMA) || thirdBand.Candle.Low <= float32(thirdBand.SMA)
	isBandCrossWithUpper := firstBand.Candle.Low <= float32(firstBand.Upper) || secondBand.Candle.Low <= float32(secondBand.Upper) || thirdBand.Candle.Low <= float32(thirdBand.Upper)
	if isBandCrossWithLower && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].Lower {
		return 0.35
	} else if isBandCrossWithSMA && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].SMA {
		return 0.3
	} else if isBandCrossWithUpper && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].Upper {
		return 0.25
	}

	return 0.2
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

	return 0.3
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

func getPositionWeight(bands []models.Band, trend int8, masterTrend int8, isLongInterval bool) float32 {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]

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
		if trend == models.TREND_UP {
			return 0.325
		}
	}

	// hight menyentuh SMA tp close dibaawh SMA
	if lastBand.Candle.Hight >= float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.225
		}
	}

	// close menyentuh SMA tp open dibaawh SMA
	if lastBand.Candle.Close >= float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.4
		}
	}

	// open menyentuh SMA tp low dibaawh SMA
	if lastBand.Candle.Open >= float32(lastBand.SMA) && lastBand.Candle.Low < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.475
		}
	}

	// low hight dibawah Upper
	if lastBand.Candle.Hight < float32(lastBand.Upper) {
		if trend == models.TREND_UP {
			return 0.275
		}
	}

	if masterTrend != models.TREND_DOWN && !isLongInterval {

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
