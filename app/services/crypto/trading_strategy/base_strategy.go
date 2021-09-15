package trading_strategy

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

type TradingStrategy interface {
	Execute(currentTime time.Time)
	InitService()
	Shutdown()
}

var masterCoin *models.BandResult
var masterCoinLongInterval *models.BandResult
var waitMasterCoin bool

func StartCheckMasterCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		checkCryptoMasterCoinPrice()
	}
}

func checkCryptoMasterCoinPrice() {
	waitMasterCoin = true

	log.Println("starting crypto check price master coin worker")

	responseChan := make(chan crypto.CandleResponse)

	masterCoinConfig, err := repositories.GetMasterCoinConfig()
	if err != nil {
		log.Println("error: ", err.Error())
		waitMasterCoin = false
		return
	}

	request := crypto.CandleRequest{
		Symbol:       masterCoinConfig.Symbol,
		EndDate:      GetEndDate(nil),
		Limit:        40,
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := crypto.MakeCryptoRequest(*masterCoinConfig, request)

	masterCoin = result

	request = crypto.CandleRequest{
		Symbol:       masterCoinConfig.Symbol,
		EndDate:      GetEndDate(nil),
		Limit:        40,
		Resolution:   "1h",
		ResponseChan: responseChan,
	}

	result = crypto.MakeCryptoRequest(*masterCoinConfig, request)

	masterCoinLongInterval = result

	log.Println("crypto check price worker is done")
	waitMasterCoin = false
}

func checkCryptoHoldCoinPrice() []models.BandResult {
	log.Println("starting crypto check price hold coin worker")

	holdCoin := []models.BandResult{}

	endDate := GetEndDate(nil)

	responseChan := make(chan crypto.CandleResponse)

	condition := map[string]interface{}{"is_on_hold": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil)

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Limit:        40,
			EndDate:      endDate,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := crypto.MakeCryptoRequest(data, request)
		if result == nil {
			continue
		}

		holdCoin = append(holdCoin, *result)
	}

	log.Println("crypto check price worker is done")

	return holdCoin
}

func checkCryptoAltCoinPrice(baseTime *time.Time) []models.BandResult {
	log.Println("starting crypto check price for alt coin worker")

	altCoin := []models.BandResult{}

	endDate := GetEndDate(baseTime)

	responseChan := make(chan crypto.CandleResponse)

	limit := 120
	condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit)

	waitMasterCoinProcessed()
	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			EndDate:      endDate,
			Limit:        40,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := crypto.MakeCryptoRequest(data, request)
		if result == nil || result.Direction == analysis.BAND_DOWN {
			continue
		}

		result.Weight = analysis.CalculateWeight(result, *masterCoin)
		if !analysis.IsIgnored(result) && result.Weight > 1.49 {
			altCoin = append(altCoin, *result)
		}

		//log.Println(analysis.GetWeightLog())
	}

	log.Println("crypto check price worker is done")

	return altCoin
}

func GetEndDate(baseTime *time.Time) int64 {
	var endTime time.Time
	if baseTime == nil {
		endTime = time.Now()
	} else {
		endTime = *baseTime
	}
	unixTime := endTime.Unix()

	return unixTime * 1000
}

func waitMasterCoinProcessed() {
	for waitMasterCoin {
		time.Sleep(1 * time.Second)
	}
}
