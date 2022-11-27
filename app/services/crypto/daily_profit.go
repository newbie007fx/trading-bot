package crypto

import (
	"fmt"
	"log"
	"strconv"
	"telebot-trading/app/helper"
	"time"
)

const profitKey string = "dailyProfit"
const timeKey string = "currentDay"
const profitThreshold float32 = 6

func SetProfit(profit float32) {
	currentTimeKey := getKeyTime()
	storedTimeKey := populateKeyTimeFromStorage()

	storedProfit := populateCurrentProfitFromStorage()
	totalProfit := storedProfit + profit
	storeProfitToStorage(totalProfit)

	if currentTimeKey != storedTimeKey {
		storeKeyTimeToStorage(currentTimeKey)
	}

	log.Printf("accumulated profit on time %d is = %f ", currentTimeKey, totalProfit)
}

func IsProfitMoreThanThreshold() bool {
	currentTimeKey := getKeyTime()
	storedTimeKey := populateKeyTimeFromStorage()
	if currentTimeKey == storedTimeKey {
		storedProfit := populateCurrentProfitFromStorage()
		return storedProfit > profitThreshold
	}

	storeProfitToStorage(0)
	storeKeyTimeToStorage(currentTimeKey)

	return false
}

func populateCurrentProfitFromStorage() float32 {
	profitString := helper.GetSimpleStore().Get(profitKey)
	if profitString != nil {
		if s, err := strconv.ParseFloat(*profitString, 32); err == nil {
			return float32(s)
		}
	}

	return float32(0)
}

func populateKeyTimeFromStorage() int64 {
	timeString := helper.GetSimpleStore().Get(timeKey)
	if timeString != nil {
		if s, err := strconv.ParseInt(*timeString, 10, 64); err == nil {
			return s
		}
	}

	return int64(0)
}

func storeProfitToStorage(profit float32) {
	s := fmt.Sprintf("%f", profit)

	helper.GetSimpleStore().Set(profitKey, s)
}

func storeKeyTimeToStorage(timeVal int64) {
	s := fmt.Sprintf("%d", timeVal)

	helper.GetSimpleStore().Set(timeKey, s)
}

func getKeyTime() int64 {
	now := time.Now()
	keyDate := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return keyDate.Unix()
}
