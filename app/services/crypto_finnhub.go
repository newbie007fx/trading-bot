package services

import (
	"context"
	"telebot-trading/utils"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type CandleData struct {
	Open      float32
	Close     float32
	Low       float32
	Hight     float32
	Volume    float32
	Timestamp int64
}

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

func (client *FinnhubClient) GetCandlesData(symbol string, startTime, endTime int64) ([]CandleData, error) {
	var candlesData []CandleData
	cryptoCandles, _, err := client.service.CryptoCandles(client.contextAuth, symbol, "15", startTime, endTime)
	if err == nil && cryptoCandles.S == "ok" {
		candlesData = client.convertCandleDataMap(cryptoCandles)
	}

	return candlesData, err
}

func (FinnhubClient) convertCandleDataMap(cryptoCanldes finnhub.CryptoCandles) []CandleData {

	candlesData := []CandleData{}
	size := len(cryptoCanldes.O)

	for i := 0; i < size; i++ {
		candleData := CandleData{
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
