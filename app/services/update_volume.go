package services

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"time"
)

func StartUpdateVolumeService(updateVolumeChan chan bool) {
	for <-updateVolumeChan {
		updateVolume()
	}
}

func updateVolume() {
	log.Println("starting update volume worker ")

	currentTime := time.Now().Unix()
	startTime := currentTime - (60 * 60 * 12) - 60

	responseChan := make(chan crypto.CandleResponse)

	currency_configs := repositories.GetCurrencyNotifConfigs(nil, nil)
	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Start:        startTime,
			End:          currentTime,
			Resolution:   "60",
			ResponseChan: responseChan,
		}

		crypto.DispatchRequestJob(request)

		response := <-responseChan
		if response.Err != nil {
			log.Println("error: ", response.Err.Error())
			continue
		}
		vol := countVolume(response.CandleData)

		err := repositories.UpdateCurrencyNotifConfig(data.ID, map[string]interface{}{
			"volume": vol,
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
