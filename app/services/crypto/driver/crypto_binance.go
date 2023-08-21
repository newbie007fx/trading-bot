package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"syscall"
	"telebot-trading/app/models"
	"telebot-trading/utils"
	"time"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceClient struct {
	binanceService *binance.Client
}

func (client *BinanceClient) init() {
	apiKey := utils.Env("BINANCE_API_KEY", "")
	secretKey := utils.Env("BINANCE_SECRET_KEY", "")
	client.binanceService = binance.NewClient(apiKey, secretKey)
}

func (client *BinanceClient) GetCandlesData(symbol string, limit int, startDate, endDate int64, resolution string) ([]models.CandleData, error) {
	var candlesData []models.CandleData
	conteks := context.Background()
	var cryptoCandles []*binance.Kline
	var err error

	count := 0
	for {
		count++
		conteks = context.Background()
		cryptoCandles, err = client.callGetCandleService(symbol, limit, startDate, endDate, resolution, conteks)
		if (err != nil && errors.Is(err, syscall.ECONNRESET)) || (conteks.Err() != nil && errors.Is(conteks.Err(), syscall.ECONNRESET)) {
			log.Print("This is connection reset by peer error")
			if count <= 3 {
				log.Println("sleep one second before retry")
				time.Sleep(1 * time.Second)
				continue
			}
		}

		break
	}

	if err == nil {
		candlesData = client.convertCandleDataMap(cryptoCandles)
	}

	if len(cryptoCandles) != limit {
		log.Println(symbol, " ,", limit, " ,", startDate, " ,", endDate, " ,", resolution)
		if conteks.Err() != nil {
			log.Println(conteks.Err().Error())
		}

		if err == nil {
			err = errors.New("invalid candle length")
		}

		s, _ := json.MarshalIndent(cryptoCandles, "", "\t")
		log.Print(string(s))
	}

	return candlesData, err
}

func (client *BinanceClient) callGetCandleService(symbol string, limit int, startDate, endDate int64, resolution string, conteks context.Context) (res []*binance.Kline, err error) {
	service := client.binanceService.NewKlinesService().Symbol(symbol).Limit(limit).Interval(resolution)
	if startDate > 1 {
		service = service.StartTime(startDate)
	}
	if endDate > 1 {
		service = service.EndTime(endDate)
	}
	return service.Do(conteks)
}

func (client *BinanceClient) GetBlanceInfo() (*[]models.AssetBalance, error) {
	assetBalances := []models.AssetBalance{}
	account, err := client.binanceService.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	balances := account.Balances
	for _, balance := range balances {
		s, _ := strconv.ParseFloat(balance.Free, 32)
		assetBalance := models.AssetBalance{AssetName: balance.Asset, Balance: float32(s)}
		assetBalances = append(assetBalances, assetBalance)
	}

	return &assetBalances, err
}

func (client *BinanceClient) GetExchangeInformation(symbols *[]string) (*[]models.MarketSymbol, error) {
	service := client.binanceService.NewExchangeInfoService()
	if symbols != nil {
		service = service.Symbols(*symbols...)
	} else {
		service = service.Permissions("SPOT")
	}

	res, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	return convertMarketSymbols(res.Symbols), nil
}

func (client *BinanceClient) ListPriceChangeStats(symbols *[]string) (*[]models.PriceChangeStats, error) {
	service := client.binanceService.NewListPriceChangeStatsService()
	service = service.Symbols(*symbols)

	res, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	return convertPriceChangeStats(res), nil
}

func (client *BinanceClient) CreateBuyOrder(symbol string, quantity float32) (*models.CreateOrderResponse, error) {
	return client.createOrder(binance.SideTypeBuy, symbol, quantity)
}

func (client *BinanceClient) CreateSellOrder(symbol string, quantity float32) (*models.CreateOrderResponse, error) {
	return client.createOrder(binance.SideTypeSell, symbol, quantity)
}

