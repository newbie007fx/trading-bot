package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"time"
)

func UpdateVolume() {
	log.Println("starting update volume worker ")

	responseChan := make(chan CandleResponse)

	condition := map[string]interface{}{
		"status": models.STATUS_ACTIVE,
	}

	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil, nil, nil)
	for _, data := range *currency_configs {
		time.Sleep(1 * time.Second)

		request := CandleRequest{
			Symbol:       data.Symbol,
			Limit:        24,
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
			"volume":        vol,
			"price_changes": 0,
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
		if candle.Open < candle.Close {
			volume += candle.Volume
		} else {
			volume -= candle.Volume
		}

		lastPrice = candle.Close
	}

	return volume * lastPrice
}
