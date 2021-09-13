package crypto

import (
	"log"
	"strconv"
	"sync"
	"telebot-trading/app/services"
	"telebot-trading/app/services/external"
	"time"
)

var mutex *sync.Mutex
var listNotifData []notifData

type notifData struct {
	msg           string
	timeToExecute int64
}

func StartSendNotifService() {
	mutex = &sync.Mutex{}
	listNotifData = []notifData{}
	for {
		currentTime := time.Now()

		tmp := []notifData{}
		mutex.Lock()
		if len(listNotifData) > 0 {
			for _, data := range listNotifData {
				dataTime := time.Unix(data.timeToExecute, 0)
				if dataTime.Minute() == currentTime.Minute() {
					SendNotif(data.msg)
				} else {
					tmp = append(tmp, data)
				}
			}
		}
		listNotifData = tmp
		mutex.Unlock()

		time.Sleep(time.Duration(60) * time.Second)
	}
}

func SendNotifWithDelay(msg string, delayInSecond int64) {
	data := notifData{
		msg:           msg,
		timeToExecute: time.Now().Unix() + delayInSecond,
	}

	mutex.Lock()
	listNotifData = append(listNotifData, data)
	mutex.Unlock()
}

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
