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
			Resolution:   "1h",
			ResponseChan: responseChan,
		}

		DispatchRequestJob(request)

		response := <-responseChan
		if response.Err != nil {
			log.Println("error: ", response.Err.Error())
			continue
		}

		vol := countVolume(response.CandleData[len(response.CandleData)-4:])
		pricePercent := priceChanges(response.CandleData[len(response.CandleData)-4:])
		priceToVolume := vol + (vol * pricePercent / 100)

		bollinger := analysis.GenerateBollingerBands(response.CandleData)
		if bollinger.AllTrend.SecondTrend == models.TREND_UP && bollinger.AllTrend.ShortTrend == models.TREND_UP {
			countTrendUp++
		}

		err := repositories.UpdateCurrencyNotifConfig(data.ID, map[string]interface{}{
			"volume":        vol + priceToVolume,
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

	var firstOpen float32 = 0
	if firstCandle.Open < firstCandle.Close {
		firstOpen = firstCandle.Open
	} else {
		firstOpen = firstCandle.Close
	}

	return (lastCandle.Close - firstOpen) / firstOpen * 100
}

func GetLimit() int {
	return countLimit
}
