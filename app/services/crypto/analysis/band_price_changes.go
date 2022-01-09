package analysis

import "telebot-trading/app/models"

const BAND_UP int8 = 1
const BAND_DOWN int8 = 2

func CalculateBandPriceChangesPercent(bands models.Bands, direction int8) float32 {
	data := bands.Data

	dataLength := len(data)
	lastCandle := data[dataLength-1]
	threeCandleBeforeLast := data[dataLength-4 : dataLength-1]

	if CalculateTrendShort(data[len(data)-4:]) == models.TREND_DOWN {
		return percentDownCandle(threeCandleBeforeLast, lastCandle)
	}

	return percentUpCandle(threeCandleBeforeLast, lastCandle)
}

func percentUpCandle(threeCandleBeforeLast []models.Band, lastCandle models.Band) float32 {
	lowest := threeCandleBeforeLast[0].Candle.Close
	for _, val := range threeCandleBeforeLast {
		if lowest > val.Candle.Close {
			lowest = val.Candle.Close
		}
	}

	if lowest > lastCandle.Candle.Open {
		lowest = lastCandle.Candle.Open
	}

	changes := lastCandle.Candle.Close - lowest

	return changes / lowest * 100
}

func percentDownCandle(threeCandleBeforeLast []models.Band, lastCandle models.Band) float32 {
	higest := higestFromBand(threeCandleBeforeLast[0])
	for _, val := range threeCandleBeforeLast {
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
