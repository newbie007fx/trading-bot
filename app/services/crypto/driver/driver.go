package driver

import (
	"telebot-trading/app/models"
)

type CryptoDriver interface {
	init()
	GetCandlesData(symbol string, limit int, endDate int64, resolution string) ([]models.CandleData, error)
	GetBlanceInfo() (*[]models.AssetBalance, error)
	CreateBuyOrder(symbol string, quantity float32) (*models.CreateOrderResponse, error)
	CreateSellOrder(symbol string, quantity float32) (*models.CreateOrderResponse, error)
}

var crypto CryptoDriver

func GetCrypto() CryptoDriver {
	if crypto == nil {
		crypto = new(BinanceClient)
		crypto.init()
	}

	return crypto
}
