package analysis

import (
	"log"
	"telebot-trading/app/models"
)

func IsIgnored(result, masterCoin *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) {
		return true
	}

	if lastBandHeadDoubleBody(result) {
		return true
	}

	if isContaineBearishEngulfing(result) {
		return true
	}

	if whenHeadCrossBandAndMasterDown(result, masterCoin) {
		return true
	}

	return ignored(result, masterCoin)
}

func IsIgnoredLongInterval(result *models.BandResult) bool {
	if isInAboveUpperBandAndDownTrend(result) {
		return true
	}

	if lastBandHeadDoubleBody(result) {
		return true
	}

	return isContaineBearishEngulfing(result)
}

func IsIgnoredMasterDown(result, masterCoin *models.BandResult) bool {
	if result.Position != models.BELOW_LOWER && result.Position != models.BELOW_SMA {
		if result.Position == models.ABOVE_SMA && masterCoin.PriceChanges > 0.33 {
			return false
		}

		return true
	}

	return false
}

func isInAboveUpperBandAndDownTrend(result *models.BandResult) bool {
	lastFiveData := result.Bands[len(result.Bands)-5:]
	if isHeighestOnHalfEndAndAboveUpper(result) && CalculateTrends(lastFiveData) == models.TREND_DOWN {
		return true
	}

	return false
}

func isHeighestOnHalfEndAndAboveUpper(result *models.BandResult) bool {
	hiIndex := getHighestIndex(result)
	if hiIndex >= len(result.Bands)-5 {
		return result.Bands[hiIndex].Candle.Close > float32(result.Bands[hiIndex].Upper)
	}

	return false
}

func isContaineBearishEngulfing(result *models.BandResult) bool {
	hiIndex := getHighestIndex(result)
	if hiIndex >= len(result.Bands)/4 {
		return BearishEngulfing(result.Bands[hiIndex:])
	}

	return false
}

func getHighestIndex(result *models.BandResult) int {
	hiIndex := 0
	for i, band := range result.Bands {
		if result.Bands[hiIndex].Candle.Close < band.Candle.Close {
			hiIndex = i
		}
	}

	return hiIndex
}

func lastBandHeadDoubleBody(result *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Close > lastBand.Candle.Open && result.Position == models.ABOVE_UPPER {
		head := lastBand.Candle.Hight - lastBand.Candle.Close
		body := lastBand.Candle.Close - lastBand.Candle.Open
		return head > body
	}
	return false
}

func ignored(result, masterCoin *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	if lastBand.Candle.Low <= float32(lastBand.SMA) && lastBand.Candle.Hight >= float32(lastBand.Upper) {
		log.Println("reset to 0 with criteria 1")
		return true
	}

	highest := getHigestPrice(result.Bands)
	lowest := getLowestPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent < 2 {
		log.Println("reset to 0 with criteria 2")
		return true
	}

	if masterCoin.Trend == models.TREND_DOWN {
		lastFourData := result.Bands[len(result.Bands)-4:]
		if CalculateTrends(lastFourData) != models.TREND_UP {
			log.Println("reset to 0 with criteria 3")
			return true
		}
	}

	return false
}

func whenHeadCrossBandAndMasterDown(result, masterCoin *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	crossSMA := lastBand.Candle.Close < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA)
	crossUpper := lastBand.Candle.Close < float32(lastBand.Upper) && lastBand.Candle.Hight > float32(lastBand.Upper)
	if (crossSMA || crossUpper) && CalculateTrends(masterCoin.Bands[len(masterCoin.Bands)-5:]) == models.TREND_DOWN {
		return true
	}

	return false
}
