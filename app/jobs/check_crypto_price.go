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

	log.Println("starting crypto check price worker")

	counter := 0
	currentTime := time.Now().Unix()

	altCoin := []models.BandResult{}
	holdCoin := []models.BandResult{}
	masterCoin := models.BandResult{}

	currency_configs := repositories.GetCurrencyNotifConfigs()
	for _, data := range *currency_configs {
		if !isTimeMultipleFifteenMinute(currentTime) && !data.IsMaster && !data.IsOnHold {
			continue
		}

		checkCounter(&counter, currentTime)

		bands, err := services.GetCurrentBollingerBands(data.Symbol)
		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}

		direction := services.BAND_UP
		if !services.CheckLastCandleIsUp(bands.Data) {
			direction = services.BAND_DOWN
		}

		lastBand := bands.Data[len(bands.Data)-1]
		result := models.BandResult{
			Symbol:        data.Symbol,
			Direction:     direction,
			CurrentPrice:  lastBand.Candle.Close,
			CurrentVolume: lastBand.Candle.Volume,
			Trend:         bands.Trend,
			PriceChanges:  bands.PriceChanges,
			VolumeChanges: bands.VolumeAverageChanges,
			Weight:        bands.PriceChanges + (bands.VolumeAverageChanges * 20 / 100),
		}

		if data.IsMaster || data.IsOnHold {
			if result.Direction == services.BAND_UP {
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
			if result.Direction == services.BAND_UP {
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

	if len(holdCoin) > 0 {
		msg := "List coin yang dihold</br>"
		for _, coin := range holdCoin {
			msg += generateMsg(coin)
			msg += "</br>"
		}
		msg += "untuk master coin"
		msg += generateMsg(masterCoin)

		services.SendToTelegram(clientID, msg)
	}

	if len(altCoin) > 0 {
		if len(altCoin) > 5 {
			altCoin = sortAndGetTopFive(altCoin)
		}

		msg := "top gain coin</br>"
		for _, coin := range altCoin {
			msg += generateMsg(coin)
			msg += "</br>"
		}

		services.SendToTelegram(clientID, msg)
	}
}

func sortAndGetTopFive(coins []models.BandResult) []models.BandResult {
	sort.Slice(coins, func(i, j int) bool { return coins[i].Weight > coins[j].Weight })

	return coins[0:5]
}

func generateMsg(coinResult models.BandResult) string {
	format := "Coin name: %s</br>Direction: %s</br>Price: %.2f</br>Volume: %.2f</br>Trend: %s</br>Price Changes: %.2f%%</br>Volume Average Changes: %.2f%%</br>Notes: %s</br>"
	msg := fmt.Sprintf(format, coinResult.Symbol, directionString(coinResult.Direction), coinResult.CurrentPrice, coinResult.CurrentVolume, trendString(coinResult.Trend), coinResult.PriceChanges, coinResult.VolumeChanges, coinResult.Note)
	return msg
}

func upTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	note := ""
	if services.CheckPositionOnUpperBand(bands.Data) {
		note += "Posisi diupper band"
	}

	if services.CheckPositionSMAAfterLower(bands) {
		note += "Posisi diSMA"
	}

	if services.CheckPositionAfterLower(bands.Data) {
		note += "Posisi lower"
	}

	if services.IsPriceIncreaseAboveThreshold(bands, data.IsMaster) {
		note += "Naik diatas threshold"
	}

	return note
}

func downTrendChecking(data models.CurrencyNotifConfig, bands models.Bands) string {
	note := ""
	if services.CheckPositionOnLowerBand(bands.Data) {
		note += "Posisi lower"
	}

	if services.CheckPositionSMAAfterUpper(bands) {
		note += "Posisi SMA"
	}

	if services.CheckPositionAfterUpper(bands.Data) {
		note += "Posisi Upper"
	}

	if services.IsPriceDecreasebelowThreshold(bands, data.IsMaster) {
		note += "Turun dibawah threshold"
	}

	return note
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

func isTimeMultipleFifteenMinute(currentTime int64) bool {
	fifteenMinutes := int64(60 * 15)

	return currentTime%fifteenMinutes == 0
}

func checkCounter(counter *int, startTime int64) {
	currentTime := time.Now().Unix()
	difference := currentTime - startTime
	sleep := 0
	quotaLeft := difference % 60
	if *counter == 60 {
		if quotaLeft > 0 {
			sleep = int(quotaLeft) + 1
		}
		*counter = 0
	}

	if difference%60 == 0 {
		*counter = 0
		sleep = 1
	}

	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	*counter++
}
