package jobs

import (
	"fmt"
	"log"
	"strconv"
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

		if !services.CheckLastCandleIsUp(bands) {
			continue
		}

		if services.CheckPositionOnUpperBand(bands) {
			msg += fmt.Sprintf("Koin %s sekarang sudah melewati upper band, check dulu gan. ", data.Symbol)
		}

		if services.CheckPositionSMAAfterLower(bands) {
			msg += fmt.Sprintf("Koin %s sekarang sudah melewati SMA band, Pantau gan. ", data.Symbol)
		}

		if services.CheckPositionAfterLower(bands) {
			msg += fmt.Sprintf("Koin %s sekarang sudah ijo setelah melewati lower band, Urgent gan. ", data.Symbol)
		}
	}

	services.SendToTelegram(clientID, msg)

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
