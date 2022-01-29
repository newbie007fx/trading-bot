package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/services/crypto/analysis"
)

func MakeCryptoRequest(data models.CurrencyNotifConfig, request CandleRequest) *models.BandResult {
	DispatchRequestJob(request)

	response := <-request.ResponseChan
	if response.Err != nil {
		log.Println("error: ", response.Err.Error())
		return nil
	}

	if len(response.CandleData) < 20 {
		log.Println("invalid candle data value")
	}

	bands := analysis.GetCurrentBollingerBands(response.CandleData)
	if len(bands.Data) < 13 {
		log.Println("invalid number of band, skipped")
		return nil
	}

	direction := analysis.BAND_UP
	if !analysis.CheckLastCandleIsUp(bands.Data) {
		direction = analysis.BAND_DOWN
	}

	lastBand := bands.Data[len(bands.Data)-1]

	result := models.BandResult{
		Symbol:        request.Symbol,
		Direction:     direction,
		CurrentPrice:  lastBand.Candle.Close,
		CurrentVolume: lastBand.Candle.Volume,
		AllTrend:      bands.AllTrend,
		PriceChanges:  bands.PriceChanges,
		VolumeChanges: bands.VolumeAverageChanges,
		Position:      bands.Position,
		Bands:         bands.Data,
	}

	return &result
}

func MakeCryptoRequestUpdateLasCandle(data models.CurrencyNotifConfig, request CandleRequest, close, hight, low float32) *models.BandResult {
	DispatchRequestJob(request)

	response := <-request.ResponseChan
	if response.Err != nil {
		log.Println("error: ", response.Err.Error())
		return nil
	}

	if len(response.CandleData) < 20 {
		log.Println("invalid candle data value")
	}

	bands := analysis.GetCurrentBollingerBands(updateLastCandle(response.CandleData, close, hight, low))
	if len(bands.Data) < 13 {
		log.Println("invalid number of band, skipped")
		return nil
	}

	direction := analysis.BAND_UP
	if !analysis.CheckLastCandleIsUp(bands.Data) {
		direction = analysis.BAND_DOWN
	}

	lastBand := bands.Data[len(bands.Data)-1]

	result := models.BandResult{
		Symbol:        request.Symbol,
		Direction:     direction,
		CurrentPrice:  lastBand.Candle.Close,
		CurrentVolume: lastBand.Candle.Volume,
		AllTrend:      bands.AllTrend,
		PriceChanges:  bands.PriceChanges,
		VolumeChanges: bands.VolumeAverageChanges,
		Position:      bands.Position,
		Bands:         bands.Data,
	}

	return &result
}

func updateLastCandle(candles []models.CandleData, close, hight, low float32) []models.CandleData {
	lastCandle := candles[len(candles)-1]
	lastCandle.Close = close
	if hight > lastCandle.Open {
		lastCandle.Hight = hight
	} else {
		lastCandle.Hight = lastCandle.Open
	}

	if low < lastCandle.Open {
		lastCandle.Low = low
	} else {
		lastCandle.Low = lastCandle.Open
	}
	candles[len(candles)-1] = lastCandle

	return candles
}