func (client *BinanceClient) createOrder(sideType binance.SideType, symbol string, quantity float32) (*models.CreateOrderResponse, error) {
	quoteQty := fmt.Sprintf("%f", quantity)
	orderService := client.binanceService.NewCreateOrderService().Symbol(symbol).Side(sideType).Type(binance.OrderTypeMarket).QuoteOrderQty(quoteQty)

	order, err := orderService.Do(context.Background())
	if err != nil {
		return nil, err
	}

	reponse := convertCreateOrderReponse(order)
	return &reponse, nil
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

func convertCreateOrderReponse(response *binance.CreateOrderResponse) models.CreateOrderResponse {
	var totalFillPrice float32 = 0
	for _, fill := range response.Fills {
		totalFillPrice += convertToFloat32(fill.Price)
	}

	return models.CreateOrderResponse{
		Symbol:   response.Symbol,
		Price:    totalFillPrice / float32(len(response.Fills)),
		Quantity: convertToFloat32(response.ExecutedQuantity),
		Status:   string(response.Status),
	}
}

func convertMarketSymbols(symbols []binance.Symbol) *[]models.MarketSymbol {
	marketSymbols := []models.MarketSymbol{}
	for _, symbol := range symbols {
		marketSymbol := models.MarketSymbol{
			Symbol:                     symbol.Symbol,
			Status:                     symbol.Status,
			BaseAsset:                  symbol.BaseAsset,
			BaseAssetPrecision:         symbol.BaseAssetPrecision,
			QuoteAsset:                 symbol.QuoteAsset,
			QuotePrecision:             symbol.QuotePrecision,
			QuoteAssetPrecision:        symbol.QuoteAssetPrecision,
			BaseCommissionPrecision:    symbol.BaseCommissionPrecision,
			QuoteCommissionPrecision:   symbol.QuoteCommissionPrecision,
			OrderTypes:                 symbol.OrderTypes,
			IcebergAllowed:             symbol.IcebergAllowed,
			OcoAllowed:                 symbol.OcoAllowed,
			QuoteOrderQtyMarketAllowed: symbol.QuoteOrderQtyMarketAllowed,
			IsSpotTradingAllowed:       symbol.IsSpotTradingAllowed,
			IsMarginTradingAllowed:     symbol.IsMarginTradingAllowed,
			Filters:                    symbol.Filters,
			Permissions:                symbol.Permissions,
		}
		marketSymbols = append(marketSymbols, marketSymbol)
	}

	return &marketSymbols
}

func convertPriceChangeStats(data []*binance.PriceChangeStats) *[]models.PriceChangeStats {
	priceChangeStats := []models.PriceChangeStats{}
	for _, priceChange := range data {
		stats := models.PriceChangeStats{
			Symbol:             priceChange.Symbol,
			PriceChange:        priceChange.PriceChange,
			PriceChangePercent: priceChange.PriceChangePercent,
			WeightedAvgPrice:   priceChange.WeightedAvgPrice,
			PrevClosePrice:     priceChange.PrevClosePrice,
			LastPrice:          priceChange.LastPrice,
			LastQty:            priceChange.LastQty,
			BidPrice:           priceChange.BidPrice,
			BidQty:             priceChange.BidQty,
			AskPrice:           priceChange.AskPrice,
			AskQty:             priceChange.AskQty,
			OpenPrice:          priceChange.OpenPrice,
			HighPrice:          priceChange.HighPrice,
			LowPrice:           priceChange.LowPrice,
			Volume:             priceChange.Volume,
			QuoteVolume:        priceChange.QuoteVolume,
			OpenTime:           priceChange.OpenTime,
			CloseTime:          priceChange.CloseTime,
			FristID:            priceChange.FristID,
			LastID:             priceChange.LastID,
			Count:              priceChange.Count,
		}
		priceChangeStats = append(priceChangeStats, stats)
	}

	return &priceChangeStats
}

func convertToFloat32(data string) float32 {
	f, err := strconv.ParseFloat(data, 32)
	if err != nil {
		fmt.Println(err.Error())
	}

	return float32(f)
}
