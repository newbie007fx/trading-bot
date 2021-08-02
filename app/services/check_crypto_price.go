package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"telebot-trading/app/helper"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var currentTime int64

func StartCheckMasterCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		checkCryptoMasterCoinPrice()
	}
}

func StartCheckHoldCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		checkCryptoHoldCoinPrice()
	}
}

func StartCheckAltCoinPriceService(checkPriceChan chan bool) {
	for <-checkPriceChan {
		checkCryptoAltCoinPrice()
	}
}

func checkCryptoMasterCoinPrice() {
	log.Println("starting crypto check price master coin worker")

	startTime, endTime := getStartEndTime()

	responseChan := make(chan crypto.CandleResponse)

	masterCoinConfig, err := repositories.GetMasterCoinConfig()
	if err != nil {
		log.Println("error: ", err.Error())
		return
	}

	request := crypto.CandleRequest{
		Symbol:       masterCoinConfig.Symbol,
		Start:        startTime,
		End:          endTime,
		Resolution:   "15",
		ResponseChan: responseChan,
	}

	crypto.DispatchRequestJob(request)

	response := <-responseChan
	if response.Err != nil {
		log.Println("error: ", response.Err.Error(), " master")
		return
	}

	bands := analysis.GetCurrentBollingerBands(response.CandleData)
	result := buildResult(masterCoinConfig.Symbol, bands)

	if result.Direction == analysis.BAND_UP {
		result.Note = upTrendChecking(*masterCoinConfig, bands)
	} else {
		result.Note = downTrendChecking(*masterCoinConfig, bands)
	}

	setMasterCoin(result)

	log.Println("crypto check price worker is done")
}

func checkCryptoHoldCoinPrice() {
	log.Println("starting crypto check price hold coin worker")

	startTime, endTime := getStartEndTime()

	holdCoin := []models.BandResult{}

	responseChan := make(chan crypto.CandleResponse)

	condition := map[string]interface{}{"is_on_hold": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil)
	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Start:        startTime,
			End:          endTime,
			Resolution:   "15",
			ResponseChan: responseChan,
		}

		crypto.DispatchRequestJob(request)

		response := <-responseChan
		if response.Err != nil {
			log.Println("error: ", response.Err.Error(), " hold")
			continue
		}

		bands := analysis.GetCurrentBollingerBands(response.CandleData)
		result := buildResult(data.Symbol, bands)

		if result.Direction == analysis.BAND_UP {
			result.Note = upTrendChecking(data, bands)
		} else {
			result.Note = downTrendChecking(data, bands)
		}

		holdCoin = append(holdCoin, result)
	}

	msg := ""
	if len(holdCoin) > 0 {
		msg = "List coin yang dihold:\n"
		haveNote := false
		for _, coin := range holdCoin {
			if coin.Note != "" {
				msg += generateMsg(coin)
				msg += "\n"
				haveNote = true
			}
		}
		if !haveNote {
			msg = ""
		}

		masterCoin := getMasterCoin()
		if masterCoin != nil && msg != "" {
			msg += "untuk master coin:\n"
			msg += generateMsg(*masterCoin)
		}
	}

	sendNotif(msg)

	log.Println("crypto check price worker is done")
}

func checkCryptoAltCoinPrice() {
	log.Println("starting crypto check price for alt coin worker")

	startTime, endTime := getStartEndTime()

	altCoin := []models.BandResult{}

	responseChan := make(chan crypto.CandleResponse)

	limit := 85
	condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit)
	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Start:        startTime,
			End:          endTime,
			Resolution:   "15",
			ResponseChan: responseChan,
		}

		crypto.DispatchRequestJob(request)

		response := <-responseChan
		if response.Err != nil {
			log.Println("error: ", response.Err.Error())
			continue
		}

		bands := analysis.GetCurrentBollingerBands(response.CandleData)
		result := buildResult(data.Symbol, bands)

		if (result.Direction == analysis.BAND_UP || result.Trend == models.TREND_UP) && result.PriceChanges > 0.5 {
			result.Weight += getPositionWeight(result.Position)
			altCoin = append(altCoin, result)
		}

	}

	msg := ""
	if len(altCoin) > 0 {
		if len(altCoin) > 5 {
			altCoin = sortAndGetTopFive(altCoin)
		}

		msg += "top gain coin:\n"
		for _, coin := range altCoin {
			msg += generateMsg(coin)
			msg += "\n"
		}

		masterCoin := getMasterCoin()
		if masterCoin != nil {
			msg += "untuk master coin:\n"
			msg += generateMsg(*masterCoin)
		}
	}

	sendNotif(msg)

	log.Println("crypto check price worker is done")
}

func getStartEndTime() (int64, int64) {
	currentTime = time.Now().Unix()
	requestTime := currentTime
	if isTimeMultipleFifteenMinute(requestTime) {
		requestTime -= 1
	}
	startTime := requestTime - (60 * 15 * 29)

	return startTime, requestTime
}

func getMasterCoin() *models.BandResult {
	key := "last:master:coin:check"
	store := helper.GetSimpleStore()

	resultString := store.Get(key)
	if resultString != nil {
		masterCoin := models.BandResult{}
		err := json.Unmarshal([]byte(*resultString), &masterCoin)
		if err == nil {
			return &masterCoin
		}
	}

	return nil
}

