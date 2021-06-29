package jobs

import (
	"fmt"
	"log"
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

	clintIDString := services.GetConfigValueByName("chat_id")
	if clintIDString == nil {
		log.Println("client id belum diset")
		return
	}

	msg := ""

	clientID, _ := strconv.ParseInt(*clintIDString, 10, 64)
	currency_configs := repositories.GetCurrencyNotifConfigs()
	for _, data := range *currency_configs {
		bands, err := services.GetCurrentBollingerBands(data.Symbol)
		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}

		if !services.CheckLastCandleIsUp(bands.Data) {
			continue
		}

		if services.CheckPositionOnUpperBand(bands.Data) {
			msg += fmt.Sprintf("Koin %s sekarang sudah melewati upper band, check dulu gan (%s). ", data.Symbol, trendsString(bands.Trend))
		}

		if services.CheckPositionSMAAfterLower(bands.Data) {
			msg += fmt.Sprintf("Koin %s sekarang sudah melewati SMA band, Pantau gan (%s). ", data.Symbol, trendsString(bands.Trend))
		}

		if services.CheckPositionAfterLower(bands.Data) {
			msg += fmt.Sprintf("Koin %s sekarang sudah ijo setelah melewati lower band, Urgent gan (%s). ", data.Symbol, trendsString(bands.Trend))
		}
	}

	if msg != "" {
		services.SendToTelegram(clientID, msg)
	}

	log.Println("crypto check price worker is done")
}

func checkIsInSleepHours() bool {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(loc)

	date := now.Format("02-01-2006")

	sleepTimeStart, _ := time.ParseInLocation("02-01-2006 15:04:05", date+" 00:00:01", loc)
	sleepTimeEnd, _ := time.ParseInLocation("02-01-2006 15:04:05", date+" 04:30:00", loc)

	return sleepTimeStart.Unix() < now.Unix() && now.Unix() < sleepTimeEnd.Unix()
}

func trendsString(trend int8) string {
	if trend == models.TREND_UP {
		return "trend up"
	} else if trend == models.TREND_DOWN {
		return "trend down"
	} else {
		return "trend sideway"
	}
}
