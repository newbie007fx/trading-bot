package crypto

import (
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var countLimit int = 0
var modeChecking string = ""

func StartUpdateVolumeService(updateVolumeChan chan bool) {
	for <-updateVolumeChan {
		currentTime := time.Now()
		if !IsProfitMoreThanThreshold() || (currentTime.Hour() == 23 && currentTime.Minute() > 45) {
			updateVolume()
		}
	}
}

func updateVolume() {
	log.Println("starting update volume worker ")

	responseChan := make(chan CandleResponse)

	currency_configs := repositories.GetCurrencyNotifConfigs(nil, nil, nil)
	countTrendUp := 0
	countTrendUpSignifican := 0
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

		pricePercent := analysis.CalculateBandPriceChangesPercent(bollinger, direction, 11)
		vol := countVolume(response.CandleData[len(response.CandleData)-11:])

		if bollinger.AllTrend.SecondTrend == models.TREND_UP && pricePercent > 3 {
			countTrendUp++
		}

		if bollinger.AllTrend.SecondTrend == models.TREND_UP && pricePercent > 5 {
			countTrendUpSignifican++
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

	log.Println(fmt.Sprintf("count trend up %d, count significan trend up %d", countTrendUp, countTrendUpSignifican))

	if countTrendUp/countTrendUpSignifican <= 6 && countTrendUp >= 40 {
		modeChecking = "trend_up"
	} else {
		modeChecking = "not_trend_up"
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

func GetModeChecking() string {
	return modeChecking
}
