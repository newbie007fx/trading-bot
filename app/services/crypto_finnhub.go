package services

import (
	"context"
	"errors"
	"telebot-trading/app/models"
	"telebot-trading/utils"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type FinnhubClient struct {
	service     *finnhub.DefaultApiService
	contextAuth context.Context
}

func (client *FinnhubClient) init() {
	api_key := utils.Env("FINNHUB_KEY", "")
	client.service = finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	client.contextAuth = context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: api_key,
	})
}

func (client *FinnhubClient) GetCandlesData(symbol string, startTime, endTime int64) ([]models.CandleData, error) {
	var candlesData []models.CandleData
	cryptoCandles, _, err := client.service.CryptoCandles(client.contextAuth, symbol, "15", startTime, endTime-1)
	if err == nil && cryptoCandles.S == "ok" {
		candlesData = client.convertCandleDataMap(cryptoCandles)
	} else if err == nil && cryptoCandles.S != "ok" {
		err = errors.New(cryptoCandles.S)
	}

	return candlesData, err
}

func (FinnhubClient) convertCandleDataMap(cryptoCanldes finnhub.CryptoCandles) []models.CandleData {

	candlesData := []models.CandleData{}
	size := len(cryptoCanldes.O)

	for i := 0; i < size; i++ {
		candleData := models.CandleData{
			Open:      cryptoCanldes.O[i],
			Close:     cryptoCanldes.C[i],
			Low:       cryptoCanldes.L[i],
			Hight:     cryptoCanldes.H[i],
			Volume:    cryptoCanldes.V[i],
			Timestamp: cryptoCanldes.T[i],
		}
		candlesData = append(candlesData, candleData)
	}

	return candlesData
}

var crypto *FinnhubClient

func GetCrypto() *FinnhubClient {
	if crypto == nil {
		crypto = new(FinnhubClient)
		crypto.init()
	}

	return crypto
}
