package driver

import (
	"telebot-trading/app/models"
)

type CryptoDriver interface {
	init()
	GetCandlesData(symbol string, limit int, startDate, endDate int64, resolution string) ([]models.CandleData, error)
	GetBlanceInfo() (*[]models.AssetBalance, error)
	GetExchangeInformation() (*[]models.MarketSymbol, error)
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
