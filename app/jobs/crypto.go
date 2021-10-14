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
		if !isMuted() {

			strategy.Execute(currentTime)
		}

		if isTimeToUpdateVolume(currentTime) {
			updateVolumeChan <- true
		}

		log.Println("sleep")
		sleep := 60 - currentTime.Second()
		time.Sleep(time.Duration(sleep) * time.Second)
		log.Println("recheck")
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
	} else {
		strategy = &trading_strategy.ManualTradingStrategy{}
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
	hour := time.Hour()
	if (hour == 5 || hour == 17) && minute == 1 {
		return true
	}

	return false
}
