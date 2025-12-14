package indicator

import "errors"

// EMASeries returns EMA value for each candle index
func EMASeries(prices []float64, period int) ([]float64, error) {
	if len(prices) < period {
		return nil, errors.New("not enough data for EMA")
	}

	ema := make([]float64, len(prices))

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	k := 2.0 / float64(period+1)

	for i := period; i < len(prices); i++ {
		ema[i] = prices[i]*k + ema[i-1]*(1-k)
	}

	return ema, nil
}
