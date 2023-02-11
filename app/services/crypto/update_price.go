package crypto

import (
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var countLimit int = 1
var offset int = 0
var limit int = 130
var modeChecking string = ""
var autoDebugMode bool = false

func StartUpdatePriceService(updatePriceChan chan bool) {
	for <-updatePriceChan {
		currentTime := time.Now()
		if !IsProfitMoreThanThreshold() || (currentTime.Hour() == 23 && currentTime.Minute() > 45) {
			updatePrice()
		}
	}
}

func updatePrice() {
	log.Println("starting update price worker ")

	responseChan := make(chan CandleResponse)

	condition := map[string]interface{}{
		"status": models.STATUS_ACTIVE,
	}

	if offset >= countTotal || offset >= 250 {
		offset = 0
	}

	orderBy := "volume desc"
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit, &offset, &orderBy)
	countTrendUp := 0
	countTrendUpSignifican := 0
	for i, data := range *currency_configs {
		if i%20 == 0 {
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
		if autoDebugMode {
			log.Println("coin: ", data.Symbol)
			log.Printf("response: %v", response)
		}
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
			"price_changes": pricePercent,
		})

		if err != nil {
			log.Println("error: ", err.Error())
		}
	}

	log.Println(fmt.Sprintf("count trend up %d, count significan trend up %d", countTrendUp, countTrendUpSignifican))
	log.Println("total checked data: ", len(*currency_configs))

	if countTrendUp == 0 && countTrendUpSignifican == 0 && countLimit == 0 {
		autoDebugMode = true
	} else {
		autoDebugMode = false
	}

	if countTrendUpSignifican > 0 && countTrendUp/countTrendUpSignifican <= 6 && countTrendUp >= 10 {
		modeChecking = models.MODE_TREND_UP
	} else {
		modeChecking = models.MODE_TREND_NOT_UP
	}

	countLimit = countTrendUp
	offset += limit

	log.Println("update price worker done")
}

func GetLimit() int {
	if countLimit > 0 {
		return countLimit
	}

	return 1
}

func GetModeChecking() string {
	return modeChecking
}
