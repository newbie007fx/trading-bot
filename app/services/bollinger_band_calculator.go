package services

import (
	"math"
)

type Band struct {
	Candle *CandleData
	SMA    float64
	Upper  float64
	Lower  float64
}

const SMA_DAYS = 20

const STANDARD_DEVIATIONS = 2

func GenerateBollingerBands(historical []CandleData, graphData int) (bands []Band) {
	start := 0
	end := SMA_DAYS

	for i := 0; i < graphData; i++ {
		bands = append(bands, getBandData(historical[start:end]))
		start++
		end++
	}

	return
}

func getBandData(historical []CandleData) (result Band) {
	size := len(historical)

	sum := float64(0)
	for _, h := range historical {
		sum += float64(h.Close)
	}

	sma := sum / float64(size)

	squares := float64(0)
	for i := 0; i < size; i++ {
		squares += math.Pow((float64(historical[i].Close) - sma), 2)
	}

	dev := math.Sqrt(squares / float64(size))

	upper := sma + (STANDARD_DEVIATIONS * dev)

	lower := sma - (STANDARD_DEVIATIONS * dev)

	return Band{&historical[size-1], sma, upper, lower}
}
