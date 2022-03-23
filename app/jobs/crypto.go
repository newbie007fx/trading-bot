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
			log.Println("executing")
			strategy.Execute(currentTime)
		}

		if isTimeToUpdateVolume(currentTime) {
			updateVolumeChan <- true
		}

		sleep := 60 - time.Now().Second()
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
	hour := time.Hour()
	if (hour == 0 || hour == 3 || hour == 6 || hour == 9 || hour == 12 || hour == 15 || hour == 18 || hour == 21) && minute == 3 {
		return true
	}

	return false
}
