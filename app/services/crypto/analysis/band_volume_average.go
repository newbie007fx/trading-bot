package analysis

import "telebot-trading/app/models"

func CalculateVolumeAverage(data []models.Band) float32 {
	if len(data) < 3 {
		return 0
	}

	lastVolume := data[len(data)-2].Candle.BuyVolume
	if lastVolume < data[len(data)-1].Candle.BuyVolume {
		lastVolume = data[len(data)-1].Candle.BuyVolume
	}

	var total float32 = 0
	for i := len(data) - 1; i > len(data)-4; i-- {
		total += data[i-1].Candle.BuyVolume
	}
	average := total / float32(3)
	if lastVolume > average {
		difference := lastVolume - average
		return difference / average * 100
	} else {
		difference := average - lastVolume
		return -(difference / average * 100)
	}
}
