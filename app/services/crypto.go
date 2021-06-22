package services

import (
	"time"
)

func GetCurrentBollingerBands(symbol string) (bands []Band, err error) {
	end := time.Now().Unix()
	start := end - (60 * 15 * 23)

	crypto := GetCrypto()
	candlesData, err := crypto.GetCandlesData(symbol, start, end)
	if err == nil {
		bands = GenerateBollingerBands(candlesData, 1)
	}

	return
}

func CheckLastCandleIsUp(bollingerBands []Band) bool {
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

func CheckPositionInUpperBand(bollingerBands []Band) bool {
	//candle posisi sekrang diatas upper band
	size := len(bollingerBands)
	if size > 0 {
		band := bollingerBands[size-1]
		if band.Candle.Open <= float32(band.Upper) && float32(band.Upper) <= band.Candle.Close {
			return true
		}
	}

	return false
}

func CheckPositionSMAAfterLower(bollingerBands []Band) bool {
	//candle posisi sekarang diatas sma, candle sebelumnya open close dibawah sma.
	size := len(bollingerBands)
	if size > 1 {
		lastBand := bollingerBands[size-1]
		if lastBand.Candle.Open <= float32(lastBand.SMA) && float32(lastBand.SMA) <= lastBand.Candle.Close {
			secondLastBand := bollingerBands[size-2]
			if secondLastBand.Candle.Open < float32(secondLastBand.SMA) && secondLastBand.Candle.Close < float32(secondLastBand.SMA) {
				return true
			}
		}
	}

	return false
}

func CheckPositionAfterLower(bollingerBands []Band) bool {
	//candle posisi ditas lower band. setelah open/low dibaawah lowerband atau candle sebelumnya meyentuh loweband
	size := len(bollingerBands)
	if size > 2 {
		lastBand := bollingerBands[size-1]
		if lastBand.Candle.Open <= float32(lastBand.Lower) && float32(lastBand.Lower) <= lastBand.Candle.Close {
			return true
		}

		secondLastBand := bollingerBands[size-2]
		if secondLastBand.Candle.Open >= float32(secondLastBand.Lower) && float32(secondLastBand.SMA) >= secondLastBand.Candle.Close {
			return true
		}

		thirdLastBand := bollingerBands[size-2]
		if thirdLastBand.Candle.Open >= float32(thirdLastBand.Lower) && float32(thirdLastBand.SMA) >= thirdLastBand.Candle.Close {
			return true
		}
	}

	return false
}
