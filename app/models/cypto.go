package models

const TREND_UP int8 = 1
const TREND_DOWN int8 = 2
const TREND_SIDEWAY int8 = 3

type CandleData struct {
	Open      float32
	Close     float32
	Low       float32
	Hight     float32
	Volume    float32
	Timestamp int64
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
}
