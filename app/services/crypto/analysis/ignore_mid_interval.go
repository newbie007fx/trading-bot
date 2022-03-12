package analysis

import (
	"fmt"
	"telebot-trading/app/models"
)

func IsIgnoredMidInterval(result *models.BandResult, shortInterval *models.BandResult) bool {
	lastBand := result.Bands[len(result.Bands)-1]
	secondLastBand := result.Bands[len(result.Bands)-2]
	hMidFirstBand := result.HeuristicBand.FirstBand

	percentFromUpper := (hMidFirstBand.Upper - float64(hMidFirstBand.Candle.Close)) / float64(hMidFirstBand.Candle.Close) * 100
	percentFromSMA := (hMidFirstBand.SMA - float64(hMidFirstBand.Candle.Close)) / float64(hMidFirstBand.Candle.Close) * 100

	shortLastBand := shortInterval.Bands[len(shortInterval.Bands)-1]
	hShortSecondBand := shortInterval.HeuristicBand.SecondBand
	hShortFourthBand := shortInterval.HeuristicBand.FourthBand
	shortPercentFromUpper := (float32(hShortSecondBand.Upper) - hShortSecondBand.Candle.Close) / hShortSecondBand.Candle.Close * 100
	shortHFourthPercentFromUpper := (hShortFourthBand.Upper - float64(hShortFourthBand.Candle.Close)) / float64(hShortFourthBand.Candle.Close) * 100

	if isInAboveUpperBandAndDownTrend(result) && !isLastBandOrPreviousBandCrossSMA(result.Bands) && !reversalFromLower(*shortInterval) {
		ignoredReason = "isInAboveUpperBandAndDownTrend"
		return true
	}

	if lastBandHeadDoubleBody(result) && lastBand.Candle.Close > float32(lastBand.SMA) && !(shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)/2:], false)) {
		ignoredReason = "lastBandHeadDoubleBody"
		return true
	}

	lastThreeExceptLastBand := result.Bands[len(result.Bands)-4 : len(result.Bands)-1]
	shortUp := CalculateTrendShort(result.Bands[len(result.Bands)-5:]) == models.TREND_UP || (CalculateTrendShort(result.Bands[len(result.Bands)-4:]) == models.TREND_UP && (isHasCrossSMA(lastThreeExceptLastBand, true) || isHasCrossLower(lastThreeExceptLastBand, false)))
	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && !shortUp {
		ignoredReason = "first trend up and second sideway"
		return true
	}

	if result.Position == models.ABOVE_UPPER && shortInterval.AllTrend.Trend != models.TREND_UP {
		if result.AllTrend.SecondTrend != models.TREND_UP {
			ignoredReason = "position above upper trend not up"
			return true
		}

		if (CountUpBand(result.Bands[len(result.Bands)-3:]) < 2 || !isHasCrossUpper(result.Bands[len(result.Bands)-4:len(result.Bands)-1], false)) && (result.AllTrend.SecondTrend != models.TREND_UP) {
			ignoredReason = "position above upper but previous band not upper or count up bellow 3"
			return true
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER {
		if CalculateTrendShort(result.Bands[len(result.Bands)-4:]) != models.TREND_UP || shortInterval.AllTrend.Trend != models.TREND_UP {
			ignoredReason = "short interval above upper and mid trend down or trend not up"
			return true
		}

		checkTrend := (result.AllTrend.FirstTrend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN) || (result.AllTrend.FirstTrend == models.TREND_SIDEWAY && result.AllTrend.SecondTrend == models.TREND_SIDEWAY)
		numberBelowSMA := countBelowSMA(result.Bands[len(result.Bands)/2:], true)
		if numberBelowSMA < 2 && checkTrend && !(isHasCrossLower(result.Bands[len(result.Bands)/2:], false) && lastBand.Candle.Close > float32(lastBand.SMA) && result.AllTrend.Trend != models.TREND_DOWN) {
			ignoredReason = "short interval above upper and mid inteval trend not up"
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	if isContaineBearishEngulfing(result) && !isLastBandOrPreviousBandCrossSMA(result.Bands) && shortPercentFromUpper < 3.2 {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if isLastBandOrPreviousBandCrossSMA(result.Bands) && result.AllTrend.SecondTrend == models.TREND_UP {
		if shortInterval.Position == models.ABOVE_SMA {
			if shortPercentFromUpper < 1.3 {
				ignoredReason = "mid cross band and short candle near upper band"
				return true
			}
		}

		if isLastBandOrPreviousBandCrossSMA(shortInterval.Bands) {
			ignoredReason = "mid cross band and short candle cross band"
			return true
		}
	}

	if (lastBand.Candle.Open >= float32(lastBand.Upper) && lastBand.Candle.Close > float32(lastBand.Upper)) || (secondLastBand.Candle.Close >= float32(secondLastBand.Upper) && secondLastBand.Candle.Open >= float32(secondLastBand.Upper)) {
		ignoredReason = fmt.Sprintf("open close above upper, %.4f, %.4f", lastBand.Candle.Open, lastBand.Upper)
		return true
	}

	if shortInterval.Position == models.ABOVE_SMA && result.AllTrend.Trend != models.TREND_UP && !isReversal(result.Bands) && !isReversal(shortInterval.Bands) {
		if !isHasCrossLower(result.Bands[len(result.Bands)-7:], false) {
			ignoredReason = "above sma and trend not up"
			return true
		}
	}

	if shortInterval.AllTrend.FirstTrend != models.TREND_UP && shortInterval.AllTrend.SecondTrend == models.TREND_UP && (result.AllTrend.Trend != models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP) {
		shortIntervalHalfBands := shortInterval.Bands[:len(shortInterval.Bands)/2]
		if (shortInterval.AllTrend.FirstTrendPercent <= 50 || isHasCrossLower(shortIntervalHalfBands, false)) && shortInterval.AllTrend.SecondTrendPercent <= 50 {
			if shortInterval.Position == models.ABOVE_SMA {
				if result.Position == models.BELOW_SMA {
					if shortPercentFromUpper < 3.2 && percentFromSMA < 3.2 {
						ignoredReason = "up down, and margin form up bellow threshold"
						return true
					}
				}

				if result.Position == models.ABOVE_SMA {
					percentFromUpper := (float32(lastBand.Upper) - lastBand.Candle.Close) / lastBand.Candle.Close * 100
					if shortPercentFromUpper < 3.2 && percentFromUpper < 3.4 {
						ignoredReason = "up down, and margin form up bellow threshold"
						return true
					}
				}
			}
		}
	}

	highest := GetHigestHightPrice(result.Bands)
	lowest := GetLowestLowPrice(result.Bands)
	difference := highest - lowest
	percent := difference / lowest * 100
	if percent <= 2 {
		ignoredReason = "hight and low bellow 3"
		return true
	}

	if result.AllTrend.Trend == models.TREND_UP && result.AllTrend.SecondTrendPercent < 5 {
		longestIndex := getLongestCandleIndex(result.Bands[len(result.Bands)/2:])
		secondWaveTrendDetail := CalculateTrendsDetail(result.Bands[longestIndex+len(result.Bands)/2:])
		if secondWaveTrendDetail.FirstTrend != models.TREND_UP || secondWaveTrendDetail.SecondTrend != models.TREND_UP {
			ignoredReason = "after significan up and not up up"
			return true
		}
	}

	if isDoubleUp(result.Bands) {
		ignoredReason = "has double up"
		return true
	}

	if result.AllTrend.SecondTrendPercent < 7 && lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close < float32(lastBand.SMA) {
		if CalculateTrendShort(result.Bands[len(result.Bands)-5:]) != models.TREND_UP && !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) {
			if !isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-4:], false) {
				ignoredReason = "significan down and below sma, last 5 not up trend, not cross lower"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_UPPER && !isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-5:], true) {
		if shortPercentFromUpper < 3.2 {
			ignoredReason = "above upper and short interval 5 band not cross upper and margin below 3"
			return true
		}
	}

	significanDown := result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 10 && result.AllTrend.SecondTrend == models.TREND_DOWN
	secondWaweBelowSMA := result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 15 && countBelowSMA(result.Bands[len(result.Bands)/2:], false) >= len(result.Bands)/2
	if significanDown || secondWaweBelowSMA {
		if !isHasCrossSMA(result.Bands[len(result.Bands)-5:], true) && !isHasCrossLower(result.Bands[len(result.Bands)-5:], false) && result.AllTrend.ShortTrend != models.TREND_UP {
			ignoredReason = "first wave significan down, second still donw, below sma, not cross sma or lower"
			return true
		}
	}

	if shortLastBand.Candle.Hight > float32(shortLastBand.Upper) && result.Position == models.ABOVE_SMA {
		if !isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) && isHasCrossLower(result.Bands[len(result.Bands)-5:], false) {
			if percentFromUpper < 3.4 {
				ignoredReason = "short interval cross upper, mid not cross upper but lower and margin below 3"
				return true
			}
		}
	}

	if crossSMAAndPreviousBandNotHaveAboveSMA(result.Bands) {
		ignoredReason = "cross sma and all previous band do not above sma"
		return true
	}

	if isHasCrossUpper(result.Bands[:len(result.Bands)/2], true) && !isHasCrossUpper(result.Bands[len(result.Bands)/2:], true) {
		if result.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_UPPER && percentFromUpper < 3.4 {
			ignoredReason = "short interval above upper, mid interval above sma and margin below 3"
			return true
		}
	}

	if getHighestIndex(result.Bands) == len(result.Bands)-1 && percentFromUpper < 3.4 && result.Position == models.ABOVE_SMA {
		if isHasCrossLower(shortInterval.Bands[:len(shortInterval.Bands)/2], false) && isHasCrossUpper(shortInterval.Bands[:len(shortInterval.Bands)/2], true) {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)/2:], true) {
				ignoredReason = "short interval above upper, mid interval above sma and margin below 3 - second"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && countBelowSMA(result.Bands[len(result.Bands)/2:], false) >= len(result.Bands)/2 {
		if lastBand.Candle.Open < float32(lastBand.Lower) && lastBand.Candle.Close < float32(lastBand.Lower) && percentFromSMA < 3.4 {
			ignoredReason = "open close below lower and margin below 3.4"
			return true
		}
	}

	margin := lastBand.Candle.Close - lastBand.Candle.Open
	if lastBand.Candle.Open+margin < float32(lastBand.Lower) {
		ignoredReason = "open close below avg lower"
		return true
	}

	shortPercentFromSMA := (shortLastBand.SMA - float64(shortLastBand.Candle.Close)) / float64(shortLastBand.Candle.Close) * 100
	if countAboveSMA(result.Bands) == 0 && result.AllTrend.ShortTrend != models.TREND_UP {
		if countAboveSMA(shortInterval.Bands) == 0 && shortPercentFromSMA < 3.1 && !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) {
			ignoredReason = "trend down, short margin from sma < 3"
			return true
		}
	}

	if secondLastBand.Candle.Close < secondLastBand.Candle.Open {
		if shortInterval.PriceChanges > 3 && countCrossUpper(shortInterval.Bands) == 0 && IsHammer(secondLastBand) {
			ignoredReason = "previous band down and hammer pattern"
			return true
		}
	}

	if result.AllTrend.ShortTrend != models.TREND_UP {
		shortLastBandChanges := (shortLastBand.Candle.Close - shortLastBand.Candle.Open) / shortLastBand.Candle.Open * 100
		if shortLastBandChanges > 3 {
			if shortInterval.Position == models.BELOW_SMA && shortPercentFromSMA < 3.1 {
				ignoredReason = "mid short interval down. short below sma high changes, percen below threshold"
				return true
			}

			if shortInterval.Position == models.ABOVE_SMA && shortPercentFromUpper < 3.2 {
				ignoredReason = "mid short interval down. short above sma high changes, percen below threshold"
				return true
			}
		}
	}

	if shortInterval.PriceChanges > 3 && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
		if isHasCrossSMA(result.Bands[len(result.Bands)-2:], true) && result.Position == models.ABOVE_SMA {
			if lastBand.Candle.Open > float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA) && countAboveSMA(result.Bands) == 1 {
				ignoredReason = "above sma but percent from upper below 3"
				return true
			}
		}
	}

	if lastBand.Candle.Close > float32(lastBand.SMA) {
		bodyShort := shortLastBand.Candle.Close - shortLastBand.Candle.Open
		if shortInterval.Position == models.ABOVE_UPPER && shortInterval.PriceChanges > 3 && shortInterval.PriceChanges < 5 {
			if bodyShort/2+shortLastBand.Candle.Open > float32(shortLastBand.Upper) {
				ignoredReason = "short Above upper"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
		if isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)/2:], false) && !isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)/2:], true) {
			if shortPercentFromUpper < 3.2 {
				ignoredReason = "short above sma and percent below 3"
				return true
			}
		}
	}

	if result.Position == models.BELOW_SMA && countCrossLower(result.Bands[len(result.Bands)-5:]) > 3 && result.AllTrend.ShortTrend != models.TREND_UP {
		if shortInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(shortInterval.Bands, true) {
			if shortPercentFromUpper < 3.2 && percentFromHigest(result.Bands) < 3 {
				ignoredReason = "short above sma and percent below 3 2nd"
				return true
			}
		}
	}

	if (lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA)) || result.Position == models.BELOW_SMA {
		if result.AllTrend.SecondTrend == models.TREND_DOWN && !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
			if shortLastBand.Candle.Open < float32(shortLastBand.SMA) && shortLastBand.Candle.Close > float32(shortLastBand.SMA) && shortPercentFromUpper < 3.2 {
				ignoredReason = "down, not cross lowe. short cross sma but percent from upper below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent < 5 {
		if secondLastBand.Candle.Open < float32(secondLastBand.Lower) || secondLastBand.Candle.Close < float32(secondLastBand.Lower) {
			if shortPercentFromUpper < 3.2 {
				ignoredReason = "significan down. short percent from sma below 3"
				return true
			}
		}
	}

	shortSecondLastBand := shortInterval.Bands[len(shortInterval.Bands)-2]
	if result.Position == models.ABOVE_UPPER || percentFromUpper < 3.4 {
		if shortLastBand.Candle.Open >= float32(shortLastBand.Upper) || (shortSecondLastBand.Candle.Open >= float32(shortSecondLastBand.Upper) && shortSecondLastBand.Candle.Close >= float32(shortSecondLastBand.Upper)) {
			ignoredReason = "short lasband or previous band open close above upper"
			return true
		}
	}

	if lastBand.Candle.Close < float32(lastBand.Upper) && percentFromUpper < 3.4 {
		if shortInterval.Position == models.ABOVE_UPPER && upperSmaMarginBelow3(*result) {
			ignoredReason = "short cross upper, mid after cross upper hit sma but percent below 3"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.FirstTrendPercent < 1.5 && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if !isHasCrossUpper(shortInterval.Bands, true) && isHasCrossLower(shortInterval.Bands, false) {
			if shortInterval.Position == models.ABOVE_SMA && shortPercentFromUpper < 3.2 {
				ignoredReason = "after significan up, short cross lower, above sma but percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_SIDEWAY && percentFromUpper < 3.4 && result.AllTrend.ShortTrend == models.TREND_DOWN {
		ignoredReason = "trend sideway, but percent from upper below 3"
		return true
	}

	if countBelowLower(result.Bands[len(result.Bands)-3:], false) > 0 && percentFromSMA < 3.1 {
		ignoredReason = "contain open close below lower"
		return true
	}

	if shortInterval.Position == models.ABOVE_UPPER || isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
		if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 10 && result.AllTrend.SecondTrend != models.TREND_DOWN {
			if countCrossUpper(result.Bands[len(result.Bands)/2:]) == 1 || countCrossSMA(result.Bands[len(result.Bands)/2:]) == 1 {
				ignoredReason = "after down tren up cross sma"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if result.Position == models.ABOVE_SMA && !isHasCrossSMA(result.Bands[len(result.Bands)-2:], false) {
			ignoredReason = "second trend up, short down. above sma but not cross sma"
			return true
		}
	}

	if countBelowLower(result.Bands[len(result.Bands)-3:], false) > 0 {
		if shortPercentFromSMA < 3.1 {
			ignoredReason = "have below lower and short interval percent below 3"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend != models.TREND_UP {
		if (result.Position == models.BELOW_SMA && percentFromSMA < 3) || (result.Position == models.ABOVE_SMA && percentFromUpper < 3.4) {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-2:], true) {
				ignoredReason = "percent below 3 and short interval cross upper"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP {
		if shortInterval.AllTrend.SecondTrend != models.TREND_UP && shortInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.Position == models.BELOW_SMA && !isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-2:], false) {
				ignoredReason = "possibility down, skip"
				return true
			}
		}
	}

	if isHasCrossLower(result.Bands[len(result.Bands)-2:], true) {
		if shortHFourthPercentFromUpper < 3.2 {
			ignoredReason = "down and just minor up"
			return true
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossLower(result.Bands[len(result.Bands)-1:], true) {
		if isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-4:], true) {
			ignoredReason = "trend down and cros lower on body"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands, false) {
		if shortHFourthPercentFromUpper < 3.2 && shortPercentFromUpper < 3.2 {
			ignoredReason = "trend down and not cross lower and percent from upper below 3"
			return true
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossLower(result.Bands[len(result.Bands)-2:], true) && countDownBand(result.Bands[len(result.Bands)-2:]) > 0 {
		ignoredReason = "trend down and contain band cross lower on body"
		return true
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossSMA(result.Bands[len(result.Bands)-2:], false) {
		if isHasOpenCloseAboveUpper(result.Bands[len(result.Bands)-4:]) && !isHasCrossLower(result.Bands[len(result.Bands)-2:], false) {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.Position == models.BELOW_SMA {
				ignoredReason = "trend down and contain band cross lower on body 2nd"
				return true
			}
		}
	}

	if !isHasCrossLower(result.Bands[len(result.Bands)-2:], false) {
		if shortInterval.AllTrend.SecondTrend != models.TREND_UP && shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortPercentFromUpper < 3 {
			lowPrice := getLowestPrice(shortInterval.Bands)
			hightPrice := getHigestPrice(shortInterval.Bands)
			percent := (hightPrice - lowPrice) / lowPrice * 100
			if percent < 3.2 && shortInterval.Position != models.ABOVE_UPPER && !(shortInterval.Position == models.BELOW_SMA && isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-2:], false)) {
				ignoredReason = "sideway"
				return true
			}
		}

		if result.AllTrend.SecondTrend != models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && percentFromUpper < 3.4 {
			lowPrice := getLowestPrice(result.Bands[len(result.Bands)/2:])
			hightPrice := getHigestPrice(result.Bands[len(result.Bands)/2:])
			percent := (hightPrice - lowPrice) / lowPrice * 100
			if percent < 3 && result.Position != models.ABOVE_UPPER {
				ignoredReason = "sideway 2nd"
				return true
			}
		}
	}

	return false
}
