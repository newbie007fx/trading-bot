package driver

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"telebot-trading/app/models"
	"telebot-trading/utils"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceClient struct {
	klineService *binance.KlinesService
}

func (client *BinanceClient) init() {
	apiKey := utils.Env("BINANCE_API_KEY", "")
	secretKey := utils.Env("BINANCE_SECRET_KEY", "")
	service := binance.NewClient(apiKey, secretKey)
	client.klineService = service.NewKlinesService()
}

func (client *BinanceClient) GetCandlesData(symbol string, limit int, endDate int64, resolution string) ([]models.CandleData, error) {
	var candlesData []models.CandleData
	service := client.klineService.Symbol(formatSymbol(symbol)).Limit(limit).Interval(resolution)
	if endDate > 1 {
		service = service.EndTime(endDate)
	}
	cryptoCandles, err := service.Do(context.Background())
	if err == nil {
		candlesData = client.convertCandleDataMap(cryptoCandles)
	}

	return candlesData, err
}

func (BinanceClient) convertCandleDataMap(cryptoCanldes []*binance.Kline) []models.CandleData {
	candlesData := []models.CandleData{}

	for _, candle := range cryptoCanldes {
		candleData := models.CandleData{
			Open:      convertToFloat32(candle.Open),
			Close:     convertToFloat32(candle.Close),
			Low:       convertToFloat32(candle.Low),
			Hight:     convertToFloat32(candle.High),
			Volume:    convertToFloat32(candle.Volume),
			BuyVolume: convertToFloat32(candle.TakerBuyBaseAssetVolume),
			OpenTime:  candle.OpenTime,
			CloseTime: candle.CloseTime,
		}
		candlesData = append(candlesData, candleData)
	}

	return candlesData
}

func formatSymbol(symbol string) string {
	result := strings.Split(symbol, ":")
	if len(result) > 1 {
		return result[1]
	}

	return symbol
}

func convertToFloat32(data string) float32 {
	f, err := strconv.ParseFloat(data, 32)
	if err != nil {
		fmt.Println(err.Error())
	}

	return float32(f)
}
