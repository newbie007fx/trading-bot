package analysis

import (
	"telebot-trading/app/models"
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

func CheckIsNeedSellOnTrendUp(currencyConfig *models.CurrencyNotifConfig, shortInterval models.BandResult) bool {
	if currencyConfig.HoldPrice > shortInterval.CurrentPrice {
		changes := currencyConfig.HoldPrice - shortInterval.CurrentPrice
		changesInPercent := changes / currencyConfig.HoldPrice * 100
		if shortInterval.Direction == BAND_DOWN && changesInPercent > 3 {
			reason = "sell on profit"
			return true
		}
	} else {
		changes := shortInterval.CurrentPrice - currencyConfig.HoldPrice
		changesInPercent := changes / currencyConfig.HoldPrice * 100
		if changesInPercent > 3 {
			reason = "sell on profit"
			return true
		}
	}

	return false
}