func setMasterCoin(coin models.BandResult) error {
	key := "last:master:coin:check"
	store := helper.GetSimpleStore()

	result, err := json.Marshal(coin)
	if err == nil {
		resultString := string(result)
		store.Set(key, resultString)
		return nil
	}

	return err
}

func buildResult(symbol string, bands models.Bands) models.BandResult {

	direction := analysis.BAND_UP
	if !analysis.CheckLastCandleIsUp(bands.Data) {
		direction = analysis.BAND_DOWN
	}

	lastBand := bands.Data[len(bands.Data)-1]

	weight := bands.PriceChanges
	if bands.VolumeAverageChanges > 0 {
		weight += (bands.VolumeAverageChanges * 0.2 / 100)
	}

	result := models.BandResult{
		Symbol:        symbol,
		Direction:     direction,
		CurrentPrice:  lastBand.Candle.Close,
		CurrentVolume: lastBand.Candle.Volume,
		Trend:         bands.Trend,
		PriceChanges:  bands.PriceChanges,
		VolumeChanges: bands.VolumeAverageChanges,
		Weight:        weight,
		Position:      bands.Position,
	}

	return result
}

func sendNotif(msg string) {
	if msg == "" {
		return
	}

	clintIDString := GetConfigValueByName("chat_id")
	if clintIDString == nil {
		log.Println("client id belum diset")
		return
	}

	clientID, _ := strconv.ParseInt(*clintIDString, 10, 64)

	err := SendToTelegram(clientID, msg)
	if err != nil {
		log.Println(err.Error())
	}
}

func sortAndGetTopFive(coins []models.BandResult) []models.BandResult {
	sort.Slice(coins, func(i, j int) bool { return coins[i].Weight > coins[j].Weight })

	return coins[0:5]
}

func generateMsg(coinResult models.BandResult) string {
	format := "Coin name: <b>%s</b> \nDirection: <b>%s</b> \nPrice: <b>%f</b> \nVolume: <b>%f</b> \nTrend: <b>%s</b> \nPrice Changes: <b>%.2f%%</b> \nVolume Average Changes: <b>%.2f%%</b> \nNotes: <b>%s</b> \nPosition: <b>%s</b> \n"
	msg := fmt.Sprintf(format, coinResult.Symbol, directionString(coinResult.Direction), coinResult.CurrentPrice, coinResult.CurrentVolume, trendString(coinResult.Trend), coinResult.PriceChanges, coinResult.VolumeChanges, coinResult.Note, positionString(coinResult.Position))
	return msg
}

func upTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	if analysis.CheckPositionOnUpperBand(bands.Data) {
		return "Posisi naik upper band"
	}

	if analysis.CheckPositionSMAAfterLower(bands) {
		return "Posisi naik ke SMA"
	}

	if analysis.CheckPositionAfterLower(bands.Data) {
		return "Posisi lower"
	}

	if analysis.IsPriceIncreaseAboveThreshold(bands, data.IsMaster) {
		return "Naik diatas threshold"
	}

	if analysis.IsTrendUpAfterTrendDown(data.Symbol, bands) {
		return "Trend Up after down"
	}

	return ""
}

func downTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	if analysis.CheckPositionOnLowerBand(bands.Data) {
		return "Posisi turun dibawah lower"
	}

	if analysis.CheckPositionSMAAfterUpper(bands) {
		return "Posisi turun dibawah SMA"
	}

	if analysis.CheckPositionAfterUpper(bands.Data) {
		return "Posisi turun dari Upper"
	}

	if analysis.IsPriceDecreasebelowThreshold(bands, data.IsMaster) {
		return "Turun dibawah threshold"
	}

	if analysis.IsTrendDownAfterTrendUp(data.Symbol, bands) {
		return "Trend Down after up"
	}

	if data.IsOnHold && (bands.Position == models.ABOVE_SMA || bands.Position == models.ABOVE_UPPER) {
		if isTimeMultipleFifteenMinute(currentTime) {
			lastDown := countLastDownCandle(bands.Data)
			return fmt.Sprintf("Turun gan siaga !!! jumlah down %d", lastDown)
		}
	}

	return ""
}

func countLastDownCandle(data []models.Band) int {
	count := 0
	for i := len(data) - 1; i >= 0; i-- {
		band := data[i]
		if band.Candle.Close < band.Candle.Open {
			count++
		} else {
			break
		}
	}

	return count
}

func trendString(trend int8) string {
	if trend == models.TREND_UP {
		return "trend up"
	} else if trend == models.TREND_DOWN {
		return "trend down"
	}

	return "trend sideway"
}

func directionString(direction int8) string {
	if direction == analysis.BAND_UP {
		return "UP"
	}

	return "DOWN"
}

func positionString(position int8) string {
	if position == models.ABOVE_UPPER {
		return "above upper"
	} else if position == models.ABOVE_SMA {
		return "above sma"
	} else if position == models.BELOW_SMA {
		return "below sma"
	}

	return "below lower"
}

func isTimeMultipleFifteenMinute(currentTime int64) bool {
	fifteenMinutes := int64(60 * 15)

	return currentTime%fifteenMinutes == 0
}

func getPositionWeight(position int8) float32 {
	var weight float32 = 0
	if position == models.BELOW_SMA {
		weight = 0.5
	} else if position == models.ABOVE_SMA {
		weight = 0.25
	}

	return weight
}
