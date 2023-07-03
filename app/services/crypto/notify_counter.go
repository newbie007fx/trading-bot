package crypto

import (
	"fmt"
	"log"
	"strconv"
	"telebot-trading/app/repositories"
)

const notifyCounterKey = "notify_counter"

func SetNotifyCounter(counter uint) {
	notifyCounter := PopulateCurrentNotifyCounter()
	if counter > uint(notifyCounter) {
		storeNotifyCounterToStorage(counter)
	}
}

func IsShouldNotify() bool {
	notifyCounter := PopulateCurrentNotifyCounter()
	if notifyCounter == 0 {
		return false
	}

	return true
}

func NotifyCounterDecrement() {
	notifyCounter := PopulateCurrentNotifyCounter()
	if notifyCounter == 0 {
		return
	}

	notifyCounter--
	storeNotifyCounterToStorage(notifyCounter)
}

func PopulateCurrentNotifyCounter() uint {
	notifyCounterString := repositories.GetConfigValueByName(notifyCounterKey)
	if notifyCounterString != nil {
		if notifyCounter, err := strconv.ParseInt(*notifyCounterString, 10, 64); err == nil {
			return uint(notifyCounter)

		}
	}

	return uint(0)
}

func storeNotifyCounterToStorage(counter uint) {
	s := fmt.Sprintf("%d", counter)

	err := repositories.SetConfigByName(notifyCounterKey, s)
	if err != nil {
		log.Println(err.Error())
	}
}
