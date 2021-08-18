package driver

import "telebot-trading/app/models"

type CryptoDriver interface {
	init()
	GetCandlesData(symbol string, limit int, endDate int64, resolution string) ([]models.CandleData, error)
}

var crypto CryptoDriver

func GetCrypto() CryptoDriver {
	if crypto == nil {
		crypto = new(BinanceClient)
		crypto.init()
	}

	return crypto
}
