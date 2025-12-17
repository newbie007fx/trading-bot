package indicator

import "errors"

func RSISeries(closes []float64, period int) ([]float64, error) {
	if len(closes) < period+1 {
		return nil, errors.New("not enough data for RSI")
	}

	rsi := make([]float64, len(closes))

	var gainSum, lossSum float64

	// 1️⃣ Initial SMA (first period)
	for i := 1; i <= period; i++ {
		delta := closes[i] - closes[i-1]
		if delta >= 0 {
			gainSum += delta
		} else {
			lossSum -= delta
		}
	}

	avgGain := gainSum / float64(period)
	avgLoss := lossSum / float64(period)

	// First RSI value
	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// 2️⃣ Wilder smoothing
	for i := period + 1; i < len(closes); i++ {
		delta := closes[i] - closes[i-1]

		var gain, loss float64
		if delta >= 0 {
			gain = delta
		} else {
			loss = -delta
		}

		avgGain = (avgGain*(float64(period-1)) + gain) / float64(period)
		avgLoss = (avgLoss*(float64(period-1)) + loss) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi, nil
}
