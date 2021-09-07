package crypto

import (
	"log"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
)

func SyncBalance() {
	cryptoDriver := driver.GetCrypto()

	balances, err := cryptoDriver.GetBlanceInfo()
	if err != nil {
		log.Println("error with message: ", err.Error())
	}

	for _, balance := range *balances {
		if balance.AssetName == "USDT" {
			SetBalance(balance.Balance)
		} else {
			symbol := balance.AssetName + "USDT"
			repositories.UpdateCurrencyNotifConfigBySymbol(symbol, map[string]interface{}{"balance": balance.Balance})
		}
	}
}
