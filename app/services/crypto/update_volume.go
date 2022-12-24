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
var countTotal int = 0
var countDown int = 0
var modeChecking string = ""

func StartUpdateVolumeService(updateVolumeChan chan bool) {
	for <-updateVolumeChan {
		currentTime := time.Now()
		if !IsProfitMoreThanThreshold() || (currentTime.Hour() == 23 && currentTime.Minute() > 45) {
			updateVolume(currentTime)
		}
	}
}

func updateVolume(currentTime time.Time) {
	if countDown > 0 {
		log.Println("skipped update volume worker, after massive update ")
		countDown--
		return
	}

	log.Println("starting update volume worker ")

	responseChan := make(chan CandleResponse)

	condition := map[string]interface{}{
		"status": models.STATUS_ACTIVE,
	}

	var limit *int = nil
	if currentTime.Minute() > 5 {
		if countTotal == 0 {
			countTotal = 300
		}
		temp := countTotal / 2
		limit = &temp
	} else {
		countTotal = 0
		countDown = 1
	}

	orderBy := "price_changes desc"
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, limit, nil, &orderBy)
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

		pricePercent := analysis.CalculateBandPriceChangesPercent(bollinger, direction, 21)
		vol := countVolume(response.CandleData[len(response.CandleData)-11:])
		if bollinger.AllTrend.ShortTrend != models.TREND_UP {
			pricePercent = -pricePercent
		}
		if bollinger.AllTrend.SecondTrend == models.TREND_UP && bollinger.AllTrend.ShortTrend == models.TREND_UP && pricePercent > 2.1 {
			countTrendUp++
		}

		if bollinger.AllTrend.SecondTrend == models.TREND_UP && bollinger.AllTrend.ShortTrend == models.TREND_UP && pricePercent > 4.1 {
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
	log.Println("total checked data: ", len(*currency_configs))

	if countTrendUpSignifican > 0 && countTrendUp/countTrendUpSignifican <= 6 && countTrendUp >= 25 {
		modeChecking = models.MODE_TREND_UP
	} else {
		modeChecking = models.MODE_TREND_NOT_UP
	}

	countLimit = countTrendUp
	if countTotal == 0 {
		countTotal = len(*currency_configs)
	}

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
