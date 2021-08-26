package analysis

import "telebot-trading/app/models"

func CalculateWeight(result *models.BandResult, masterTrend int8) float32 {
	weight := result.PriceChanges
	if weight < 0.5 {
		return 0
	} else if weight > 1.5 {
		weight = 1.5
	}

	if result.VolumeChanges > 0 {
		weight += getVolumeAverageChangesWeight(result.VolumeChanges)
	}

	weight += getPositionWeight(result.Bands, result.Trend)

	weight += getPriceMarginWithUpperBandWeight(result.Bands)

	weight += getPatternWeight(result.Bands)

	weight += reversalWeight(result)

	if masterTrend == models.TREND_UP {
		weight += 0.15
	}

	if result.Trend == models.TREND_UP {
		weight += 0.15
	}

	return weight
}

func CalculateWeightLongInterval(result *models.BandResult) float32 {
	var weight float32 = 0

	weight += getPositionWeight(result.Bands, result.Trend)

	weight += getPriceMarginWithUpperBandWeight(result.Bands)

	weight += getPatternWeight(result.Bands)

	if result.Trend == models.TREND_UP {
		weight += 0.25
	}

	return weight
}

func reversalWeight(result *models.BandResult) float32 {
	lastFiveData := result.Bands[len(result.Bands)-5 : len(result.Bands)]
	if result.Trend != models.TREND_DOWN || CalculateTrends(lastFiveData) != models.TREND_UP {
		return 0
	}

	isBandCrossWithLower := lastFiveData[0].Candle.Low <= float32(lastFiveData[0].Lower) || lastFiveData[1].Candle.Low <= float32(lastFiveData[1].Lower)
	isBandCrossWithSMA := lastFiveData[0].Candle.Low <= float32(lastFiveData[0].SMA) || lastFiveData[1].Candle.Low <= float32(lastFiveData[1].SMA)
	isBandCrossWithUpper := lastFiveData[0].Candle.Low <= float32(lastFiveData[0].Upper) || lastFiveData[1].Candle.Low <= float32(lastFiveData[1].Upper)
	if isBandCrossWithLower && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].Lower {
		return 0.35
	} else if isBandCrossWithSMA && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].SMA {
		return 0.3
	} else if isBandCrossWithUpper && float64(lastFiveData[4].Candle.Low) > lastFiveData[4].Upper {
		return 0.25
	}

	return 0.1
}

func getPatternWeight(bands []models.Band) float32 {
	listMatchPattern := GetCandlePattern(bands)

	var weight float32 = 0
	if len(listMatchPattern) > 0 {
		weight += 0.5 * float32(len(listMatchPattern))
	}

	return weight
}

func getPriceMarginWithUpperBandWeight(bands []models.Band) float32 {
	lastBand := bands[len(bands)-1]
	var percent float32 = 2.1

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
	if percent >= 5.1 {
		return 0.5
	} else if percent >= 4.1 {
		return 0.45
	} else if percent >= 3.1 {
		return 0.375
	} else if percent >= 2.1 {
		return 0.3
	} else if percent >= 1.1 {
		return 0.225
	}

	return 0.15
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

func getPositionWeight(bands []models.Band, trend int8) float32 {
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

	// hight menyentuh Upper tp close dibaawh Upper
	if lastBand.Candle.Hight >= float32(lastBand.Upper) && lastBand.Candle.Close < float32(lastBand.Upper) {
		return 0.175
	}

	// close menyentuh Upper tp open dibaawh Upper
	if lastBand.Candle.Close >= float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) {
		return 0.35
	}

	// open menyentuh Upper tp low dibaawh Upper
	if lastBand.Candle.Open >= float32(lastBand.Upper) && lastBand.Candle.Low < float32(lastBand.Upper) {
		return 0.425
	}

	return 0.15
}
