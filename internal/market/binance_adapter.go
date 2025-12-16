package market

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/newbie007fx/trading-bot/internal/model"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceAdapter struct {
	client *binance.Client
}

func NewBinanceAdapter(binanceKey, binanceSecret string) *BinanceAdapter {
	return &BinanceAdapter{
		client: binance.NewClient(binanceKey, binanceSecret), // public data only
	}
}

func (b *BinanceAdapter) GetCandles(
	ctx context.Context,
	symbol string,
	interval string,
	limit int,
	endTime *int64,
) ([]model.CandleData, error) {
	klinesService := b.client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit)

	if endTime != nil {
		klinesService.EndTime(*endTime)
	}

	klines, err := klinesService.Do(ctx)

	if err != nil {
		log.Printf("[BINANCE ERROR] %v", err)
		return nil, err
	}

	candles := make([]model.CandleData, 0, len(klines))
	for _, k := range klines {
		candles = append(candles, model.CandleData{
			OpenTime:  k.OpenTime,
			Open:      convertToFloat64(k.Open),
			Close:     convertToFloat64(k.Close),
			Low:       convertToFloat64(k.Low),
			Hight:     convertToFloat64(k.High),
			Volume:    convertToFloat64(k.Volume),
			BuyVolume: convertToFloat64(k.TakerBuyBaseAssetVolume),
			CloseTime: k.CloseTime,
		})
	}

	return candles, nil
}

func convertToFloat64(data string) float64 {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	return f
}
