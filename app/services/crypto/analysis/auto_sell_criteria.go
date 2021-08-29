package analysis

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
)

var reason string = ""

func IsNeedToSell(result models.BandResult) bool {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(result.Symbol)
	if err != nil {
		log.Panicln(err.Error())
	}

	lastBand := result.Bands[len(result.Bands)-1]
	if currencyConfig.HoldPrice > result.CurrentPrice {
		changes := currencyConfig.HoldPrice - result.CurrentPrice
		changesInPercent := changes / currencyConfig.HoldPrice * 100
		if changesInPercent >= 3 && result.Direction == BAND_DOWN {
			if result.Position == models.BELOW_LOWER {
				reason = "sell on down with criteria 1"
				return true
			} else if result.Position == models.BELOW_SMA {
				changesFromLower := currencyConfig.HoldPrice - float32(lastBand.Lower)
				changesFromLowerPercent := changesFromLower / currencyConfig.HoldPrice * 100
				if changesFromLowerPercent > 4 {
					reason = "sell on down with criteria 2"
					return true
				}
			} else if result.Position == models.ABOVE_SMA {
				changesFromSMA := currencyConfig.HoldPrice - float32(lastBand.SMA)
				changesFromSMAPercent := changesFromSMA / currencyConfig.HoldPrice * 100
				if changesFromSMAPercent > 4 {
					reason = "sell on down with criteria 3"
					return true
				}
			} else if result.Position == models.ABOVE_UPPER {
				changesFromUpper := currencyConfig.HoldPrice - float32(lastBand.Upper)
				changesFromUpperPercent := changesFromUpper / currencyConfig.HoldPrice * 100
				if changesFromUpperPercent > 4 {
					reason = "sell on down with criteria 4"
					return true
				}
			}
		} else if changesInPercent > 4 {
			reason = "sell on down with criteria 5"
			return true
		}
	} else {
		changes := result.CurrentPrice - currencyConfig.HoldPrice
		changesInPercent := changes / currencyConfig.HoldPrice * 100
		highest := getHigestPrice(result.Bands)
		highestChangePercent := changes / (highest - currencyConfig.HoldPrice) * 100
		lastFourData := result.Bands[len(result.Bands)-4 : len(result.Bands)]

		if highestChangePercent >= 34 && changesInPercent >= 3 && CalculateTrends(lastFourData) == models.TREND_DOWN && result.Direction == BAND_DOWN {

			secondLastBand := result.Bands[len(result.Bands)-2]
			if result.Position == models.BELOW_LOWER {
				if lastBand.Candle.Open > float32(lastBand.Lower) && float32(lastBand.Lower) > result.CurrentPrice {
					changesOnLower := result.CurrentPrice - float32(lastBand.Lower)
					changesOnLowerPercent := changesOnLower / float32(lastBand.Lower) * 100
					if changesOnLowerPercent >= 3 {
						reason = "sell on up with criteria 1"
						return true
					}
				}

				if secondLastBand.Candle.Close > float32(secondLastBand.Lower) || secondLastBand.Candle.Open > float32(secondLastBand.Lower) {
					if lastBand.Candle.Close < float32(lastBand.Lower) && lastBand.Candle.Open < float32(lastBand.Lower) && lastBand.Candle.Open > lastBand.Candle.Close {
						reason = "sell on up with criteria 2"
						return true
					}
				}
			} else if result.Position == models.BELOW_SMA {
				changesFromLower := result.CurrentPrice - float32(lastBand.Lower)
				changesFromLowerPercent := changesFromLower / result.CurrentPrice * 100
				if changesFromLowerPercent <= 1 {
					return false
				}

				if lastBand.Candle.Open > float32(lastBand.SMA) && float32(lastBand.SMA) > result.CurrentPrice {
					changesOnSMA := result.CurrentPrice - float32(lastBand.SMA)
					changesOnSMAPercent := changesOnSMA / float32(lastBand.SMA) * 100
					if changesOnSMAPercent >= 3 {
						reason = "sell on up with criteria 3"
						return true
					}
				}

				if secondLastBand.Candle.Close > float32(secondLastBand.SMA) || secondLastBand.Candle.Open > float32(secondLastBand.SMA) {
					if lastBand.Candle.Close < float32(lastBand.SMA) && lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Open > lastBand.Candle.Close {
						reason = "sell on up with criteria 4"
						return true
					}
				}
			} else if result.Position == models.ABOVE_SMA {
				changesFromSMA := result.CurrentPrice - float32(lastBand.SMA)
				changesFromSMAPercent := changesFromSMA / result.CurrentPrice * 100
				if changesFromSMAPercent <= 1 {
					return false
				}

				if lastBand.Candle.Open > float32(lastBand.Upper) && float32(lastBand.Upper) > result.CurrentPrice {
					changesOnUpper := result.CurrentPrice - float32(lastBand.Upper)
					changesOnUpperPercent := changesOnUpper / float32(lastBand.Upper) * 100
					if changesOnUpperPercent >= 3 {
						reason = "sell on up with criteria 5"
						return true
					}
				}

				if secondLastBand.Candle.Close > float32(secondLastBand.Upper) || secondLastBand.Candle.Open > float32(secondLastBand.Upper) {
					if lastBand.Candle.Close < float32(lastBand.Upper) && lastBand.Candle.Open < float32(lastBand.Upper) && lastBand.Candle.Open > lastBand.Candle.Close {
						reason = "sell on up with criteria 6"
						return true
					}
				}
			} else if result.Position == models.ABOVE_UPPER {
				changesFromUpper := result.CurrentPrice - float32(lastBand.Upper)
				changesFromUpperPercent := changesFromUpper / result.CurrentPrice * 100
				if changesFromUpperPercent <= 1 {
					return false
				}
			}

			reason = "sell on up with criteria 7"
			return true

		} else if highestChangePercent < 34 && changesInPercent >= 3 {
			reason = "sell on up with criteria 8"
			return true
		}
	}

	return false
}

func getHigestPrice(bands []models.Band) float32 {
	var highest float32 = 0
	for _, band := range bands {
		if highest < band.Candle.Close {
			highest = band.Candle.Close
		}
	}

	return highest
}

func GetSellReason() string {
	return reason
}
