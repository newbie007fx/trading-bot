package analysis

import (
	"fmt"
	"telebot-trading/app/models"
	"time"
)

var reason string = ""

func GetHigestHightPrice(bands []models.Band) float32 {
	var highest float32 = 0
	for _, band := range bands {
		if highest < band.Candle.Hight {
			highest = band.Candle.Hight
		}
	}

	return highest
}

func GetLowestLowPrice(bands []models.Band) float32 {
	var lowest float32 = bands[0].Candle.Low
	for _, band := range bands {
		if lowest > band.Candle.Low {
			lowest = band.Candle.Low
		}
	}

	return lowest
}

func GetSellReason() string {
	return reason
}

func CheckIsNeedSellOnTrendUp(currencyConfig *models.CurrencyNotifConfig, shortInterval models.BandResult, currentTime time.Time) bool {
	holdTime := time.Unix(currencyConfig.HoldedAt, 0)
	holdedHour := calculateHoldTimeInHour(holdTime, currentTime)
	var threshold float32 = 3
	if holdedHour > 5 {
		threshold = 1
	} else if holdedHour > 2 {
		threshold = 2
	}

	if currencyConfig.HoldPrice > shortInterval.CurrentPrice {
		changes := currencyConfig.HoldPrice - shortInterval.CurrentPrice
		changesInPercent := changes / currencyConfig.HoldPrice * 100
		if (shortInterval.Direction == BAND_DOWN && changesInPercent > 3) || holdedHour > 11 {
			reason = fmt.Sprintf("sell on defisit after holded %d hours", holdedHour)
			return true
		}
	} else {
		changes := shortInterval.CurrentPrice - currencyConfig.HoldPrice
		changesInPercent := changes / currencyConfig.HoldPrice * 100
		if changesInPercent > threshold || holdedHour > 11 {
			reason = fmt.Sprintf("sell on profit after holded %d hours", holdedHour)
			return true
		}
	}

	return false
}

func calculateHoldTimeInHour(holdTime, currentTime time.Time) int {
	var utcZone, _ = time.LoadLocation("UTC")
	holdTime = holdTime.In(utcZone)
	currentTime = currentTime.In(utcZone)

	result := currentTime.Sub(holdTime)

	return int(result.Hours())
}
