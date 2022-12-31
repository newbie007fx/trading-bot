package jobs

import (
	"log"
	"strconv"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/trading_strategy"
	"time"
)

var updatePriceChan chan bool
var updateCurrencyChan chan bool
var strategy trading_strategy.TradingStrategy

func StartCryptoWorker() {
	startService()

	defer func() {
		strategy.Shutdown()
		close(updatePriceChan)
		close(updateCurrencyChan)
	}()

	for {
		currentTime := time.Now()
		second := currentTime.Second()

		if !isMuted() {
			strategy.Execute(currentTime)
		}

		if isTimeToUpdatePrice(currentTime) && second < 15 {
			updatePriceChan <- true
		}

		if currentTime.Hour() == 0 && currentTime.Minute() == 0 && second < 15 {
			log.Println("dispatch service update currency")
			updateCurrencyChan <- true
		}

		if second >= 15 {
			second = second % 15
		}

		sleep := 15 - second
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func startService() {
	updatePriceChan = make(chan bool)
	updateCurrencyChan = make(chan bool)

	go crypto.RequestCandleService()
	go crypto.StartUpdatePriceService(updatePriceChan)
	go crypto.StartUpdateCurrencyService(updateCurrencyChan)
	go crypto.StartSyncBalanceService()

	setStrategy()
	strategy.InitService()

	log.Println("waiting start up")
	time.Sleep(10 * time.Second)
}

func setStrategy() {
	mode := "manual"

	result := repositories.GetConfigValueByName("mode")
	if result != nil {
		mode = *result
	}

	if mode == "automatic" || mode == "simulation" {
		strategy = &trading_strategy.AutomaticTradingStrategy{}
	}

	log.Println("mode: " + mode)
}

func ChangeStrategy() {
	strategy.Shutdown()
	setStrategy()
	strategy.InitService()
}

func isMuted() bool {
	key := "is-muted"
	resultString := repositories.GetConfigValueByName(key)
	if resultString != nil {
		tmp, err := strconv.ParseBool(*resultString)
		if err != nil {
			return false
		}
		return tmp
	}
	return false
}

func isTimeToUpdatePrice(currentTime time.Time) bool {
	minute := currentTime.Minute()
	return minute == 3 || minute == 13 || minute == 18 || minute == 28 || minute == 33 || minute == 43 || minute == 48 || minute == 58
}
