package analysis

import (
	"log"
	"math"
	"strconv"
	"strings"
	"telebot-trading/app/models"
)

const SMA_DAYS = 20

const STANDARD_DEVIATIONS = 2

func GenerateBollingerBands(historical []models.CandleData) (bands models.Bands) {
	start := 0
	end := SMA_DAYS

	bands.Data = []models.Band{}

	if len(historical) < SMA_DAYS {
		log.Println("invalid historycal data with len: ", len(historical))
	}

	graphData := len(historical) - SMA_DAYS
	for i := 0; i <= graphData; i++ {
		bands.Data = append(bands.Data, getBandData(historical[start:end]))
		start++
		end++
	}

	bands.AllTrend = CalculateTrendsDetail(bands.Data)
	bands.Position = getPosition(bands.Data[len(bands.Data)-1])

	return
}

func getPosition(band models.Band) int8 {
	position := models.BELOW_LOWER
	if band.Candle.Close >= float32(band.Upper) {
		position = models.ABOVE_UPPER
	} else if band.Candle.Close >= float32(band.SMA) {
		position = models.ABOVE_SMA
	} else if band.Candle.Close >= float32(band.Lower) {
		position = models.BELOW_SMA
	}

	return position
}

func getBandData(historical []models.CandleData) (result models.Band) {
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

	candle := &historical[size-1]

	return models.Band{
		Candle: candle,
		SMA:    updatePrecision(candle, sma),
		Upper:  updatePrecision(candle, upper),
		Lower:  updatePrecision(candle, lower),
	}
}

func updatePrecision(data *models.CandleData, value float64) float64 {
	basePrecision := getMaxNumPrec(data)
	multiplier := math.Pow(10, float64(basePrecision))

	return math.Floor(value*multiplier) / multiplier
}

func getMaxNumPrec(data *models.CandleData) float64 {
	base := numFloatPlaces(float64(data.Low))
	tmpBase := numFloatPlaces(float64(data.Open))
	if tmpBase > base {
		base = tmpBase
	}

	tmpBase = numFloatPlaces(float64(data.Close))
	if tmpBase > base {
		base = tmpBase
	}

	tmpBase = numFloatPlaces(float64(data.Hight))
	if tmpBase > base {
		base = tmpBase
	}

	return float64(base)
}

func numFloatPlaces(v float64) int {
	s := strconv.FormatFloat(v, 'f', -1, 32)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}

	return 0
}
