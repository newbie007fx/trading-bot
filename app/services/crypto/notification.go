package crypto

import (
	"log"
	"strconv"
	"telebot-trading/app/services"
	"telebot-trading/app/services/external"
)

func SendNotif(msg string) {
	if msg == "" {
		return
	}

	clintIDString := services.GetConfigValueByName("chat_id")
	if clintIDString == nil {
		log.Println("client id belum diset")
		return
	}

	clientID, _ := strconv.ParseInt(*clintIDString, 10, 64)

	err := external.SendToTelegram(clientID, msg)
	if err != nil {
		log.Println(err.Error())
	}
}
