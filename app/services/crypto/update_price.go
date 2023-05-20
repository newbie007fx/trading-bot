package crypto

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var countLimit int = 1
var limit int = 200
var modeChecking string = ""

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
	currentTime := time.Now()
	endDate := getEndDate(currentTime)

	responseChan := make(chan CandleResponse)

	orderBy := "volume desc"
	condition := map[string]interface{}{
		"is_on_hold = ?": false,
		"status = ?":     models.STATUS_ACTIVE,
	}
	ignoredCoins := services.GetIgnoredCurrencies()
	currencyConfigs := repositories.GetCurrencyNotifConfigsIgnoredCoins(&condition, &limit, ignoredCoins, &orderBy)

	countTrendUp := 0
	countTrendUpSignifican := 0
	trendUpCoins := ""
	for i, data := range *currencyConfigs {
		if i%20 == 0 {
			time.Sleep(1 * time.Second)
		}

		var result *models.BandResult
		request := CandleRequest{
			Symbol:       data.Symbol,
			EndDate:      endDate,
			Limit:        40,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result = MakeCryptoRequest(request)

		if result == nil {
			log.Println("empty result ")
			continue
		}

		pricePercent := result.PriceChanges
		if result.AllTrend.ShortTrend != models.TREND_UP {
			pricePercent = -pricePercent
		}
		if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && pricePercent > 2.5 {
			trendUpCoins = trendUpCoins + ", " + data.Symbol
			countTrendUp++
		}

		if result.AllTrend.Trend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && pricePercent > 4 {
			countTrendUpSignifican++
		}

		err := repositories.UpdateCurrencyNotifConfig(data.ID, map[string]interface{}{
			"price_changes": pricePercent,
		})

		if err != nil {
			log.Println("error: ", err.Error())
		}
		if result.Direction == analysis.BAND_DOWN {
			services.SetIgnoredCurrency(data.Symbol, 1)
		}
	}

	log.Printf("count trend up %d, count significan trend up %d\n", countTrendUp, countTrendUpSignifican)
	log.Println("list trend up coin: ", trendUpCoins)
	log.Println("total checked data: ", len(*currencyConfigs))

	if countTrendUpSignifican > 0 && countTrendUp/countTrendUpSignifican <= 9 && countTrendUp >= 15 {
		modeChecking = models.MODE_TREND_UP
	} else {
		modeChecking = models.MODE_TREND_NOT_UP
	}

	countLimit = countTrendUp

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

func getEndDate(baseTime time.Time) int64 {
	unixTime := baseTime.Unix()

	if baseTime.Minute()%15 == 0 {
		unixTime = unixTime - (int64(baseTime.Second()) % 60) - 1
	}

	return unixTime * 1000
}
