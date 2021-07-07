package services

import (
	"fmt"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
)

func GetCurrentBollingerBands(symbol string, time int64) (bands models.Bands, err error) {
	end := time
	start := end - (60 * 15 * 27)

	crypto := GetCrypto()
	candlesData, err := crypto.GetCandlesData(symbol, start, end)
	if err == nil {
		bands = GenerateBollingerBands(candlesData)
		direction := BAND_DOWN
		if CheckLastCandleIsUp(bands.Data) {
			direction = BAND_UP
		}

		bands.PriceChanges = CalculateBandPriceChangesPercent(bands.Data, direction)
		bands.VolumeAverageChanges = CalculateVolumeAverage(bands.Data)
	}

	return
}

func CheckLastCandleIsUp(bollingerBands []models.Band) bool {
	//candle posisi sekarang up, close diatas open
	size := len(bollingerBands)
	if size > 0 {
		candle := bollingerBands[size-1].Candle
		if candle.Close > candle.Open {
			return true
		}
	}

	return false
}

func CheckPositionOnUpperBand(bollingerBands []models.Band) bool {
	//candle posisi sekrang  diupper band
	size := len(bollingerBands)
	if size > 0 {
		band := bollingerBands[size-1]
		if band.Candle.Open <= float32(band.Upper) && float32(band.Upper) <= band.Candle.Close {
			return true
		}
	}

	return false
}

func CheckPositionSMAAfterLower(bands models.Bands) bool {
	//candle posisi sekarang diatas sma, trend up.
	lastBand := bands.Data[len(bands.Data)-1]
	if lastBand.Candle.Open <= float32(lastBand.SMA) && float32(lastBand.SMA) <= lastBand.Candle.Close {
		if bands.Trend == models.TREND_UP {
			return true
		}
	}

	return false
}

func CheckPositionAfterLower(bollingerBands []models.Band) bool {
	//candle posisi ditas lower band. setelah open/low dibaawah lowerband atau candle sebelumnya meyentuh loweband
	size := len(bollingerBands)
	if size > 2 {
		lastBand := bollingerBands[size-1]
		isLastBandOnLower := lastBand.Candle.Open <= float32(lastBand.Lower) && float32(lastBand.Lower) <= lastBand.Candle.Close
		isLastBandBelowLower := lastBand.Candle.Open <= float32(lastBand.Lower) && lastBand.Candle.Close <= float32(lastBand.Lower)
		if isLastBandOnLower || isLastBandBelowLower {
			return true
		}

		secondLastBand := bollingerBands[size-2]
		if secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Low < float32(secondLastBand.Lower) {
			return true
		}
	}

	return false
}

func IsPriceIncreaseAboveThreshold(bands models.Bands, isMaster bool) bool {
	var threshold float32 = 3
	if isMaster {
		threshold = 1.5
	}

	return bands.PriceChanges > threshold
}

func IsTrendUpAfterTrendDown(symbol string, bands models.Bands) bool {
	result := false
	var trend int8 = 0
	if bands.Trend == models.TREND_UP {
		resultString := repositories.GetConfigValueByName(symbol)
		if resultString != nil {
			tmp, err := strconv.ParseInt(*resultString, 10, 8)
			if err == nil {
				trend = int8(tmp)
			}
		}
	}

	result = trend == models.TREND_DOWN

	repositories.SetConfigByName(symbol, fmt.Sprint(bands.Trend))

	return result
}

func IsTrendDownAfterTrendUp(symbol string, bands models.Bands) bool {
	result := false
	var trend int8 = 0
	if bands.Trend == models.TREND_DOWN {
		resultString := repositories.GetConfigValueByName(symbol)
		if resultString != nil {
			tmp, err := strconv.ParseInt(*resultString, 10, 8)
			if err == nil {
				trend = int8(tmp)
			}
		}
	}

	result = trend == models.TREND_UP

	repositories.SetConfigByName(symbol, fmt.Sprint(bands.Trend))

	return result
}

func CheckPositionOnLowerBand(bollingerBands []models.Band) bool {
	//candle posisi sekrang  dilower band
	size := len(bollingerBands)
	if size > 0 {
		band := bollingerBands[size-1]
		if band.Candle.Open >= float32(band.Lower) && float32(band.Lower) >= band.Candle.Close {
			return true
		}
	}

	return false
}

func CheckPositionSMAAfterUpper(bands models.Bands) bool {
	//candle posisi sekarang diatas sma, trend down.
	lastBand := bands.Data[len(bands.Data)-1]
	if lastBand.Candle.Open >= float32(lastBand.SMA) && float32(lastBand.SMA) >= lastBand.Candle.Close {
		if bands.Trend == models.TREND_DOWN {
			return true
		}
	}

	return false
}

func CheckPositionAfterUpper(bollingerBands []models.Band) bool {
	//candle posisi dibawah upper band. setelah open/hight diatas upper band atau candle sebelumnya meyentuh upper band
	size := len(bollingerBands)
	if size > 2 {
		lastBand := bollingerBands[size-1]
		isLastBandOnUpper := lastBand.Candle.Open >= float32(lastBand.Upper) && float32(lastBand.Upper) >= lastBand.Candle.Close
		isLastBandAboveUpper := lastBand.Candle.Open >= float32(lastBand.Upper) && lastBand.Candle.Close >= float32(lastBand.Upper)
		if isLastBandOnUpper || isLastBandAboveUpper {
			return true
		}

		secondLastBand := bollingerBands[size-2]
		if secondLastBand.Candle.Open < secondLastBand.Candle.Close && secondLastBand.Candle.Hight > float32(secondLastBand.Upper) {
			return true
		}
	}

	return false
}

func IsPriceDecreasebelowThreshold(bands models.Bands, isMaster bool) bool {
	var threshold float32 = 2
	if isMaster {
		threshold = 1
	}

	return bands.PriceChanges > threshold
}
