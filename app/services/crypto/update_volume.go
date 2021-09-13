package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
)

func StartUpdateVolumeService(updateVolumeChan chan bool) {
	for <-updateVolumeChan {
		updateVolume()
	}
}

func updateVolume() {
	log.Println("starting update volume worker ")

	responseChan := make(chan CandleResponse)

	currency_configs := repositories.GetCurrencyNotifConfigs(nil, nil)
	for _, data := range *currency_configs {
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
