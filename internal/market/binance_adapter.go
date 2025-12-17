package market

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/newbie007fx/trading-bot/internal/model"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceAdapter struct {
	client      *binance.Client
	symbol      string
	minNotional float64
	stepSize    float64
	minQty      float64
}

func NewBinanceAdapter(ctx context.Context, binanceKey, binanceSecret string, symbol string) *BinanceAdapter {
	binanceAdapter := &BinanceAdapter{
		client: binance.NewClient(binanceKey, binanceSecret), // public data only
		symbol: symbol,
	}

	binanceAdapter.loadSymbolFilter(ctx)

	return binanceAdapter
}

func (b *BinanceAdapter) GetCandles(
	ctx context.Context,
	interval string,
	limit int,
	endTime *int64,
) ([]model.CandleData, error) {
	klinesService := b.client.NewKlinesService().
		Symbol(b.symbol).
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

func (b *BinanceAdapter) Buy(ctx context.Context, price float64) error {
	usdtBalance, err := b.GetFreeBalance(ctx, "USDT")
	if err != nil {
		return err
	}

	const feeBuffer = 0.003 // 0.3%

	maxSpend := usdtBalance * (1 - feeBuffer)
	orderUSDT := math.Floor(maxSpend*100) / 100 // FLOOR

	if orderUSDT < b.minNotional {
		return errors.New("insufficient USDT balance")
	}

	log.Printf(
		"[LIVE] BUY %s orderUSDT=%.2f",
		b.symbol, orderUSDT,
	)

	order, err := b.client.NewCreateOrderService().
		Symbol(b.symbol).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		QuoteOrderQty(fmt.Sprintf("%.2f", orderUSDT)).
		Do(ctx)

	if err != nil {
		return err
	}

	log.Printf(
		"[BINANCE LIVE] BUY FILLED orderId=%d status=%s",
		order.OrderID,
		order.Status,
	)

	return nil
}

func (b *BinanceAdapter) Sell(ctx context.Context) error {
	ethBalance, err := b.GetFreeBalance(ctx, "ETH")
	if err != nil {
		return err
	}

	if ethBalance <= 0 {
		return errors.New("no ETH to sell")
	}

	const feeBuffer = 0.003 // 0.3%
	sellable := ethBalance * (1 - feeBuffer)

	sellQty := adjustToStepSize(sellable, b.stepSize)

	if sellQty < b.minQty {
		return errors.New("sell qty below minQty (dust)")
	}

	avgPrice, err := b.client.NewAveragePriceService().
		Symbol(b.symbol).
		Do(ctx)
	if err != nil {
		return err
	}

	averagePrice, _ := strconv.ParseFloat(avgPrice.Price, 64)
	notional := sellQty * averagePrice
	if notional < b.minNotional {
		return errors.New("sell notional below minNotional")
	}

	log.Printf(
		"[LIVE] SELL %s qty=%.6f balance=%.6f",
		b.symbol, sellQty, ethBalance,
	)

	order, err := b.client.NewCreateOrderService().
		Symbol(b.symbol).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.6f", sellQty)).
		Do(ctx)

	if err != nil {
		return err
	}

	log.Printf(
		"[BINANCE LIVE] SELL FILLED orderId=%d status=%s",
		order.OrderID,
		order.Status,
	)

	return nil
}

func (b *BinanceAdapter) GetFreeBalance(ctx context.Context, asset string) (float64, error) {

	acc, err := b.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return 0, err
	}

	for _, bal := range acc.Balances {
		if bal.Asset == asset {
			return strconv.ParseFloat(bal.Free, 64)
		}
	}

	return 0, errors.New("asset not found")
}

func (b *BinanceAdapter) loadSymbolFilter(ctx context.Context) {

	info, err := b.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range info.Symbols {
		if s.Symbol == b.symbol {
			for _, f := range s.Filters {
				switch f["filterType"] {
				case "LOT_SIZE":
					b.stepSize, _ = strconv.ParseFloat(f["stepSize"].(string), 64)
					b.minQty, _ = strconv.ParseFloat(f["minQty"].(string), 64)
				case "NOTIONAL":
					b.minNotional, _ = strconv.ParseFloat(f["minNotional"].(string), 64)
				}
			}
			break
		}
	}
}

func convertToFloat64(data string) float64 {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	return f
}

func adjustToStepSize(qty float64, step float64) float64 {
	return math.Floor(qty/step) * step
}
