package jobs

import (
	"log"
	"strconv"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/trading_strategy"
	"time"
)

var updateVolumeChan chan bool
var strategy trading_strategy.TradingStrategy

func StartCryptoWorker() {
	startService()

	defer func() {
		strategy.Shutdown()
		close(updateVolumeChan)
	}()

	for {
		currentTime := time.Now()
		second := time.Now().Second()

		if !isMuted() {
			log.Println("executing")
			strategy.Execute(currentTime)
		}

		if isTimeToUpdateVolume(currentTime) && second < 15 {
			updateVolumeChan <- true
		}

		if second >= 15 {
			second = second % 15
		}

		sleep := 15 - second
		log.Println("sleep: ", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func startService() {
	updateVolumeChan = make(chan bool)

	go crypto.RequestCandleService()
	go crypto.StartUpdateVolumeService(updateVolumeChan)
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

func isTimeToUpdateVolume(time time.Time) bool {
	minute := time.Minute()
	return minute == 13 || minute == 28 || minute == 43 || minute == 58
}
