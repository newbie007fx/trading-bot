package market

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/newbie007fx/trading-bot/internal/config"
	"github.com/newbie007fx/trading-bot/internal/infra/secret"
	"github.com/newbie007fx/trading-bot/internal/model"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceAdapter struct {
	client *binance.Client
}

func NewBinanceAdapter(ctx context.Context, cfg config.Config) *BinanceAdapter {
	secretLoader, err := secret.NewLoader(ctx, cfg.ProjectID, cfg.Location)
	if err != nil {
		log.Fatal(err)
	}
	defer secretLoader.Close()

	binanceKey, err := secretLoader.Get(ctx, "BINANCE_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	binanceSecret, err := secretLoader.Get(ctx, "BINANCE_SECRET_KEY")
	if err != nil {
		log.Fatal(err)
	}

	return &BinanceAdapter{
		client: binance.NewClient(binanceKey, binanceSecret), // public data only
	}
}

func (b *BinanceAdapter) GetCandles(
	ctx context.Context,
	symbol string,
	interval string,
	limit int,
) ([]model.CandleData, error) {

	log.Printf(
		"[BINANCE] symbol=%s interval=%s limit=%d",
		symbol, interval, limit,
	)

	klines, err := b.client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(ctx)

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
