package crypto

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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

	bands := analysis.GetCurrentBollingerBands(response.CandleData)

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
		Trend:         bands.Trend,
		PriceChanges:  bands.PriceChanges,
		VolumeChanges: bands.VolumeAverageChanges,
		Position:      bands.Position,
		Bands:         bands.Data,
	}

	if data.IsMaster || data.IsOnHold {
		if result.Direction == analysis.BAND_UP {
			result.Note = upTrendChecking(data, bands)
		} else {
			result.Note = downTrendChecking(data, bands, request)
		}
	}

	return &result
}

func upTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	if analysis.CheckPositionOnUpperBand(bands.Data) {
		return "Posisi naik upper band"
	}

	if analysis.CheckPositionSMAAfterLower(bands) {
		return "Posisi naik ke SMA"
	}

	if analysis.CheckPositionAfterLower(bands.Data) {
		return "Posisi lower"
	}

	if analysis.IsPriceIncreaseAboveThreshold(bands, data.IsMaster) {
		return "Naik diatas threshold"
	}

	if analysis.IsTrendUpAfterTrendDown(data.Symbol, bands) {
		return "Trend Up after down"
	}

	return ""
}

func downTrendChecking(data models.CurrencyNotifConfig, bands models.Bands, request CandleRequest) string {
	if analysis.CheckPositionOnLowerBand(bands.Data) {
		return "Posisi turun dibawah lower"
	}

	if analysis.CheckPositionSMAAfterUpper(bands) {
		return "Posisi turun dibawah SMA"
	}

	if analysis.CheckPositionAfterUpper(bands.Data) {
		return "Posisi turun dari Upper"
	}

	if analysis.IsPriceDecreasebelowThreshold(bands, data.IsMaster) {
		return "Turun dibawah threshold"
	}

	if analysis.IsTrendDownAfterTrendUp(data.Symbol, bands) {
		return "Trend Down after up"
	}

	if data.IsOnHold && (bands.Position == models.ABOVE_SMA || bands.Position == models.ABOVE_UPPER) {
		if isLastBandComplete(bands.Data, request.Resolution) {
			lastDown := countLastDownCandle(bands.Data)
			return fmt.Sprintf("Turun gan siaga !!! jumlah down %d", lastDown)
		}
	}

	return ""
}

func countLastDownCandle(data []models.Band) int {
	count := 0
	for i := len(data) - 1; i >= 0; i-- {
		band := data[i]
		if band.Candle.Close < band.Candle.Open {
			count++
		} else {
			break
		}
	}

	return count
}

func isLastBandComplete(bands []models.Band, resolution string) bool {
	lastBand := bands[len(bands)-1]

	closeTime := lastBand.Candle.CloseTime
	openTime := lastBand.Candle.OpenTime

	different := closeTime - openTime

	miliseconds := convertResolutionToSeconds(resolution)
	return different == miliseconds
}

func convertResolutionToSeconds(resolution string) int64 {
	var regex, err = regexp.Compile(`^([0-9]{1,2})([a-zA-Z])$`)

	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	var res = regex.FindStringSubmatch(resolution)

	base := 0
	switch res[2] {
	case "m":
		base = 60
	case "H":
		base = 60 * 60
	}
	i, _ := strconv.Atoi(res[1])

	return int64(base * i * 1000)
}
