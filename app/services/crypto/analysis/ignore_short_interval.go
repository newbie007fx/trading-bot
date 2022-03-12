package analysis

import (
	"telebot-trading/app/models"
	"time"
)

func IsIgnored(result *models.BandResult, requestTime time.Time) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	percentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100
	percentFromSMA := (lastBand.SMA - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100

	if lastBandHeadDoubleBody(result) {
		if lastBand.Candle.Close > float32(lastBand.SMA) || (result.AllTrend.ShortTrend == models.TREND_DOWN && lastBand.Candle.Close < float32(lastBand.SMA)) {
			ignoredReason = "lastBandHeadDoubleBody"
			return true
		}
	}

	if isContaineBearishEngulfing(result) && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if whenHeightTripleAverage(result) {
		ignoredReason = "whenHeightTripleAverage"
		return true
	}

	secondLastBand := result.Bands[len(result.Bands)-2]
	if IsHammer(lastBand) && result.AllTrend.SecondTrend == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "hammer pattern"
		return true
	}

	if (IsDoji(lastBand, true) || secondAlgDoji(lastBand)) && result.AllTrend.SecondTrend == models.TREND_UP && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "doji pattern"
		return true
	}

	if ThreeBlackCrowds((result.Bands[len(result.Bands)-4:])) {
		ignoredReason = "three black crowds pattern"
		return true
	}

	if isUpSignificanAndNotUp(result) {
		ignoredReason = "after up significan and trend not up"
		return true
	}

	if isLastBandChangeMoreThan5AndHeadMoreThan3(lastBand) {
		ignoredReason = "last band change more than 5 and head more than 3"
		return true
	}

	if secondLastBand.Candle.Open > float32(secondLastBand.Upper) && secondLastBand.Candle.Close < float32(secondLastBand.Upper) {
		if lastBand.Candle.Close < float32(lastBand.Upper) {
			ignoredReason = "below sma and previous band down from upper"
			return true
		}
	}

	if isHasCrossUpper(result.Bands[:len(result.Bands)/2], true) && result.AllTrend.Trend == models.TREND_UP && result.Position == models.ABOVE_SMA {
		if !isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) && countBelowSMA(result.Bands[len(result.Bands)-5:], true) >= 1 {
			if result.PriceChanges > 3 && percentFromUpper < 1.5 {
				ignoredReason = "up down, above sma and margin < 1.5"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA {
		if isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) && result.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "above sma and short trend not up"
			return true
		}

		if !isHasCrossUpper(result.Bands[len(result.Bands)/2:], true) && !isHasCrossSMA(result.Bands[len(result.Bands)/2:], false) {
			if percentFromUpper < 3.2 && result.AllTrend.SecondTrend < 35 {
				ignoredReason = "above sma and 10 band not cross upper or sma"
				return true
			}
		}

		if !isHasCrossUpper(result.Bands[len(result.Bands)-10:], true) && isHasCrossUpper(result.Bands[len(result.Bands)-15:], true) {
			if percentFromUpper < 3.2 {
				ignoredReason = "above sma and percent below 3 2nd logic"
				return true
			}
		}
	}

	if isHasCrossLower(result.Bands[len(result.Bands)-5:], false) && isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) {
		lowPrice := getLowestPrice(result.Bands[len(result.Bands)-5:])
		hightPrice := getHigestPrice(result.Bands[len(result.Bands)-5:])
		percent := (hightPrice - lowPrice) / lowPrice * 100
		if percent > 3.5 {
			ignoredReason = "percen changes more than 3.5"
			return true
		}
	}

	if upperLowerMarginBelow3(*result) {
		ignoredReason = "upper lower margin below 3"
		return true
	}

	if lastBand.Candle.Open > float32(lastBand.Upper) {
		ignoredReason = "open close above uper"
		return true
	}

	bandPercent := (lastBand.Candle.Close - lastBand.Candle.Open) / lastBand.Candle.Open * 100
	if countDownBand(result.Bands[len(result.Bands)-4:]) > 2 && bandPercent < 0.2 {
		ignoredReason = "significant down and just minor up"
		return true
	}

	if result.AllTrend.SecondTrend == BAND_DOWN || result.AllTrend.Trend == models.TREND_DOWN {
		if countBelowLower(result.Bands[len(result.Bands)-2:], false) > 0 || isHasCrossLower(result.Bands[len(result.Bands)-2:], true) {
			if percentFromSMA < 3.1 {
				ignoredReason = "down trend minor up, and percent below 3"
				return true
			}
		}
	}

	return ignored(result)
}
