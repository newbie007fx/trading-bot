package indicator

import "telebot-trading/internal/model"

func ExtractClosePrices(candles []model.CandleData) []float64 {
	closes := make([]float64, 0, len(candles))
	for _, c := range candles {
		closes = append(closes, c.Close)
	}
	return closes
}

func LastN(values []float64, n int) []float64 {
	if len(values) < n {
		return nil
	}
	return values[len(values)-n:]
}
