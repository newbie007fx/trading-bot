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

	clintIDString := services.GetConfigValueByKey("chat_id")
	if clintIDString == nil {
		log.Println("client id belum diset")
		return
	}

	clientID, _ := strconv.ParseInt(*clintIDString, 10, 64)
	currency_notif_configs := repositories.GetCurrencyNotifConfigs()
	for _, data := range *currency_notif_configs {
		bands, err := services.GetCurrentBollingerBands(data.Symbol)
		if err != nil {
			log.Println("error: ", err.Error())
			continue
		}

		if !services.CheckLastCandleIsUp(bands) {
			continue
		}

		if services.CheckPositionInUpperBand(bands) {
			msg := fmt.Sprintf("Koin %s sekarang sudah melewati upper band, check dulu gan", data.Symbol)
			services.SendToTelegram(clientID, msg)
		}

		if services.CheckPositionSMAAfterLower(bands) {
			msg := fmt.Sprintf("Koin %s sekarang sudah melewati SMA band, Pantau gan", data.Symbol)
			services.SendToTelegram(clientID, msg)
		}

		if services.CheckPositionAfterLower(bands) {
			msg := fmt.Sprintf("Koin %s sekarang sudah melewati SMA band, Urgent gan", data.Symbol)
			services.SendToTelegram(clientID, msg)
		}
	}
}

func checkIsInSleepHours() bool {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(loc)

	sleepTimeStart, _ := time.ParseInLocation("15:04:05", "00:00:01", loc)
	sleepTimeEnd, _ := time.ParseInLocation("15:04:05", "00:04:30", loc)

	return sleepTimeStart.Unix() < now.Unix() && now.Unix() < sleepTimeEnd.Unix()
}
