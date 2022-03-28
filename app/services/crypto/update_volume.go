package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"time"
)

func StartUpdateVolumeService(updateVolumeChan chan bool) {
	for <-updateVolumeChan {
		updateVolume()
	}
}

func updateVolume() {
	log.Println("starting update volume worker ")

	responseChan := make(chan CandleResponse)

	var limit *int = nil
	currentTime := time.Time{}
	if currentTime.Hour()%4 != 0 {
		defaultLimit := 50
		limit = &defaultLimit
	}

	currency_configs := repositories.GetCurrencyNotifConfigs(nil, limit, nil)
	for _, data := range *currency_configs {
		time.Sleep(2 * time.Second)

		request := CandleRequest{
			Symbol:       data.Symbol,
			Limit:        12,
			Resolution:   "1h",
			ResponseChan: responseChan,
		}

		DispatchRequestJob(request)

		response := <-responseChan
		if response.Err != nil {
			log.Println("error: ", response.Err.Error())
			continue
		}
		vol := countVolume(response.CandleData)
		pricePercent := priceChanges(response.CandleData)
		priceToVolume := vol + (vol * pricePercent / 100)

		err := repositories.UpdateCurrencyNotifConfig(data.ID, map[string]interface{}{
			"volume":        vol + priceToVolume,
			"price_changes": pricePercent,
		})

		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}
	}

	log.Println("update volume worker done")
}

func countVolume(candles []models.CandleData) float32 {
	var volume float32 = 0
	var lastPrice float32 = 1
	for _, candle := range candles {
		volume += candle.Volume
		lastPrice = candle.Close
	}

	return volume / float32(len(candles)) * lastPrice
}

func priceChanges(candles []models.CandleData) float32 {
	firstCandle := candles[0]
	lastCandle := candles[len(candles)-1]

	return (lastCandle.Close - firstCandle.Close) / firstCandle.Close * 100
}
