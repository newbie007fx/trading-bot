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

	if masterTrend == models.TREND_UP {
		weight += 0.15
	}

	return weight
}

func getPatternWeight(bands []models.Band) float32 {
	listMatchPattern := GetCandlePattern(bands)

	var weight float32 = 0
	for _ = range listMatchPattern {
		weight += 0.5
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
		return 0.4
	} else if percent >= 3.1 {
		return 0.325
	} else if percent >= 2.1 {
		return 0.25
	} else if percent >= 1.1 {
		return 0.175
	} else if percent >= 0.1 {
		return 0.1
	}

	return 0
}

func getVolumeAverageChangesWeight(volumeAverageChanges float32) float32 {
	if volumeAverageChanges >= 101 {
		return 0.3
	} else if volumeAverageChanges >= 81 {
		return 0.25
	} else if volumeAverageChanges >= 61 {
		return 0.2
	} else if volumeAverageChanges >= 41 {
		return 0.15
	} else if volumeAverageChanges >= 21 {
		return 0.1
	} else if volumeAverageChanges >= 1 {
		return 0.05
	}

	return 0
}

func getPositionWeight(bands []models.Band, trend int8) float32 {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]

	// low hight dibawah lower
	if lastBand.Candle.Hight < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.2
		}
	}

	// hight menyentuh lower tp close dibaawh lower
	if lastBand.Candle.Hight >= float32(lastBand.Lower) && lastBand.Candle.Close < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.15
		}
	}

	// close menyentuh lower tp open dibaawh lower
	if lastBand.Candle.Close >= float32(lastBand.Lower) && lastBand.Candle.Open < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.25
		}
	}

	// open menyentuh lower tp low dibaawh lower
	if lastBand.Candle.Open >= float32(lastBand.Lower) && lastBand.Candle.Low < float32(lastBand.Lower) {
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close {
			return 0.35
		}
	}

	// low hight dibawah SMA
	if lastBand.Candle.Hight < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.18
		}
	}

	// hight menyentuh SMA tp close dibaawh SMA
	if lastBand.Candle.Hight >= float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.13
		}
	}

	// close menyentuh SMA tp open dibaawh SMA
	if lastBand.Candle.Close >= float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.23
		}
	}

	// open menyentuh SMA tp low dibaawh SMA
	if lastBand.Candle.Open >= float32(lastBand.SMA) && lastBand.Candle.Low < float32(lastBand.SMA) {
		if trend == models.TREND_UP {
			return 0.33
		}
	}

	// low hight dibawah Upper
	if lastBand.Candle.Hight < float32(lastBand.Upper) {
		if trend == models.TREND_UP {
			return 0.15
		}
	}

	// hight menyentuh Upper tp close dibaawh Upper
	if lastBand.Candle.Hight >= float32(lastBand.Upper) && lastBand.Candle.Close < float32(lastBand.Upper) {
		return 0.1
	}

	// close menyentuh Upper tp open dibaawh Upper
	if lastBand.Candle.Close >= float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) {
		return 0.2
	}

	// open menyentuh Upper tp low dibaawh Upper
	if lastBand.Candle.Open >= float32(lastBand.Upper) && lastBand.Candle.Low < float32(lastBand.Upper) {
		return 0.3
	}

	return 0.13
}
