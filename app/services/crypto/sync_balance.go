package crypto

import (
	"log"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
)

var syncBalanceChan chan bool

func StartSyncBalanceService() {
	syncBalanceChan = make(chan bool, 10)
	for <-syncBalanceChan {
		SyncBalance()
	}
}

func RequestSyncBalance() {
	syncBalanceChan <- true
}

func SyncBalance() {
	cryptoDriver := driver.GetCrypto()

	balances, err := cryptoDriver.GetBlanceInfo()
	if err != nil {
		log.Println("error with message: ", err.Error())
		return
	}

	for _, balance := range *balances {
		if balance.Balance == 0 {
			continue
		}

		if balance.AssetName == "USDT" {
			SetBalance(balance.Balance)
		} else {
			symbol := balance.AssetName + "USDT"
			repositories.UpdateCurrencyNotifConfigBySymbol(symbol, map[string]interface{}{"balance": balance.Balance})
		}
	}
}
