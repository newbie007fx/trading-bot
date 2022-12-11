package analysis

import "telebot-trading/app/models"

const BAND_UP int8 = 1
const BAND_DOWN int8 = 2

func CalculateBandPriceChangesPercent(bands models.Bands, direction int8, candleLen int) float32 {
	data := bands.Data

	dataLength := len(data)
	lastCandle := data[dataLength-1]
	candleBeforeLast := data[dataLength-candleLen : dataLength-1]

	if bands.AllTrend.ShortTrend != models.TREND_UP {
		return percentDownCandle(candleBeforeLast, lastCandle)
	}

	return percentUpCandle(candleBeforeLast, lastCandle)
}

func percentUpCandle(candleBeforeLast []models.Band, lastCandle models.Band) float32 {
	lowest := candleBeforeLast[0].Candle.Close
	for _, val := range candleBeforeLast {
		if lowest > lowestFromBand(val) {
			lowest = lowestFromBand(val)
		}
	}

	if lowest > lastCandle.Candle.Open {
		lowest = lastCandle.Candle.Open
	}

	changes := lastCandle.Candle.Close - lowest

	return changes / lowest * 100
}

func percentDownCandle(candleBeforeLast []models.Band, lastCandle models.Band) float32 {
	higest := higestFromBand(candleBeforeLast[0])
	for _, val := range candleBeforeLast {
		if higest < higestFromBand(val) {
			higest = higestFromBand(val)
		}
	}

	if higest < lastCandle.Candle.Open {
		higest = lastCandle.Candle.Open
	}

	changes := higest - lastCandle.Candle.Close

	return changes / higest * 100
}
