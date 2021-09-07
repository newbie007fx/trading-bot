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
	Data                 []Band
	Trend                int8
	PriceChanges         float32
	VolumeAverageChanges float32
	Position             int8
}

type BandResult struct {
	Symbol        string
	Direction     int8
	CurrentPrice  float32
	CurrentVolume float32
	Trend         int8
	PriceChanges  float32
	VolumeChanges float32
	Weight        float32
	Note          string
	Position      int8
	Bands         []Band
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
