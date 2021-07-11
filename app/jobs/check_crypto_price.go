package jobs

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"time"
)

func CheckCryptoPrice() {
	if checkIsInSleepHours() {
		return
	}

	muted := isMuted()

	log.Println("starting crypto check price worker")

	counter := 0
	currentTime := time.Now().Unix()
	requestTime := currentTime
	if isTimeMultipleFifteenMinute(requestTime) {
		requestTime -= 1
	}

	altCoin := []models.BandResult{}
	holdCoin := []models.BandResult{}
	masterCoin := models.BandResult{}

	currency_configs := repositories.GetCurrencyNotifConfigs(nil)
	for _, data := range *currency_configs {
		if (muted || !isTimeMultipleFifteenMinute(currentTime)) && !data.IsMaster && !data.IsOnHold {
			continue
		}

		checkCounter(&counter, currentTime)

		bands, err := services.GetCurrentBollingerBands(data.Symbol, requestTime)
		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}

		direction := services.BAND_UP
		if !services.CheckLastCandleIsUp(bands.Data) {
			direction = services.BAND_DOWN
		}

		lastBand := bands.Data[len(bands.Data)-1]

		weight := bands.PriceChanges
		if bands.VolumeAverageChanges > 0 {
			weight += (bands.VolumeAverageChanges * 0.2 / 100)
		}

		result := models.BandResult{
			Symbol:        data.Symbol,
			Direction:     direction,
			CurrentPrice:  lastBand.Candle.Close,
			CurrentVolume: lastBand.Candle.Volume,
			Trend:         bands.Trend,
			PriceChanges:  bands.PriceChanges,
			VolumeChanges: bands.VolumeAverageChanges,
			Weight:        weight,
		}

		if data.IsMaster || data.IsOnHold {
			if result.Direction == services.BAND_UP {
				if muted {
					log.Println("On muted, exclude non urgent notif")
					return
				}

				result.Note = upTrendChecking(data, bands)
			} else {
				result.Note = downTrendChecking(data, bands)
			}

			if result.Note == "" && data.IsOnHold {
				continue
			}

			if data.IsMaster {
				masterCoin = result
			} else {
				holdCoin = append(holdCoin, result)
			}
		} else {
			if result.Direction == services.BAND_UP && result.PriceChanges > 1.0 {
				altCoin = append(altCoin, result)
			}
		}
	}

	sendNotif(masterCoin, holdCoin, altCoin)

	log.Println("crypto check price worker is done")
}

func sendNotif(masterCoin models.BandResult, holdCoin []models.BandResult, altCoin []models.BandResult) {
	clintIDString := services.GetConfigValueByName("chat_id")
	if clintIDString == nil {
		log.Println("client id belum diset")
		return
	}

	clientID, _ := strconv.ParseInt(*clintIDString, 10, 64)

	masterCoinMsg := "untuk master coin:\n"
	masterCoinMsg += generateMsg(masterCoin)

	if len(holdCoin) > 0 {
		msg := "List coin yang dihold:\n"
		for _, coin := range holdCoin {
			msg += generateMsg(coin)
			msg += "\n"
		}

		msg += masterCoinMsg

		err := services.SendToTelegram(clientID, msg)
		if err != nil {
			log.Println(err.Error())
		}
	}

	if len(altCoin) > 0 {
		if len(altCoin) > 5 {
			altCoin = sortAndGetTopFive(altCoin)
		}

		msg := "top gain coin:\n"
		for _, coin := range altCoin {
			msg += generateMsg(coin)
			msg += "\n"
		}

		msg += masterCoinMsg

		err := services.SendToTelegram(clientID, msg)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func sortAndGetTopFive(coins []models.BandResult) []models.BandResult {
	sort.Slice(coins, func(i, j int) bool { return coins[i].Weight > coins[j].Weight })

	return coins[0:5]
}

func generateMsg(coinResult models.BandResult) string {
	format := "Coin name: <b>%s</b> \nDirection: <b>%s</b> \nPrice: <b>%f</b> \nVolume: <b>%f</b> \nTrend: <b>%s</b> \nPrice Changes: <b>%.2f%%</b> \nVolume Average Changes: <b>%.2f%%</b> \nNotes: <b>%s</b> \n"
	msg := fmt.Sprintf(format, coinResult.Symbol, directionString(coinResult.Direction), coinResult.CurrentPrice, coinResult.CurrentVolume, trendString(coinResult.Trend), coinResult.PriceChanges, coinResult.VolumeChanges, coinResult.Note)
	return msg
}

func upTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	if services.CheckPositionOnUpperBand(bands.Data) {
		return "Posisi diupper band"
	}

	if services.CheckPositionSMAAfterLower(bands) {
		return "Posisi diSMA"
	}

	if services.CheckPositionAfterLower(bands.Data) {
		return "Posisi lower"
	}

	if services.IsPriceIncreaseAboveThreshold(bands, data.IsMaster) {
		return "Naik diatas threshold"
	}

	if services.IsTrendUpAfterTrendDown(data.Symbol, bands) {
		return "Trend Up after down"
	}

	return ""
}

func downTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	if services.CheckPositionOnLowerBand(bands.Data) {
		return "Posisi lower"
	}

	if services.CheckPositionSMAAfterUpper(bands) {
		return "Posisi SMA"
	}

	if services.CheckPositionAfterUpper(bands.Data) {
		return "Posisi Upper"
	}

	if services.IsPriceDecreasebelowThreshold(bands, data.IsMaster) {
		return "Turun dibawah threshold"
	}

	if services.IsTrendDownAfterTrendUp(data.Symbol, bands) {
		return "Trend Down after up"
	}

	return ""
}

func checkIsInSleepHours() bool {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(loc)

	date := now.Format("02-01-2006")

	sleepTimeStart, _ := time.ParseInLocation("02-01-2006 15:04:05", date+" 00:00:01", loc)
	sleepTimeEnd, _ := time.ParseInLocation("02-01-2006 15:04:05", date+" 04:30:00", loc)

	return sleepTimeStart.Unix() < now.Unix() && now.Unix() < sleepTimeEnd.Unix()
}

func trendString(trend int8) string {
	if trend == models.TREND_UP {
		return "trend up"
	} else if trend == models.TREND_DOWN {
		return "trend down"
	} else {
		return "trend sideway"
	}
}

func directionString(direction int8) string {
	if direction == services.BAND_UP {
		return "UP"
	} else {
		return "DOWN"
	}
}

func isTimeMultipleFifteenMinute(unixTime int64) bool {
	currentTime := time.Unix(unixTime, 0)
	minute := currentTime.Minute()
	if minute == 10 || minute == 25 || minute == 40 || minute == 55 {
		return true
	}

	return false
}

func isMuted() bool {
	key := "is-muted"
	resultString := repositories.GetConfigValueByName(key)
	if resultString != nil {
		tmp, err := strconv.ParseBool(*resultString)
		if err != nil {
			return false
		}
		return tmp
	}
	return false
}
