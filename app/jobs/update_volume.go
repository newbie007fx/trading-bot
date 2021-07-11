package jobs

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"time"
)

func UpdateVolume() {
	counter := 0
	currentTime := time.Now().Unix()
	startTime := currentTime - (60 * 60 * 12) - 60

	currency_configs := repositories.GetCurrencyNotifConfigs(nil)
	for _, data := range *currency_configs {

		checkCounter(&counter, currentTime)

		crypto := services.GetCrypto()
		candles, err := crypto.GetCandlesData(data.Symbol, startTime-60, currentTime, "60")
		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}
		vol := countVolume(candles)

		repositories.UpdateCurrencyNotifConfig(data.ID, map[string]interface{}{
			"volume": vol,
		})
	}
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

func checkCounter(counter *int, startTime int64) {
	currentTime := time.Now().Unix()
	difference := currentTime - startTime
	sleep := 0
	quotaLeft := difference % 60
	if *counter == 60 {
		if quotaLeft > 0 {
			sleep = int(quotaLeft) + 1
		}
		*counter = 0
	}

	if difference%60 == 0 {
		*counter = 0
		sleep = 1
	}

	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	*counter++
}
