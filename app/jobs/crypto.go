package jobs

import (
	"strconv"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/trading_strategy"
	"time"
)

var cryptoMasterCoinPriceChan chan bool
var cryptoHoldCoinPriceChan chan bool
var cryptoAltCoinPriceChan chan bool
var updateVolumeChan chan bool

func StartCryptoWorker() {
	startService()

	defer func() {
		close(cryptoMasterCoinPriceChan)
		close(cryptoHoldCoinPriceChan)
		close(cryptoAltCoinPriceChan)
		close(updateVolumeChan)
	}()

	for {
		currentTime := time.Now()
		if !isMuted() {
			minute := currentTime.Minute()
			if minute%5 == 0 {
				cryptoMasterCoinPriceChan <- true
				cryptoHoldCoinPriceChan <- true
			}

			if isTimeToCheckAltCoinPrice(currentTime) {
				cryptoAltCoinPriceChan <- true
			}
		}

		if isTimeToUpdateVolume(currentTime) {
			updateVolumeChan <- true
		}

		sleep := 60 - currentTime.Second()
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func startService() {
	cryptoMasterCoinPriceChan = make(chan bool)
	cryptoHoldCoinPriceChan = make(chan bool)
	cryptoAltCoinPriceChan = make(chan bool)
	updateVolumeChan = make(chan bool)

	go crypto.RequestCandleService()
	go trading_strategy.StartCheckMasterCoinPriceService(cryptoMasterCoinPriceChan)
	go trading_strategy.StartCheckHoldCoinPriceService(cryptoHoldCoinPriceChan)
	go trading_strategy.StartCheckAltCoinPriceService(cryptoAltCoinPriceChan)
	go services.StartUpdateVolumeService(updateVolumeChan)
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

func isTimeToCheckAltCoinPrice(time time.Time) bool {
	minute := time.Minute()
	if minute == 5 || minute == 20 || minute == 35 || minute == 50 {
		return true
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
