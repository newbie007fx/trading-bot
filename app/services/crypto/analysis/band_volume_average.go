package analysis

import "telebot-trading/app/models"

func CalculateVolumeAverage(data []models.Band) float32 {
	if len(data) < 8 {
		return 0
	}

	lastCandle := data[len(data)-1].Candle

	var total float32 = 0
	for _, val := range data {
		total += val.Candle.Volume
	}
	average := total / float32(len(data)-1)
	if lastCandle.Volume > average {
		difference := lastCandle.Volume - average
		return difference / average * 100
	} else {
		difference := average - lastCandle.Volume
		return -(difference / average * 100)
	}
}
