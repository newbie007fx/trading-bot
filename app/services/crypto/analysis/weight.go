package analysis

import (
	"telebot-trading/app/models"
)

var weightLogData map[string]float32

func CalculateWeight(result *models.BandResult) float32 {
	weightLogData = map[string]float32{}

	positionWeight := getPositionWeight(result.Bands, false)
	weightLogData["positionWeight"] = positionWeight
	weight := positionWeight

	priceMarginWeight := getPriceMarginWithUpperBandWeight(result.Bands)
	weightLogData["priceMarginWeight"] = priceMarginWeight
	weight += priceMarginWeight

	return weight
}

func GetWeightLogData() map[string]float32 {
	return weightLogData
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
