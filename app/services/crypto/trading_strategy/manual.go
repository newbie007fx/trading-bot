package trading_strategy

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var masterCoin *models.BandResult
var waitMasterCoin bool

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
		Limit:        29,
		Resolution:   "15m",
		ResponseChan: responseChan,
	}

	result := services.MakeCryptoRequest(*masterCoinConfig, request)
	if result == nil {
		return
	}

	masterCoin = result

	log.Println("crypto check price worker is done")
	waitMasterCoin = false
}

func checkCryptoHoldCoinPrice() {
	log.Println("starting crypto check price hold coin worker")

	holdCoin := []models.BandResult{}

	responseChan := make(chan crypto.CandleResponse)

	condition := map[string]interface{}{"is_on_hold": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil)

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Limit:        29,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := services.MakeCryptoRequest(data, request)
		if result == nil {
			continue
		}

		holdCoin = append(holdCoin, *result)
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

		waitMasterCoinProcessed()
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

	altCoin := []models.BandResult{}

	responseChan := make(chan crypto.CandleResponse)

	limit := 82
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
			Limit:        29,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := services.MakeCryptoRequest(data, request)
		if result == nil {
			continue
		}

		result.Weight = analysis.CalculateWeight(result, masterCoinTrend)
		if result.Direction == analysis.BAND_UP && result.Weight > 1.5 {
			altCoin = append(altCoin, *result)
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

		if masterCoin != nil {
			msg += "untuk master coin:\n"
			msg += generateMsg(*masterCoin)
		}
	}

	sendNotif(msg)

	log.Println("crypto check price worker is done")
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

func sortAndGetTopFive(coins []models.BandResult) []models.BandResult {
	sort.Slice(coins, func(i, j int) bool { return coins[i].Weight > coins[j].Weight })

	return coins[0:5]
}

func generateMsg(coinResult models.BandResult) string {
	format := "Coin name: <b>%s</b> \nDirection: <b>%s</b> \nPrice: <b>%f</b> \nVolume: <b>%f</b> \nTrend: <b>%s</b> \nPrice Changes: <b>%.2f%%</b> \nVolume Average Changes: <b>%.2f%%</b> \nNotes: <b>%s</b> \nPosition: <b>%s</b> \n"
	msg := fmt.Sprintf(format, coinResult.Symbol, directionString(coinResult.Direction), coinResult.CurrentPrice, coinResult.CurrentVolume, trendString(coinResult.Trend), coinResult.PriceChanges, coinResult.VolumeChanges, coinResult.Note, positionString(coinResult.Position))
	return msg
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

func waitMasterCoinProcessed() {
	for waitMasterCoin {
		time.Sleep(1 * time.Second)
	}
}
