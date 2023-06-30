package models

const TREND_UP int8 = 1
const TREND_DOWN int8 = 2
const TREND_SIDEWAY int8 = 3

const BELOW_LOWER int8 = 1
const BELOW_SMA int8 = 2
const ABOVE_SMA int8 = 3
const ABOVE_UPPER int8 = 4

type CandleData struct {
	Open      float32
	Close     float32
	Low       float32
	Hight     float32
	BuyVolume float32
	Volume    float32
	OpenTime  int64
	CloseTime int64
}

type Band struct {
	Candle *CandleData
	SMA    float64
	Upper  float64
	Lower  float64
}

type Bands struct {
	Data          []Band
	AllTrend      TrendDetail
	PriceChanges  float32
	Position      int8
	HeuristicBand Heuristic
}

type Heuristic struct {
	FirstBand  Band
	SecondBand Band
	ThirdBand  Band
	FourthBand Band
}

type BandResult struct {
	Symbol        string
	Direction     int8
	CurrentPrice  float32
	CurrentVolume float32
	AllTrend      TrendDetail
	PriceChanges  float32
	VolumeChanges float32
	Weight        float32
	Position      int8
	Bands         []Band
	Mid           *BandResult
	Long          *BandResult
	HeuristicBand Heuristic
}

type AssetBalance struct {
	AssetName string
	Balance   float32
}

type CreateOrderResponse struct {
	Symbol   string
	Price    float32
	Quantity float32
	Status   string
}

type TrendDetail struct {
	FirstTrend         int8
	FirstTrendPercent  float32
	SecondTrend        int8
	SecondTrendPercent float32
	Trend              int8
	ShortTrend         int8
}

type MarketSymbol struct {
	Symbol                     string
	Status                     string
	BaseAsset                  string
	BaseAssetPrecision         int
	QuoteAsset                 string
	QuotePrecision             int
	QuoteAssetPrecision        int
	BaseCommissionPrecision    int32
	QuoteCommissionPrecision   int32
	OrderTypes                 []string
	IcebergAllowed             bool
	OcoAllowed                 bool
	QuoteOrderQtyMarketAllowed bool
	IsSpotTradingAllowed       bool
	IsMarginTradingAllowed     bool
	Filters                    []map[string]interface{}
	Permissions                []string
}

type PriceChangeStats struct {
	Symbol             string
	PriceChange        string
	PriceChangePercent string
	WeightedAvgPrice   string
	PrevClosePrice     string
	LastPrice          string
	LastQty            string
	BidPrice           string
	BidQty             string
	AskPrice           string
	AskQty             string
	OpenPrice          string
	HighPrice          string
	LowPrice           string
	Volume             string
	QuoteVolume        string
	OpenTime           int64
	CloseTime          int64
	FristID            int64
	LastID             int64
	Count              int64
}
