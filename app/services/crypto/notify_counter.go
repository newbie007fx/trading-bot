package crypto

import (
	"fmt"
	"strconv"
	"telebot-trading/app/helper"
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

	notifyCounter--
	storeNotifyCounterToStorage(notifyCounter)

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
	notifyCounterString := helper.GetSimpleStore().Get(notifyCounterKey)
	if notifyCounterString != nil {
		if notifyCounter, err := strconv.ParseInt(*notifyCounterString, 10, 64); err == nil {
			return uint(notifyCounter)

		}
	}

	return uint(0)
}

func storeNotifyCounterToStorage(counter uint) {
	s := fmt.Sprintf("%d", counter)

	helper.GetSimpleStore().Set(notifyCounterKey, s)
}
