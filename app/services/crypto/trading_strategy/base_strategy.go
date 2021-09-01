package trading_strategy

import (
	"log"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
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
		return
	}

	request := crypto.CandleRequest{
		Symbol:       masterCoinConfig.Symbol,
		EndDate:      GetEndDate(),
		Limit:        33,
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := crypto.MakeCryptoRequest(*masterCoinConfig, request)
	if result == nil {
		return
	}

	masterCoin = result

	log.Println("crypto check price worker is done")
	waitMasterCoin = false
}

func checkCryptoHoldCoinPrice() []models.BandResult {
	log.Println("starting crypto check price hold coin worker")

	holdCoin := []models.BandResult{}

	endDate := GetEndDate()

	responseChan := make(chan crypto.CandleResponse)

	condition := map[string]interface{}{"is_on_hold": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil)

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Limit:        33,
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

func checkCryptoAltCoinPrice() []models.BandResult {
	log.Println("starting crypto check price for alt coin worker")

	altCoin := []models.BandResult{}

	endDate := GetEndDate()

	responseChan := make(chan crypto.CandleResponse)

	limit := 120
	condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit)

	var masterCoinTrend int8 = 0
	waitMasterCoinProcessed()
	if masterCoin != nil {
		masterCoinTrend = masterCoin.Trend
	}

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			EndDate:      endDate,
			Limit:        33,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := crypto.MakeCryptoRequest(data, request)
		if result == nil {
			continue
		}

		result.Weight = analysis.CalculateWeight(result, masterCoinTrend)
		if !analysis.IsIgnored(result) && result.Direction == analysis.BAND_UP && result.Weight > 1.79 {
			altCoin = append(altCoin, *result)
		}

		//log.Println(analysis.GetWeightLog())
	}

	log.Println("crypto check price worker is done")

	return altCoin
}

func GetEndDate() int64 {
	currentTime := time.Now()
	minute := currentTime.Minute()
	timeInMili := currentTime.Unix() * 1000
	if minute%15 == 0 {
		timeInMili -= 1
	}

	return timeInMili
}

func waitMasterCoinProcessed() {
	for waitMasterCoin {
		time.Sleep(1 * time.Second)
	}
}

func sendNotif(msg string) {
	if msg == "" {
		return
	}

	clintIDString := services.GetConfigValueByName("chat_id")
	if clintIDString == nil {
		log.Println("client id belum diset")
		return
	}

	clientID, _ := strconv.ParseInt(*clintIDString, 10, 64)

	err := services.SendToTelegram(clientID, msg)
	if err != nil {
		log.Println(err.Error())
	}
}
