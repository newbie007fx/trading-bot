package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var countLimit int = 0

func StartUpdateVolumeService(updateVolumeChan chan bool) {
	for <-updateVolumeChan {
		updateVolume()
	}
}

func updateVolume() {
	log.Println("starting update volume worker ")

	responseChan := make(chan CandleResponse)

	currency_configs := repositories.GetCurrencyNotifConfigs(nil, nil, nil)
	countTrendUp := 0
	for i, data := range *currency_configs {
		if i%15 == 0 {
			time.Sleep(1 * time.Second)
		}

		request := CandleRequest{
			Symbol:       data.Symbol,
			Limit:        40,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		DispatchRequestJob(request)

		response := <-responseChan
		if response.Err != nil {
			log.Println("error: ", response.Err.Error())
			continue
		}
		bollinger := analysis.GenerateBollingerBands(response.CandleData)
		direction := analysis.BAND_DOWN
		if analysis.CheckLastCandleIsUp(bollinger.Data) {
			direction = analysis.BAND_UP
		}

		pricePercent := analysis.CalculateBandPriceChangesPercent(bollinger, direction)
		vol := countVolume(response.CandleData[len(response.CandleData)-8:])

		if bollinger.AllTrend.SecondTrend == models.TREND_UP && bollinger.AllTrend.Trend == models.TREND_UP {
			countTrendUp++
		}

		err := repositories.UpdateCurrencyNotifConfig(data.ID, map[string]interface{}{
			"volume":        vol,
			"price_changes": pricePercent,
		})

		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}
	}

	countLimit = countTrendUp

	log.Println("update volume worker done")
}

func countVolume(candles []models.CandleData) float32 {
	var volume float32 = 0
	for _, candle := range candles {
		volume += candle.Volume
	}

	return volume
}

func GetLimit() int {
	return countLimit
}
