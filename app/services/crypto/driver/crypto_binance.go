package driver

import (
	"context"
	"fmt"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/utils"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceClient struct {
	klineService       *binance.KlinesService
	accountService     *binance.GetAccountService
	createOrderService *binance.CreateOrderService
}

func (client *BinanceClient) init() {
	apiKey := utils.Env("BINANCE_API_KEY", "")
	secretKey := utils.Env("BINANCE_SECRET_KEY", "")
	service := binance.NewClient(apiKey, secretKey)
	client.klineService = service.NewKlinesService()
	client.accountService = service.NewGetAccountService()
	client.createOrderService = service.NewCreateOrderService()
}

func (client *BinanceClient) GetCandlesData(symbol string, limit int, endDate int64, resolution string) ([]models.CandleData, error) {
	var candlesData []models.CandleData
	service := client.klineService.Symbol(symbol).Limit(limit).Interval(resolution)
	if endDate > 1 {
		service = service.EndTime(endDate)
	}
	cryptoCandles, err := service.Do(context.Background())
	if err == nil {
		candlesData = client.convertCandleDataMap(cryptoCandles)
	}

	return candlesData, err
}

func (client *BinanceClient) GetBlanceInfo() (*[]models.AssetBalance, error) {
	assetBalances := []models.AssetBalance{}
	account, err := client.accountService.Do(context.Background())
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

func (client *BinanceClient) CreateBuyOrder(symbol string, quantity float32) (*models.CreateOrderResponse, error) {
	return client.createOrder(binance.SideTypeBuy, symbol, quantity)
}

func (client *BinanceClient) CreateSellOrder(symbol string, quantity float32) (*models.CreateOrderResponse, error) {
	return client.createOrder(binance.SideTypeSell, symbol, quantity)
}

func (client *BinanceClient) createOrder(sideType binance.SideType, symbol string, quantity float32) (*models.CreateOrderResponse, error) {
	quoteQty := fmt.Sprintf("%f", quantity)
	orderService := client.createOrderService.Symbol(symbol).Side(sideType).Type(binance.OrderTypeMarket).QuoteOrderQty(quoteQty)

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

func convertToFloat32(data string) float32 {
	f, err := strconv.ParseFloat(data, 32)
	if err != nil {
		fmt.Println(err.Error())
	}

	return float32(f)
}
