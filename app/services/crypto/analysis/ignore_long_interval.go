package analysis

import (
	"telebot-trading/app/models"
)

func IsIgnoredLongInterval(result *models.BandResult, shortInterval *models.BandResult, midInterval *models.BandResult) bool {
	bandLen := len(result.Bands)

	lastBand := result.Bands[len(result.Bands)-1]
	hFirstBand := result.HeuristicBand.FirstBand
	secondLastBand := result.Bands[len(result.Bands)-2]

	midLastBand := midInterval.Bands[len(result.Bands)-1]
	hMidFirstBand := midInterval.HeuristicBand.FirstBand
	midSecondLastBand := midInterval.Bands[len(result.Bands)-2]

	shortLastBand := shortInterval.Bands[len(result.Bands)-1]
	hShortSecondBand := shortInterval.HeuristicBand.SecondBand
	hShortFourthBand := shortInterval.HeuristicBand.FourthBand
	shortSecondLastBand := shortInterval.Bands[len(shortInterval.Bands)-2]

	percentFromHeight := (hFirstBand.Candle.Hight - hFirstBand.Candle.Close) / hFirstBand.Candle.Close * 100
	hFirstPercentFromUpper := (hFirstBand.Upper - float64(hFirstBand.Candle.Close)) / float64(hFirstBand.Candle.Close) * 100
	hFirstPercentFromSMA := (hFirstBand.SMA - float64(hFirstBand.Candle.Close)) / float64(hFirstBand.Candle.Close) * 100
	percentFromUpper := (lastBand.Upper - float64(lastBand.Candle.Close)) / float64(lastBand.Candle.Close) * 100

	midHFirstPercentFromUpper := (hMidFirstBand.Upper - float64(hMidFirstBand.Candle.Close)) / float64(hMidFirstBand.Candle.Close) * 100
	midHFirstPercentFromSMA := (hMidFirstBand.SMA - float64(hMidFirstBand.Candle.Close)) / float64(hMidFirstBand.Candle.Close) * 100
	midPercentFromUpper := (midLastBand.Upper - float64(midLastBand.Candle.Close)) / float64(midLastBand.Candle.Close) * 100

	shortHSecondPercentFromUpper := (hShortSecondBand.Upper - float64(hShortSecondBand.Candle.Close)) / float64(hShortSecondBand.Candle.Close) * 100
	shortHSecondPercentFromSMA := (hShortSecondBand.SMA - float64(hShortSecondBand.Candle.Close)) / float64(hShortSecondBand.Candle.Close) * 100
	shortHFourthPercentFromUpper := (hShortFourthBand.Upper - float64(hShortFourthBand.Candle.Close)) / float64(hShortFourthBand.Candle.Close) * 100

	if isInAboveUpperBandAndDownTrend(result) && lastBand.Candle.Hight > float32(lastBand.Upper) {
		if lastBand.Candle.Hight-lastBand.Candle.Close > lastBand.Candle.Close-lastBand.Candle.Open {
			ignoredReason = "isInAboveUpperBandAndDownTrend and down from upper"
			return true
		}
	}

	if isContaineBearishEngulfing(result) && lastBand.Candle.Close > float32(lastBand.SMA) {
		ignoredReason = "isContaineBearishEngulfing"
		return true
	}

	if result.Position == models.ABOVE_UPPER && result.AllTrend.ShortTrend != models.TREND_UP {
		ignoredReason = "position above upper and not trend up"
		return true
	}

	isMidIntervalTrendNotUp := CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-4:]) == models.TREND_DOWN
	isLongIntervalTrendNotUp := CalculateTrendShort(result.Bands[len(result.Bands)-3:]) == models.TREND_DOWN
	if shortInterval.Position == models.ABOVE_UPPER || midInterval.Position == models.ABOVE_UPPER || result.Position == models.ABOVE_UPPER {
		if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN || isMidIntervalTrendNotUp || isLongIntervalTrendNotUp {
			ignoredReason = "when above upper or master trend down and trend not up"
			return true
		}

		lenData := len(result.Bands)
		hight := GetHigestHightPrice(result.Bands[lenData-lenData/4:])
		low := getLowestPrice(result.Bands[lenData-lenData/4:])
		difference := hight - low
		percent := float32(difference) / float32(low) * 100
		if percent > 45 {
			ignoredReason = "up more than 45"
			return true
		}
	}

	if isTrendUpLastThreeBandHasDoji(result) {
		ignoredReason = "isTrendUpLastThreeBandHasDoji"
		return true
	}

	if midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
		if result.Position == models.ABOVE_SMA {
			if hFirstPercentFromUpper < 3.4 && midHFirstPercentFromUpper < 3.3 && shortHSecondPercentFromUpper < 3.2 {
				if shortInterval.AllTrend.SecondTrend != models.TREND_UP || midInterval.AllTrend.SecondTrend != models.TREND_UP || isLongIntervalTrendNotUp {
					ignoredReason = "all band bellow 3.1 from upper or not up trend"
					return true
				}
			}

			if isHasCrossSMA(result.Bands[len(result.Bands)-1:len(result.Bands)], true) {
				ignoredReason = "all interval above upper and long interval cross sma"
				return true
			}
		}

		if isLastBandOrPreviousBandCrossSMA(result.Bands) && midInterval.AllTrend.Trend != models.TREND_UP {
			ignoredReason = "short and mid above sma but long interval cross sma"
			return true
		}

		if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_UP {
			lowest := getLowestIndex(result.Bands)
			if lowest != len(result.Bands)-1 {
				ignoredReason = "short and mid above sma but long interval second wave down trend"
				return true
			}
		}
	}

	if isLastBandCrossUpperAndPreviousBandNot(shortInterval.Bands) {
		if isLastBandCrossUpperAndPreviousBandNot(midInterval.Bands) {
			if isLastBandCrossUpperAndPreviousBandNot(result.Bands) {
				ignoredReason = "band above upper and just one in all interval, skip"
				return true
			}
		}
	}

	if lastBand.Candle.Open > float32(lastBand.Upper) {
		ignoredReason = "open close above upper"
		return true
	}

	allTrendUp := midInterval.AllTrend.Trend == models.TREND_UP && result.AllTrend.Trend == models.TREND_UP
	if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-4:], false) || allTrendUp {
		if result.Position == models.ABOVE_UPPER {
			highestHightIndex := getHighestHightIndex(result.Bands)
			highestIndex := getHighestIndex(result.Bands)
			higestHightBand := result.Bands[highestHightIndex]
			percent := (higestHightBand.Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
			if highestHightIndex == len(result.Bands)-1 || (percent <= 3.2 && highestIndex == len(result.Bands)-1) {
				ignoredReason = "all interval above upper or all trend up and new hight created"
				return true
			}
		}

		if midInterval.Position == models.ABOVE_UPPER && result.AllTrend.Trend == models.TREND_DOWN {
			if lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Close > float32(lastBand.SMA) {
				ignoredReason = "short and mid above upper and long down cross sma"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && percentFromHeight < 3.4 {
		if !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
			ignoredReason = "sideway, above sma and percent from upper bellow 3"
			return true
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent < 10 {
		if result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
			ignoredReason = "significan down, below sma, not cross lower yet"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.FirstTrendPercent < 20 {
		if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.SecondTrendPercent < 20 {
			index := getHighestIndex(result.Bands[len(result.Bands)/2:])
			if index < len(result.Bands)/2 && CountSquentialUpBand(result.Bands[len(result.Bands)-3:]) < 2 && !isLastBandOrPreviousBandCrossSMA(result.Bands) {
				if index < (len(result.Bands)/2)-4 {
					index = (len(result.Bands) / 2) - 4
				}
				lastDataFromHight := result.Bands[(len(result.Bands)/2)+index:]
				if CalculateTrendShort(lastDataFromHight) != models.TREND_UP {
					ignoredReason = "down and up significan then down trend"
					return true
				}
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.SecondTrend < 10 {
		if lastBand.Candle.Hight < float32(lastBand.Upper) && secondLastBand.Candle.Hight < float32(secondLastBand.Upper) && hFirstPercentFromUpper < 3.4 {
			ignoredReason = "up significan but last two band not cross upper"
			return true
		}
	}

	if midInterval.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.Trend == models.TREND_DOWN {
		if midInterval.Position == models.BELOW_SMA && result.Position == models.BELOW_SMA && hFirstPercentFromUpper <= 3.4 {
			ignoredReason = "mid and long interval down and below sma and mergin form short below 3"
			return true
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_UP {
		if countBelowSMA(result.Bands[len(result.Bands)-6:len(result.Bands)-1], false) == 5 {
			if lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA) && shortHSecondPercentFromUpper <= 3.2 {
				ignoredReason = "long interval down and below sma (5) and cross sma && mergin form short below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrendPercent < 20 {
		secondWaveBands := result.Bands[len(result.Bands)/2:]
		if countCrossUpper(secondWaveBands) > 2 {
			secondLastBand := result.Bands[len(result.Bands)-2]
			if secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Open > float32(secondLastBand.Upper) {
				ignoredReason = "previous band down from upper"
				return true
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_UP && secondLastBand.Candle.Hight < float32(secondLastBand.Upper) {
		if result.AllTrend.FirstTrendPercent < 20 && result.AllTrend.FirstTrendPercent < result.AllTrend.SecondTrendPercent {
			if isHasCrossUpper(result.Bands[len(result.Bands)-1:], true) && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-1:], true) {
				if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-1:], true) {
					ignoredReason = "all interval cross upper"
					return true
				}
			}
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && midInterval.Position == models.ABOVE_UPPER && result.Position == models.ABOVE_SMA {
		secondWaveBands := result.Bands[len(result.Bands)/2:]
		highestIndex := getHighestIndex(secondWaveBands)
		if highestIndex == len(secondWaveBands)-1 {
			secondHightIndex := getHigestIndexSecond(secondWaveBands)
			if secondHightIndex < len(secondWaveBands)-5 {
				resultTrend := CalculateTrendsDetail(secondWaveBands[secondHightIndex:])
				if resultTrend.FirstTrend == models.TREND_DOWN && resultTrend.SecondTrend == models.TREND_UP {
					higestHightIndex := getHighestHightIndex(secondWaveBands)
					percentFromHight := (secondWaveBands[higestHightIndex].Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
					if percentFromHight < 3.4 && hFirstPercentFromUpper < 3.4 {
						ignoredReason = "on up trend up, new hight below sma and below threshold"
						return true
					}
				}
			}
		}
	}

	if lastBand.Candle.Close > float32(lastBand.Upper) && !isHasCrossUpper(result.Bands[:len(result.Bands)-1], false) {
		if midSecondLastBand.Candle.Open > float32(midSecondLastBand.Upper) && midSecondLastBand.Candle.Close < float32(midSecondLastBand.Upper) {
			shortHigestIndex := getHighestIndex(shortInterval.Bands)
			if midLastBand.Candle.Close < float32(midLastBand.Upper) && shortHigestIndex == len(shortInterval.Bands)-1 {
				ignoredReason = "long interval above upper just one and mid interval below sma and previous band down from upper"
				return true
			}
		}
	}

	if result.AllTrend.FirstTrendPercent > 20 && result.AllTrend.SecondTrendPercent > 20 && !isHasCrossUpper(result.Bands, false) {
		if result.AllTrend.FirstTrend != result.AllTrend.SecondTrend || (result.AllTrend.FirstTrend == models.TREND_SIDEWAY && result.AllTrend.SecondTrend == models.TREND_SIDEWAY) {
			if result.Position == models.ABOVE_SMA && hFirstPercentFromUpper < 3.4 {
				ignoredReason = "long interval above sma sideway and percent from upper below threshold"
				return true
			}

			if result.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.Position == models.ABOVE_UPPER && result.Position == models.BELOW_SMA && hFirstPercentFromSMA < 3.3 {
				ignoredReason = "mid interval above upper and long interval cross sma sideway"
				return true
			}
		}
	}

	if isHasCrossLower(result.Bands[len(result.Bands)-4:], false) && hFirstPercentFromUpper < 3.4 {
		ignoredReason = "last 4 band cross lower && margin from upper below threshold"
		return true
	}

	if countCrossSMA(result.Bands[len(result.Bands)-3:len(result.Bands)-1]) == 2 && countCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:len(midInterval.Bands)-1]) == 2 {
		ignoredReason = "two band mid interval cross upper and two band long interval cross sma"
		return true
	}

	if result.AllTrend.SecondTrendPercent < 12 && countBelowSMA(result.Bands[len(result.Bands)-6:], true) == 6 {
		if !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) && CalculateTrendShort(result.Bands[len(result.Bands)-6:]) != models.TREND_UP {
			if result.AllTrend.FirstTrend != models.TREND_DOWN || result.AllTrend.SecondTrend != models.TREND_DOWN {
				if !reversal2nd(*midInterval) {
					ignoredReason = "significan down, last 6 below sma but not cross lower"
					return true
				}
			}
		}
	}

	if midInterval.Position == models.BELOW_SMA && result.Position == models.BELOW_SMA && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if countBelowSMA(midInterval.Bands, false) == len(midInterval.Bands) && countBelowSMA(result.Bands[len(result.Bands)-5:], false) > 2 {
			if !isHasCrossLower(result.Bands[len(result.Bands)-3:], false) && !reversal2nd((*midInterval)) {
				ignoredReason = "mid interval all band below sma, long interval below sma but not cross lower"
				return true
			}
		}
	}

	if shortInterval.Position == models.ABOVE_UPPER && shortInterval.AllTrend.Trend == models.TREND_UP {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && !isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], false) {
			if midHFirstPercentFromUpper < 3.3 && result.AllTrend.SecondTrend == models.TREND_DOWN {
				ignoredReason = "short interval trend up above upper, mid above sma not cross lower and margin < 3"
				return true
			}
		}
	}

	if countDownBand(result.Bands[len(result.Bands)-6:len(result.Bands)-1]) == 5 && CalculateTrendShort(result.Bands[len(result.Bands)-3:]) == models.TREND_DOWN {
		ignoredReason = "squential down 5 band, and short trend not up"
		return true
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.Position == models.BELOW_SMA {
		if isHasDoji(result.Bands[len(result.Bands)-3:]) || isHasHammer(result.Bands[len(result.Bands)-3:]) {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) && isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-3:], true) {
				ignoredReason = "trend down, below sma, doji or hammer pattern"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA {
		if isHasCrossUpper(result.Bands[:len(result.Bands)/2], true) {
			if countCrossUpper(shortInterval.Bands) >= 9 && midInterval.Position == models.ABOVE_UPPER {
				if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.PriceChanges > 5 {
					ignoredReason = "up down, mid and short above upper"
					return true
				}
			}
		}

		if isHasCrossLower(result.Bands, false) && !isHasCrossUpper(result.Bands, true) && hFirstPercentFromUpper < 3.4 {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:], true) {
				ignoredReason = "below sma, margin from upper below 3"
				return true
			}
		}

		if result.AllTrend.ShortTrend == models.TREND_DOWN && isHasOpenCloseAboveUpper(result.Bands[len(result.Bands)-7:]) {
			if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				if midInterval.AllTrend.FirstTrendPercent > 20 && midInterval.AllTrend.SecondTrendPercent > 20 && midInterval.Position == models.BELOW_SMA {
					if isHasCrossLower(midInterval.Bands[len(result.Bands)/2:], false) && midHFirstPercentFromSMA < 3.2 {
						ignoredReason = "mid interval from upper down to lower but not significan"
						return true
					}
				}
			}
		}

		if lastBand.Candle.Open < float32(lastBand.SMA) && midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Close > float32(midLastBand.SMA) {
			if shortLastBand.Candle.Open < float32(shortLastBand.Upper) && shortLastBand.Candle.Hight > float32(shortLastBand.Upper) {
				ignoredReason = "long interval cross sma, mid interval cross sma and short interval cross upper"
				return true
			}
		}

		if lastBand.Candle.Low < float32(lastBand.SMA) && result.AllTrend.SecondTrend == models.TREND_DOWN && headMoreThan30PrecentToBody(lastBand) {
			if midLastBand.Candle.Hight > float32(midLastBand.Upper) && headMoreThan30PrecentToBody(midLastBand) {
				if shortInterval.Position == models.ABOVE_UPPER && getHighestIndex(shortInterval.Bands) == len(shortInterval.Bands)-1 && headMoreThanBody(shortLastBand) {
					ignoredReason = "head more than 30 on long and mid, short interval head more than body"
					return true
				}
			}
		}
	}

	if lastBand.Candle.Low < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA) {
		if shortInterval.Position == models.ABOVE_UPPER && countCrossUpper(shortInterval.Bands[len(shortInterval.Bands)/2:]) > 5 {
			midAboveUpperAndJustOne := (midInterval.Position == models.ABOVE_UPPER && countCrossUpper(midInterval.Bands[len(midInterval.Bands)-5:]) == 1)
			if (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) || midAboveUpperAndJustOne {
				ignoredReason = "cross sma and short cross upper and mid margin below 3"
				return true
			}
		}
	}

	if result.Position == models.BELOW_SMA && result.AllTrend.Trend == models.TREND_DOWN && countBelowSMA(result.Bands[len(result.Bands)-7:], false) == 7 {
		if shortLastBand.Candle.Close > float32(shortLastBand.SMA) && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-7:], true) {
			if hFirstPercentFromSMA < 3.3 {
				ignoredReason = "down trend last 7 band below SMA and percent from sma < 3"
				return true
			}
		}

		if midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2 {
			if midLastBand.Candle.Open < float32(midLastBand.Lower) || midSecondLastBand.Candle.Open < float32(midSecondLastBand.Lower) {
				ignoredReason = "down trend last 7 band below SMA and mid percent from sma < 3"
				return true
			}
		}

		if midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3 {
			if midLastBand.Candle.Open < float32(midLastBand.Lower) || midSecondLastBand.Candle.Open < float32(midSecondLastBand.Lower) {
				ignoredReason = "down trend last 7 band below SMA and mid percent from upper < 3"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && hFirstPercentFromUpper < 3.4 {
		if (result.AllTrend.FirstTrend != models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP) || (getHighestIndex(result.Bands) == len(result.Bands)-1) {
			if shortLastBand.Candle.Open < float32(shortLastBand.Upper) && shortLastBand.Candle.Close > float32(shortLastBand.Upper) {
				if midLastBand.Candle.Open < float32(midLastBand.Upper) && midLastBand.Candle.Close > float32(midLastBand.Upper) {
					ignoredReason = "above sma and percent from upper below 3"
					return true
				}
			}
		}
	}

	if lastBand.Candle.Open < float32(lastBand.Upper) && lastBand.Candle.Hight > float32(lastBand.Upper) {
		if !isHasCrossUpper(result.Bands[len(result.Bands)-6:len(result.Bands)-1], true) {
			if midLastBand.Candle.Hight < float32(midLastBand.Upper) && countDownBand(midInterval.Bands[len(midInterval.Bands)-3:]) >= 1 {
				if countCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:]) >= 1 {
					if shortLastBand.Candle.Close < float32(shortLastBand.Upper) && !isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-5:], false) {
						downBand := countDownBand(shortInterval.Bands[len(shortInterval.Bands)/2:])
						upBand := CountUpBand(shortInterval.Bands[len(shortInterval.Bands)/2:])
						if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)/2:], false) && downBand > upBand {
							ignoredReason = "already cross upper and mid, short down from upper"
							return true
						}
					}
				}
			}
		}

		if midLastBand.Candle.Hight > float32(midLastBand.Upper) && !isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-7:len(midInterval.Bands)-1], true) {
			if shortLastBand.Candle.Hight > float32(shortLastBand.Upper) && shortInterval.AllTrend.FirstTrend == models.TREND_UP {
				if shortInterval.AllTrend.SecondTrend == models.TREND_UP {
					if getHighestIndex(shortInterval.Bands) == len(shortInterval.Bands)-1 && getHighestIndex(midInterval.Bands) == len(midInterval.Bands)-1 {
						ignoredReason = "all band cross upper"
						return true
					}
				}
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP && getHighestIndex(result.Bands) == len(result.Bands)-1 {
		if result.Position == models.ABOVE_SMA && midInterval.Position == models.ABOVE_SMA && shortInterval.Position == models.ABOVE_SMA {
			if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], true) && hFirstPercentFromUpper < 3.4 && midHFirstPercentFromUpper < 3.3 && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend up up new hight"
				return true
			}
		}
	}

	if result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) && !isHasCrossUpper(result.Bands[len(result.Bands)/2:], true) {
		if midLastBand.Candle.Hight > float32(midLastBand.SMA) && midLastBand.Candle.Low < float32(midLastBand.SMA) {
			if !isHasCrossLower(midInterval.Bands[len(midInterval.Bands)/2:], false) && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], true) {
				ignoredReason = "sideway, mid after up then down"
				return true
			}
		}
	}

	if downFromUpper(*result) {
		if midInterval.Position == models.BELOW_SMA && isHasCrossLower(midInterval.Bands, false) && !isHasCrossUpper(midInterval.Bands, true) {
			if midHFirstPercentFromSMA < 3.2 && !upperLowerReversal(*shortInterval) && !isHasCrossLower(shortInterval.Bands[bandLen-1:], false) {
				ignoredReason = "down from upper, mid percent from sma Bellow 3"
				return true
			}
		}

		if countBelowSMA(midInterval.Bands[len(midInterval.Bands)/2:], false) > 0 && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], true) {
			if isHasCrossLower(shortInterval.Bands, false) && !isHasCrossUpper(shortInterval.Bands, true) && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "down from upper, short percent from upper Bellow 3"
				return true
			}
		}
	}

	if downFromUpperBelowSMA(*result) && midInterval.Position == models.BELOW_SMA {
		if midLastBand.Candle.Open < float32(midLastBand.Lower) && isHasCrossLower(midInterval.Bands, false) && !isHasCrossUpper(midInterval.Bands, true) && midHFirstPercentFromSMA < 3.2 {
			ignoredReason = "down from upper, mid percent from sma Bellow 3 2nd"
			return true
		}
	}

	if result.AllTrend.SecondTrend != models.TREND_UP {
		if shortInterval.Position == models.BELOW_SMA && countBelowSMA(shortInterval.Bands[len(shortInterval.Bands)/2:], false) >= len(shortInterval.Bands)/2 {
			shortHSecondPercentFromSMA := (shortLastBand.SMA - float64(shortLastBand.Candle.Close)) / float64(shortLastBand.Candle.Close) * 100
			if !isHasCrossSMA(shortInterval.Bands[len(shortInterval.Bands)/2:], true) && shortHSecondPercentFromSMA < .13 && midHFirstPercentFromUpper < 3.3 {
				ignoredReason = "trend down from upper"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && countAboveSMA(result.Bands[len(result.Bands)-5:]) == 5 && !isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) {
		if midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3 && hFirstPercentFromUpper < 3.4 {
			if shortInterval.Position == models.ABOVE_SMA && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-5:], true) {
				ignoredReason = "short, mid and long above sma but percent below 3"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && !isHasCrossUpper(result.Bands[len(result.Bands)-7:], false) && countAboveSMA(result.Bands[len(result.Bands)-7:]) == 7 {
		if secondLastBand.Candle.Hight > float32(secondLastBand.Upper) {
			if midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3 && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
				ignoredReason = "mid and long above sma but percent below 3"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && lastBand.Candle.Open < float32(lastBand.SMA) {
		if midInterval.PriceChanges > 3 && midLastBand.Candle.Close > float32(midLastBand.SMA) {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:], true) {
				ignoredReason = "percent chnages more than3, short and mid cross upper and long cross sma"
				return true
			}
		}
	}

	if downFromUpperAboveSMA(*result) {
		midLastBandCrossSMA := (midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Hight > float32(midLastBand.SMA))
		if midLastBandCrossSMA || (midLastBand.Candle.Open < float32(midLastBand.Upper) && midLastBand.Candle.Hight > float32(midLastBand.Upper)) {
			if shortLastBand.Candle.Open < float32(shortLastBand.Upper) && shortLastBand.Candle.Hight > float32(shortLastBand.Upper) {
				ignoredReason = "down from upper, mid and short cross upper"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN {
		midlowestIndex := getLowestIndex(midInterval.Bands)
		if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-2:], true) && midlowestIndex > len(midInterval.Bands)/2 {
			if isHasCrossUpper(shortInterval.Bands[len(result.Bands)-3:], true) {
				ignoredReason = "trend down mid second wave hit lowest and then cross upper"
				return true
			}
		}

		midHigestIndex := getHighestIndex(midInterval.Bands)
		if midInterval.Position == models.ABOVE_SMA && midlowestIndex > len(midInterval.Bands)/2 {
			checking := false
			if midHigestIndex == len(midInterval.Bands)-1 {
				if midHFirstPercentFromUpper < 3.3 {
					checking = true
				}
			} else {
				midPercentFromHigest := (midInterval.Bands[midHigestIndex].Candle.Close - midLastBand.Candle.Close) / midLastBand.Candle.Close * 100
				if midPercentFromHigest < 3.3 && midHFirstPercentFromUpper < 3.3 {
					checking = true
				}
			}

			shortLastBandIsHeighestAndPercentBelow3 := (getHighestIndex(shortInterval.Bands) == len(shortInterval.Bands)-1 && shortHSecondPercentFromUpper < 3.2)
			if checking && (isHasCrossUpper(shortInterval.Bands[len(result.Bands)-3:], true) || shortLastBandIsHeighestAndPercentBelow3) {
				ignoredReason = "trend down mid second wave hit lowest and then get higest"
				return true
			}
		}

		hightIndex := getHighestIndex(midInterval.Bands[len(midInterval.Bands)/2:]) + len(midInterval.Bands)/2
		if hightIndex >= len(midInterval.Bands)/2-3 && midInterval.Bands[hightIndex].Candle.Hight > float32(midInterval.Bands[hightIndex].Upper) {
			if countBelowSMA(midInterval.Bands[hightIndex:], false) > 0 && countBelowSMA(midInterval.Bands[:len(midInterval.Bands)/2], false) > 0 {
				if midHFirstPercentFromUpper < 3.3 {
					ignoredReason = "long interval down, mid reversal from sma but margin from upper < 3"
					return true
				}
			}
		}
	}

	if result.Position == models.ABOVE_UPPER || hFirstPercentFromUpper < 3.4 {
		if midInterval.AllTrend.FirstTrend != models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrendPercent > 20 {
			if midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Hight > float32(midLastBand.SMA) {
				ignoredReason = "mid on down, not significan and cross sma"
				return true
			}
		}
	}

	if lastBand.Candle.Close < float32(lastBand.SMA) && result.AllTrend.SecondTrend == models.TREND_SIDEWAY && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && countBelowLower(midInterval.Bands[len(midInterval.Bands)-5:], false) > 0 {
			if (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) || shortInterval.Position == models.ABOVE_UPPER {
				ignoredReason = "on trend down short margin from upper below 3"
				return true
			}
		}
	}

	if isHasCrossLower(result.Bands, false) {
		belowSMAAndMarginBelow3 := result.Position == models.BELOW_SMA && hFirstPercentFromSMA < 3.3
		lastCrossLower := getLastIndexCrossLower(result.Bands)
		if lastCrossLower >= 0 && ((lastBand.Candle.Open < float32(lastBand.SMA) && lastBand.Candle.Hight > float32(lastBand.SMA)) || belowSMAAndMarginBelow3) {
			if countCrossSMA(result.Bands[lastCrossLower:]) <= 1 && midLastBand.Candle.Hight > float32(midLastBand.Upper) {
				ignoredReason = "cross sma and mid cross upper"
				return true
			}
		}

		if lastCrossLower >= 0 && result.Position == models.BELOW_SMA && countCrossSMA(result.Bands[lastCrossLower:]) <= 1 {
			if midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < .43 {
				if (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) || shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "below sma percent below, mid above sma percent below 3 and short below 3 or cross upper"
					return true
				}
			}
		}

		if lastCrossLower >= 0 && result.Position == models.ABOVE_UPPER && countCrossUpper(result.Bands[lastCrossLower:]) <= 1 {
			if (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) || midInterval.Position == models.ABOVE_UPPER {
				if (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) || shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "above upper, mid above sma percent below 3 and short below 3 or cross upper"
					return true
				}
			}
		}
	}

	if secondLastBand.Candle.Hight > float32(secondLastBand.Upper) && secondLastBand.Candle.Close < float32(secondLastBand.Upper) {
		lastCrossLower := getLastIndexCrossLower(midInterval.Bands)
		if lastCrossLower >= 0 && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], true) {
			if midInterval.Position == models.BELOW_SMA && isHasCrossSMA(midInterval.Bands[lastCrossLower:], false) {
				lastCrossLower = getLastIndexCrossLower(shortInterval.Bands)
				if isHasCrossUpper(shortInterval.Bands, true) && shortInterval.Position == models.BELOW_SMA && isHasCrossSMA(shortInterval.Bands[lastCrossLower:], false) {
					ignoredReason = "below sma, down from sma"
					return true
				}
			}
		}
	}

	if result.Position == models.ABOVE_UPPER && midInterval.Position == models.ABOVE_UPPER && shortInterval.Position == models.ABOVE_UPPER {
		ignoredReason = "all interval above upper"
		return true
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_UP {
		if getHighestIndex(result.Bands[len(result.Bands)-7:]) != len(result.Bands[len(result.Bands)-7:])-1 {
			if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP {
				if midLastBand.Candle.Open < float32(midLastBand.SMA) && midLastBand.Candle.Close > float32(midLastBand.SMA) {
					if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-5:], true) {
						ignoredReason = "donw up, short has cross upperr"
						return true
					}
				}
			}
		}
	}

	if isHasCrossUpper(result.Bands[len(result.Bands)/2:], true) && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if countBelowSMA(midInterval.Bands, true) > 13 && !isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-13:], true) {
			if (shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1) || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
				ignoredReason = "short trend down, mid have many below sma"
				return true
			}
		}
	}

	if lastBand.Candle.Hight < float32(lastBand.SMA) && CountUpBand(result.Bands[len(result.Bands)-5:]) <= 1 {
		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend down, short above sma but percent bellow 3"
				return true
			}
		}
	}

	if afterUpThenDown(result.Bands) {
		if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:], true) && isHasCrossUpper(shortInterval.Bands[len(midInterval.Bands)-3:], true) {
			ignoredReason = "up down percent below 3 and mid, short cross upper"
			return true
		}
	}

	midHigest := getHighestIndex(midInterval.Bands)
	midPercentFromHigest := (midInterval.Bands[midHigest].Candle.Close - midLastBand.Candle.Close) / midLastBand.Candle.Close * 100
	if lastBand.Candle.Close > float32(lastBand.Upper) || hFirstPercentFromUpper < 3.4 {
		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend != models.TREND_DOWN {
			if midPercentFromHigest < 3.3 && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) {
				ignoredReason = "long and short cross upper and mid percent below 3"
				return true
			}
		}
	}

	if lastBand.Candle.Hight < float32(lastBand.Upper) && countAboveUpper(result.Bands[len(result.Bands)-5:]) > 0 {
		if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.SecondTrendPercent < 15 {
			if !isHasCrossUpper(midInterval.Bands, false) && countBelowSMA(midInterval.Bands, false) == 0 && midHFirstPercentFromUpper < 3.2 {
				if shortInterval.Position == models.ABOVE_UPPER || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
					if countBelowSMA(shortInterval.Bands[len(shortInterval.Bands)/2:], false) > 0 {
						ignoredReason = "on upper, mid and short percent below 3"
						return true
					}
				}
			}
		}
	}

	if isHasCrossUpper(result.Bands[len(result.Bands)-5:], true) && result.AllTrend.ShortTrend != models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP {
		if midInterval.Position == models.ABOVE_UPPER || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
			if shortInterval.Position == models.ABOVE_UPPER || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
				ignoredReason = "on upper, mid and short percent below 3 2nd"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && secondLastBand.Candle.Open > float32(secondLastBand.Upper) && secondLastBand.Candle.Close < float32(secondLastBand.Upper) {
		if countAboveUpper(midInterval.Bands[len(midInterval.Bands)-5:]) > 0 && CalculateTrendShort(midInterval.Bands[len(midInterval.Bands)-5:]) != models.TREND_UP {
			if shortInterval.AllTrend.ShortTrend != models.TREND_UP || shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "down from upper"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.ShortTrend != models.TREND_UP && getHighestIndex(result.Bands) > len(result.Bands)-7 {
		if isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-5:], false) && midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3 {
			if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "short trend down, mid and short percent from upper below 3"
				return true
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP && result.Position == models.ABOVE_SMA {
		if result.AllTrend.ShortTrend != models.TREND_UP || CalculateTrendShort(result.Bands[len(result.Bands)-5:len(result.Bands)-1]) != models.TREND_UP {
			if midInterval.Position == models.ABOVE_SMA || (midInterval.Position == models.BELOW_SMA && midInterval.AllTrend.SecondTrend == models.TREND_DOWN) {
				if shortInterval.Position == models.ABOVE_UPPER || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
					ignoredReason = "on upper, mid and short percent below 3 3nd"
					return true
				}
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.FirstTrendPercent < 10 && result.AllTrend.SecondTrend != models.TREND_UP {
		if !isHasCrossUpper(result.Bands[len(result.Bands)-7:], true) && !isHasCrossSMA(result.Bands[len(result.Bands)-7:], false) {
			if midInterval.Position == models.ABOVE_UPPER && countBelowSMA(midInterval.Bands[len(midInterval.Bands)-7:], false) == 0 {
				if shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "mid and short above upper"
					return true
				}
			}
		}
	}

	longHigestIndex := getHighestIndex(result.Bands)
	if result.Position == models.ABOVE_SMA && longHigestIndex > len(result.Bands)/2 {
		if isHasCrossSMA(result.Bands[longHigestIndex:], false) && percentFromHeight < 3.4 {
			if midInterval.Position == models.ABOVE_UPPER || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
				if shortInterval.Position == models.ABOVE_UPPER || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
					ignoredReason = "above sma after down and percent from high below 3"
					return true
				}
			}
		}
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.ShortTrend == models.TREND_UP && countAboveUpper(result.Bands[len(result.Bands)-7:]) > 0 {
		if midInterval.Position == models.ABOVE_UPPER || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
			if shortInterval.Position == models.ABOVE_UPPER || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
				ignoredReason = "above sma after down and percent from high below 3 2nd"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.Trend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.FirstTrend == models.TREND_UP && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if shortInterval.AllTrend.FirstTrend == models.TREND_DOWN && shortInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				if midSecondLastBand.Candle.Open > midSecondLastBand.Candle.Close && (shortSecondLastBand.Candle.Open > shortSecondLastBand.Candle.Close || shortHSecondPercentFromSMA < 3.1) {
					ignoredReason = "all band start to down, up just one "
					return true
				}
			}
		}
	}

	previousBandPercentChanges := (secondLastBand.Candle.Open - secondLastBand.Candle.Close) / secondLastBand.Candle.Close * 100
	if result.AllTrend.ShortTrend != models.TREND_UP && (result.PriceChanges > 4 || previousBandPercentChanges > 4) {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "just significan down, short above sma but percent below 3 "
				return true
			}
		}

		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && shortInterval.AllTrend.Trend == models.TREND_DOWN {
			if countAboveUpper(midInterval.Bands[len(midInterval.Bands)-7:]) > 0 {
				if (shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1) || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
					ignoredReason = "just significan down, short above sma but percent below 3 2nd"
					return true
				}
			}
		}

		if result.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.PriceChanges > 4 && shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "just significan down, short above sma but percent below 3 3nd"
				return true
			}
		}

		if midInterval.AllTrend.ShortTrend != models.TREND_UP && midInterval.PriceChanges > 4 {
			if (shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1) || (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) {
				ignoredReason = "just significan down, short above sma but percent below 3 4nd"
				return true
			}
		}

		if isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-2:], true) || countBelowLower(midInterval.Bands[len(midInterval.Bands)-3:], false) > 0 {
			if countBelowLower(shortInterval.Bands[len(shortInterval.Bands)-3:], false) > 0 {
				ignoredReason = "just significan down, short have open close below lower"
				return true
			}
		}

		if midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2 && shortHSecondPercentFromUpper < 3.2 {
			if (midLastBand.Candle.Low > float32(midLastBand.Lower) || midHFirstPercentFromUpper < 3.3) && !isHasCrossLower(shortInterval.Bands[bandLen-2:], false) {
				ignoredReason = "just significan down, mid below sma and percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.SecondTrendPercent < 15 {
		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.Position == models.BELOW_SMA {
			if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 && midHFirstPercentFromUpper < 3.3 && !reversal2nd(*midInterval) {
				ignoredReason = "trend down, mid below sma but percent below 3 "
				return true
			}
		}
	}

	if isHasCrossSMA(result.Bands[len(result.Bands)-2:], false) || isHasCrossUpper(result.Bands[len(result.Bands)-2:], true) {
		if midHFirstPercentFromUpper < 3.3 && shortHSecondPercentFromUpper < 3.2 {
			ignoredReason = "all interval cross upper 2nd "
			return true
		}
	}

	higestHigest := getHighestHightIndex(result.Bands)
	if result.AllTrend.ShortTrend == models.TREND_DOWN && higestHigest != len(result.Bands)-1 && higestHigest > len(result.Bands)-5 {
		percent := (result.Bands[higestHigest].Candle.Hight - lastBand.Candle.Close) / lastBand.Candle.Close * 100
		if percent > 10 && midInterval.AllTrend.SecondTrend == models.TREND_DOWN && !isHasCrossLower(midInterval.Bands[len(midInterval.Bands)/2:], false) {
			ignoredReason = "down after significan up "
			return true
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_SIDEWAY && countAboveSMA(result.Bands[len(result.Bands)/2:]) == 0 {
		if result.Position == models.BELOW_SMA && hFirstPercentFromSMA < 3.3 {
			if (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
				shortAboveSMAorAboveUpper := (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) || shortInterval.Position == models.ABOVE_UPPER
				if (shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1) || shortAboveSMAorAboveUpper {
					ignoredReason = "sideway, short above sma but percent below 3"
					return true
				}
			}
		}
	}

	if result.Position == models.BELOW_SMA && countCrossLower(result.Bands[len(result.Bands)/2:]) > 1 && countCrossSMA(result.Bands[len(result.Bands)/2:]) > 0 {
		if hFirstPercentFromSMA < 3.3 && midHFirstPercentFromUpper < 3.3 && shortHSecondPercentFromUpper < 3.2 {
			ignoredReason = "up down sideway and percent below 3"
			return true
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP && isHasCrossUpper(result.Bands[len(result.Bands)-2:], true) {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-4:], true) {
			if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-4:], false) && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend up up and down from upper"
				return true
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
			if midInterval.AllTrend.ShortTrend == models.TREND_UP && midInterval.PriceChanges > 3.5 {
				if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-2:], false) || isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-2:], true) {
					ignoredReason = "significan up after down"
					return true
				}
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend != models.TREND_UP {
		if shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.PriceChanges > 3.5 {
			if isHasCrossSMA(shortInterval.Bands[len(shortInterval.Bands)-2:], false) || isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-2:], true) {
				ignoredReason = "significan up on trend down"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && midInterval.Position == models.BELOW_SMA {
		if secondLastBand.Candle.Open < float32(secondLastBand.Lower) || secondLastBand.Candle.Close < float32(secondLastBand.Lower) {
			if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && countAboveSMA(midInterval.Bands) == 0 && midHFirstPercentFromSMA < 3.2 {
				ignoredReason = "trend down. short below sma and percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend != models.TREND_UP {
		midAboveSMAOrAboveUpper := midInterval.Position == models.ABOVE_UPPER || (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2)
		if midAboveSMAOrAboveUpper || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
			if midLastBand.Candle.Low > float32(midLastBand.Lower) || midHFirstPercentFromUpper < 3.3 {
				if shortHSecondPercentFromUpper < 3.2 {
					ignoredReason = "trend down. short sma percent below 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if isHasCrossSMA(shortInterval.Bands[len(shortInterval.Bands)-2:], false) || isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-2:], true) {
			if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend down. short sma percent below 3 2nd"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && result.Position == models.ABOVE_SMA {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend != models.TREND_UP {
			if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-4:], true) && midInterval.Position == models.ABOVE_SMA {
				if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.AllTrend.ShortTrend == models.TREND_UP {
					if isHasCrossSMA(shortInterval.Bands[len(shortInterval.Bands)-2:], false) {
						ignoredReason = "trend up but band start to down"
						return true
					}
				}
			}
		}
	}

	if shortInterval.AllTrend.Trend != models.TREND_UP && shortLastBand.Candle.Close < float32(shortLastBand.SMA) && !isHasCrossLower(shortInterval.Bands[bandLen-2:], false) {
		midNotCrossSMAorShortTrendNotUp := (!isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-4:], false) || midInterval.AllTrend.ShortTrend != models.TREND_UP)
		if midInterval.AllTrend.SecondTrend != models.TREND_UP && midNotCrossSMAorShortTrendNotUp && !isHasCrossLower(midInterval.Bands[bandLen-4:], false) {
			if result.AllTrend.ShortTrend != models.TREND_UP && isHasCrossUpper(result.Bands[len(result.Bands)-4:], true) && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "starting down, short below sma, percent below 3"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if result.AllTrend.ShortTrend != models.TREND_DOWN && result.PriceChanges > 4 {
			if midHFirstPercentFromUpper < 3.3 && (shortInterval.Position == models.ABOVE_UPPER || shortHSecondPercentFromUpper < 3.2) {
				ignoredReason = "already up above 4% and mid percent from upper below 3"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && result.AllTrend.FirstTrend == models.TREND_UP && result.AllTrend.SecondTrend == models.TREND_UP {
		if midInterval.Position == models.ABOVE_UPPER || midHFirstPercentFromUpper < 3.3 {
			if shortInterval.Position == models.ABOVE_UPPER || shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "up trend, above sma and mid and short percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN && result.Position == models.BELOW_SMA {
		if midInterval.AllTrend.Trend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) || (midInterval.Position == models.ABOVE_UPPER && midHFirstPercentFromUpper < 3.3) {
				if shortInterval.Position == models.ABOVE_UPPER || shortHSecondPercentFromUpper < 3.2 {
					ignoredReason = "down trend, short and mid percent below 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_UP {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-2:], false) || (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) {
				if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-2:], true) {
					ignoredReason = "down trend, mid cross sma, short cross upper"
					return true
				}
			}
		}

		if result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands[len(result.Bands)/2:], false) {
			if countBelowLower(midInterval.Bands[len(midInterval.Bands)-3:], false) > 0 && midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "down trend, mid contain open close below lower, short cross percent below 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_UP && result.Direction == BAND_DOWN {
		if midInterval.AllTrend.FirstTrend == models.TREND_UP && midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "up trend, start down, short below sma percent below 3"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.Position == models.BELOW_SMA {
				if isHasCrossLower(shortInterval.Bands[bandLen-7:], false) && shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "tren up but band down"
					return true
				}
			}
		}

		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.AllTrend.ShortTrend != models.TREND_UP {
				if !isHasCrossLower(shortInterval.Bands[bandLen/2:], false) {
					ignoredReason = "band down below sma but not cross lower"
					return true
				}
			}

			if shortInterval.Position == models.BELOW_SMA {
				if shortHSecondPercentFromSMA < 3.1 && midHFirstPercentFromUpper < 3.3 {
					ignoredReason = "band down below and percent below 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP {
		if secondLastBand.Candle.Open > secondLastBand.Candle.Close && isHasCrossSMA(result.Bands[len(result.Bands)-4:], false) {
			if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.Position == models.ABOVE_SMA {
				if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], true) && countBelowSMA(midInterval.Bands[len(midInterval.Bands)/2:], false) == 0 {
					shortBelowSMAAndPercentBelow3 := shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3
					if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-2:], false) && shortInterval.AllTrend.Trend == models.TREND_DOWN && shortBelowSMAAndPercentBelow3 {
						ignoredReason = "pokok e trend down"
						return true
					}
				}
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN && result.PriceChanges > 3 && (!isHasCrossLower(shortInterval.Bands[bandLen-1:], false) || !isHasCrossLower(midInterval.Bands[bandLen-1:], false) || !(result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_UP)) {
		if midInterval.Position == models.BELOW_SMA && isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-2:], false) {
			if midHFirstPercentFromSMA < 3.2 && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "short trend down. and percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.AllTrend.ShortTrend != models.TREND_UP {
				if !isHasCrossLower(shortInterval.Bands[bandLen-1:], false) || !isHasCrossLower(midInterval.Bands[bandLen-1:], false) || !(result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_UP) {
					ignoredReason = "all interval short trend down"
					return true
				}
			}

			if result.PriceChanges > 4 {
				if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "along and mid on short trend down but short percent from sma below 3"
					return true
				}
			}
		}

		if isHasBelowLower(midInterval.Bands[len(midInterval.Bands)-2:]) {
			ignoredReason = "mid contain open close below lower"
			return true
		}

		if isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-2:], true) || midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-2:], true) && shortHSecondPercentFromUpper < 3.2 && shortInterval.AllTrend.ShortTrend != models.TREND_UP {
				ignoredReason = "mid and short cross lower"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.Position == models.BELOW_SMA && isHasCrossLower(result.Bands[len(result.Bands)-2:], false) {
		if midInterval.Position == models.BELOW_SMA && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if midHFirstPercentFromSMA < 3.2 && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "still down, on lower but mid and short percent below 3"
				return true
			}
		}
	}

	crossLowerTrendDown := isHasCrossLower(result.Bands[bandLen-2:], false) || CalculateTrendShort(result.Bands[bandLen-5:bandLen-1]) == models.TREND_DOWN
	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.Position == models.BELOW_SMA && crossLowerTrendDown {
		if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) || isHasCrossUpper(midInterval.Bands[bandLen-2:], true) {
			if countAboveUpper(shortInterval.Bands[bandLen-4:]) > 0 || isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) {
				ignoredReason = "down, still cross upper and short cross upper"
				return true
			}
		}
	}

	if result.Position == models.BELOW_SMA && (result.Direction == BAND_DOWN || (secondLastBand.Candle.Open > secondLastBand.Candle.Close)) {
		if countAboveSMA(result.Bands[5:]) == 0 && isHasCrossSMA(result.Bands[len(result.Bands)-3:], false) {
			if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && shortInterval.AllTrend.ShortTrend == models.TREND_DOWN {
				ignoredReason = "trend down, resisten after cross sma"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP && result.Position == models.ABOVE_SMA {
		if isHasCrossSMA(result.Bands[len(result.Bands)-2:], false) && isHasCrossUpper(result.Bands[len(result.Bands)-4:], true) {
			if (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) || midHFirstPercentFromSMA < 3.2 || shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "up trend but start to down"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_DOWN {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)/2:], true) {
			if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-3:], false) && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend down from upper "
				return true
			}
		}
	}

	if result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if result.AllTrend.ShortTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.AllTrend.ShortTrend == models.TREND_UP {
			if result.PriceChanges > 3 && isHasCrossUpper(midInterval.Bands[bandLen-2:], false) && isHasCrossUpper(shortInterval.Bands[bandLen-2:], false) {
				ignoredReason = "down trend but short and mid cross upper "
				return true
			}
		}

		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend != models.TREND_DOWN {
			if midLastBand.Candle.Hight > float32(midLastBand.Upper) && countCrossUpper(midInterval.Bands[len(midInterval.Bands)-5:]) == 1 {
				if isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-5:], false) {
					ignoredReason = "mid interval cross uppper"
					return true
				}
			}
		}
	}

	if lastBand.Candle.Open > lastBand.Candle.Close || secondLastBand.Candle.Open > secondLastBand.Candle.Close {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && shortInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if !isHasCrossLower(shortInterval.Bands[bandLen-1:], false) || !isHasCrossLower(midInterval.Bands[bandLen-1:], false) || !(result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_UP) {
				ignoredReason = "mid and short, short trend down "
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend != models.TREND_UP && result.PriceChanges > 4 {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.PriceChanges > 4 {
			if shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.Position == models.BELOW_SMA && shortInterval.PriceChanges > 1.4 {
				ignoredReason = "short tren down, mid short trend down. price change more than 1.4"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if shortInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
				ignoredReason = "short tren down, mid short, short trend down"
				return true
			}
		}

		if midInterval.Position == models.BELOW_SMA && midInterval.AllTrend.SecondTrend != models.TREND_UP {
			if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "short tren down, below sma not up and short above sma percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if shortInterval.AllTrend.ShortTrend == models.TREND_DOWN {
				if isHasCrossSMA(shortInterval.Bands[len(shortInterval.Bands)-5:], false) {
					if shortHSecondPercentFromUpper < 3.2 {
						ignoredReason = "mid side way but short percent from upper below 3"
						return true
					}
				} else {
					if shortHSecondPercentFromSMA < 3.1 {
						ignoredReason = "mid side way but short percent from sma below 3"
						return true
					}
				}
			}
		}

		if result.AllTrend.Trend == models.TREND_DOWN && lastBand.Candle.Close < float32(lastBand.SMA) {
			if midInterval.Position == models.BELOW_SMA && countAboveSMA(midInterval.Bands[len(midInterval.Bands)/2:]) == 0 && midHFirstPercentFromSMA < 3.2 {
				if shortInterval.AllTrend.Trend == models.TREND_DOWN && shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend down, short percent from sma below 3"
					return true
				}
			}

			if countDownBand(midInterval.Bands[len(midInterval.Bands)-4:len(midInterval.Bands)-1]) == 3 {
				if midInterval.Direction == BAND_DOWN {
					ignoredReason = "trend down, mid interval band down"
					return true
				}

				if shortHSecondPercentFromSMA < 3.2 {
					ignoredReason = "trend down, short percent from sma below 3"
					return true
				}
			}

			if midInterval.AllTrend.Trend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend != models.TREND_UP {
				if midInterval.PriceChanges > 4 && (shortHSecondPercentFromSMA < 3.1 || !isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-3:], false)) {
					ignoredReason = "trend down, bouncing minor up"
					return true
				}
			}

			lastBandOpenOrCloseBelowLower := lastBand.Candle.Open < float32(lastBand.Lower) || lastBand.Candle.Close < float32(lastBand.Lower)
			if lastBandOpenOrCloseBelowLower || secondLastBand.Candle.Open < float32(secondLastBand.Lower) || secondLastBand.Candle.Close < float32(secondLastBand.Lower) {
				if shortHSecondPercentFromSMA < 3.1 || !isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
					ignoredReason = "trend down, below lower"
					return true
				}
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && isHasCrossUpper(result.Bands[len(result.Bands)-2:], true) {
		if midInterval.AllTrend.ShortTrend != models.TREND_UP && countDownBand(midInterval.Bands[len(midInterval.Bands)-2:]) > 0 {
			if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:], true) {
				if shortInterval.AllTrend.ShortTrend != models.TREND_UP && !isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
					ignoredReason = "cross upper and starting to down"
					return true
				}
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && lastBand.Candle.Close < float32(lastBand.SMA) {
		if midInterval.AllTrend.ShortTrend != models.TREND_UP && countDownBand(midInterval.Bands[len(midInterval.Bands)-2:]) > 0 {
			if isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-3:], true) || isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-3:], false) {
				if shortInterval.AllTrend.ShortTrend != models.TREND_UP && !isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
					ignoredReason = "down, minor up  then starting to down"
					return true
				}
			}
		}

		if result.AllTrend.ShortTrend != models.TREND_UP {
			if midInterval.AllTrend.ShortTrend == models.TREND_UP && midInterval.PriceChanges > 3 {
				if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) || shortHSecondPercentFromUpper < 3.2 {
					ignoredReason = "down, start up but short already cross upper"
					return true
				}
			}
		}

		if midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if midInterval.AllTrend.FirstTrend == models.TREND_SIDEWAY && midInterval.AllTrend.SecondTrend == models.TREND_SIDEWAY {
				if shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2 {
					ignoredReason = "mid sideway and short above sma but percent below 3"
					return true
				}
			}

			if (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) || midInterval.Position == models.ABOVE_UPPER {
				if (shortInterval.Position == models.ABOVE_SMA && shortHSecondPercentFromUpper < 3.2) || shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "mid and short above sma but percent below 3"
					return true
				}
			}
		}
	}

	shortHigestHalfWave := getHigestPrice(shortInterval.Bands[len(shortInterval.Bands)/2:])
	shortPercentFromHalfHight := (shortHigestHalfWave - shortLastBand.Candle.Close) / shortLastBand.Candle.Close * 100
	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_UP && result.PriceChanges > 5 {
		if midInterval.AllTrend.Trend == models.TREND_UP && midInterval.AllTrend.ShortTrend != models.TREND_UP {
			if shortPercentFromHalfHight < 3.2 {
				ignoredReason = "just down after up and percent from last hight < 3"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend != models.TREND_UP {
		if midInterval.AllTrend.Trend != models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.Position == models.BELOW_SMA && isHasCrossLower(shortInterval.Bands[len(shortInterval.Bands)-4:], false) && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "down, start up reverlsal"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_UP && result.PriceChanges > 5 {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.Position == models.ABOVE_SMA {
			if shortInterval.Position == models.BELOW_SMA && shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.PriceChanges > 1.5 {
				ignoredReason = "down, start up reverlsal 2nd"
				return true
			}
		}
	}

	if result.Position == models.BELOW_SMA && result.AllTrend.Trend != models.TREND_UP {
		if isHasCrossSMA(midInterval.Bands[len(midInterval.Bands)-1:], false) && midInterval.AllTrend.ShortTrend == models.TREND_UP && midInterval.PriceChanges > 3 {
			if isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], false) {
				ignoredReason = "short cross upper"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP && countDownBand(result.Bands[len(result.Bands)-2:]) > 0 {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossUpper(midInterval.Bands[len(midInterval.Bands)-4:], true) {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "starting down"
				return true
			}
		}

		if midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2 && !isHasCrossLower(midInterval.Bands[len(midInterval.Bands)-3:], false) {
			if shortInterval.AllTrend.SecondTrend != models.TREND_UP && shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "starting down, mid below sma up percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_UP && result.PriceChanges > 10 {
		if midInterval.AllTrend.ShortTrend == models.TREND_UP && midInterval.PriceChanges > 4 {
			if shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.PriceChanges > 3 && isHasCrossUpper(shortInterval.Bands[len(shortInterval.Bands)-3:], true) {
				ignoredReason = "all trend up, and short cross upper"
				return true
			}
		}

		if countAboveUpper(midInterval.Bands[bandLen-7:]) > 0 && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if !isHasCrossUpper(shortInterval.Bands[bandLen-4:], false) && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "all trend up, and mid contain open close above upper"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_DOWN && result.Position == models.BELOW_SMA {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.Position == models.ABOVE_SMA {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				ignoredReason = "short trend up, and start down"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_UP && result.PriceChanges > 6 && lastBand.Candle.Hight < float32(lastBand.Upper) && isHasCrossUpper(result.Bands[bandLen-4:], true) {
		if midInterval.Position == models.ABOVE_SMA && midInterval.AllTrend.ShortTrend == models.TREND_UP && midHFirstPercentFromUpper < 3.2 {
			if isHasCrossUpper(midInterval.Bands[bandLen/2:], true) && !isHasCrossSMA(midInterval.Bands[bandLen/2:], false) {
				if shortInterval.Position == models.ABOVE_UPPER {
					ignoredReason = "short on upper and mid percent below 3"
					return true
				}

			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_UP && lastBandHeadDoubleBody(result) && result.Position == models.ABOVE_SMA {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.Position == models.ABOVE_SMA && isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
			if shortInterval.Position == models.BELOW_SMA && shortInterval.AllTrend.ShortTrend == models.TREND_DOWN && !isHasCrossLower(shortInterval.Bands[bandLen-3:], false) {
				ignoredReason = "starting down, down down"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN && result.Direction == BAND_DOWN && result.Position == models.ABOVE_SMA {
		if midInterval.Position == models.BELOW_SMA && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY && isHasCrossLower(midInterval.Bands[bandLen-4:], true) {
			if shortInterval.AllTrend.SecondTrend == models.TREND_SIDEWAY && shortInterval.AllTrend.ShortTrend == models.TREND_UP {
				if shortLastBand.Candle.Low < float32(shortLastBand.SMA) && shortLastBand.Candle.Hight > float32(shortLastBand.SMA) {
					ignoredReason = "just minor up"
					return true
				}
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) || (shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1) {
				ignoredReason = "just minor up still trend down"
				return true
			}
		}

		if result.AllTrend.SecondTrend == models.TREND_UP {
			if midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2 {
				if shortHSecondPercentFromUpper < 3.2 {
					if !isHasCrossLower(shortInterval.Bands[bandLen-1:], false) || !isHasCrossLower(midInterval.Bands[bandLen-1:], false) || !(result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_UP) {
						ignoredReason = "just minor up still trend down 2nd"
						return true
					}
				}
			}
		}
	}

	if result.AllTrend.ShortTrend != models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_DOWN {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) || (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) {
				if shortHSecondPercentFromUpper < 3.2 {
					ignoredReason = "percent from upper < 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_UP {
		if result.Direction == models.TREND_DOWN && result.PriceChanges > 4 {
			if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
				if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.Position == models.BELOW_SMA {
					if shortHSecondPercentFromSMA < 3.1 && !isHasCrossLower(shortInterval.Bands[bandLen-2:], false) {
						ignoredReason = "trend down percent from upper < 3"
						return true
					}
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP {
		if midInterval.AllTrend.FirstTrend == models.TREND_UP && midInterval.AllTrend.SecondTrend == models.TREND_UP && midHFirstPercentFromUpper < 3.3 {
			if isHasCrossUpper(shortInterval.Bands[bandLen-2:], true) {
				ignoredReason = "trend up, mid percent from upper below 3, and short cross upper"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_SIDEWAY {
		if (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
			if shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "short trend side way and percent from upper below 3"
				return true
			}
		}

		if result.Position == models.ABOVE_SMA && lastBand.Candle.Hight < float32(lastBand.Upper) && isHasCrossUpper(result.Bands[bandLen-4:], true) {
			if midInterval.Position == models.ABOVE_SMA && midInterval.AllTrend.ShortTrend == models.TREND_UP && midInterval.PriceChanges > 3 {
				if isHasCrossUpper(midInterval.Bands[bandLen/2:], true) && !isHasCrossSMA(midInterval.Bands[bandLen/2:], false) {
					if shortInterval.Position == models.ABOVE_UPPER {
						ignoredReason = "short on upper and mid percent more than 3"
						return true
					}
				}
			}
		}

		if isHasCrossUpper(result.Bands[bandLen-2:], true) || isHasCrossSMA(result.Bands[bandLen-2:], false) {
			if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.Position == models.BELOW_SMA {
				if isHasCrossLower(midInterval.Bands[bandLen-2:], false) && midHFirstPercentFromSMA < 3.2 {
					ignoredReason = "trend up but starting to down"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend != models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP && result.Position == models.BELOW_SMA {
		if (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
			if shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "short trend side way and percent from upper below 3 2nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.Position == models.BELOW_SMA && hFirstPercentFromSMA < 3.3 {
				ignoredReason = "just minor up still trend down 3nd"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_DOWN && !isHasCrossLower(midInterval.Bands[bandLen-2:], false) {
			if shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.PriceChanges > 2 && countBelowLower(shortInterval.Bands[bandLen-4:], false) >= 1 {
				ignoredReason = "just minor up still trend down 4nd"
				return true
			}
		}

		if result.PriceChanges > 4 && midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrendPercent < 10 {
			if shortInterval.AllTrend.SecondTrend == models.TREND_UP && shortHFourthPercentFromUpper < 3.2 {
				ignoredReason = "just minor up still trend down 5nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP {
		if midInterval.AllTrend.FirstTrend == models.TREND_UP && midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if midSecondLastBand.Candle.Open > midSecondLastBand.Candle.Close && getHighestHightIndex(midInterval.Bands) >= bandLen-3 {
				ignoredReason = "on upper and will starting down"
				return true
			}
		}
	}

	if hFirstPercentFromUpper < 3.4 || midHFirstPercentFromUpper < 3.3 {
		if isHasCrossUpper(shortInterval.Bands[bandLen-2:], true) {
			ignoredReason = "percent from upper below 3 and short cross upper"
			return true
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && countCrossSMA(result.Bands[bandLen-3:]) >= 2 {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend != models.TREND_UP && isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
			if countCrossUpper(midInterval.Bands[bandLen-6:]) > 1 && countCrossUpper(shortInterval.Bands[bandLen/2:]) == 1 {
				ignoredReason = "minor up thren continue down"
				return true
			}
		}

		if midInterval.AllTrend.ShortTrend != models.TREND_UP && isHasCrossUpper(midInterval.Bands[bandLen-4:], true) && midInterval.Position == models.BELOW_SMA {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 && midInterval.Direction == BAND_DOWN {
				ignoredReason = "minor up thren continue down 2nd"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN && result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands[bandLen-3:], false) {
		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.Position == models.BELOW_SMA && !isHasCrossLower(midInterval.Bands[bandLen-4:], false) {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "losng mid short interval below sma but not cross lower"
				return true
			}
		}

		if midInterval.AllTrend.ShortTrend != models.TREND_UP && midInterval.Position == models.BELOW_SMA && isHasCrossLower(midInterval.Bands[bandLen-4:], false) {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 && midHFirstPercentFromSMA < 3.2 {
				ignoredReason = "losng mid short interval below sma but not cross lower 2nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && result.Position == models.ABOVE_SMA {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP && midInterval.Position == models.ABOVE_SMA {
			if shortInterval.AllTrend.SecondTrend == models.TREND_UP && shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.Position == models.ABOVE_SMA {
				if shortHSecondPercentFromUpper < 3.2 && countHigestHightMoreThanClose(shortInterval.Bands[bandLen/2:]) > 2 && countHigestHightMoreThanClose(midInterval.Bands[bandLen/2:]) > 2 {
					ignoredReason = "trend up up but there is not significan up"
					return true
				}

				if shortHSecondPercentFromUpper < 3.2 && midPercentFromUpper < 3.3 && percentFromUpper < 3.4 {
					ignoredReason = "trend up up but all interval percent from upper below 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP && result.Direction == BAND_DOWN {
		if midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY && midLastBand.Candle.Close < float32(midLastBand.Upper) && isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
			if isOnDown(shortInterval) && shortHSecondPercentFromSMA < 3.1 && shortInterval.Position == models.BELOW_SMA {
				ignoredReason = "just down from upper, better not in"
				return true
			}
		}

		if midInterval.AllTrend.ShortTrend == models.TREND_DOWN && midInterval.Position == models.BELOW_SMA && midHFirstPercentFromUpper < 3.3 {
			ignoredReason = "band down, mid short trend down but percent from upper below 3"
			return true
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP {
		if secondLastBand.Candle.Hight-secondLastBand.Candle.Close > secondLastBand.Candle.Close-secondLastBand.Candle.Open {
			if checkDiffHightClose(midInterval.Bands[bandLen/2:]) {
				if shortHSecondPercentFromUpper < 3.2 {
					ignoredReason = "on down and percent below 3"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_UP {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if isHasCrossUpper(shortInterval.Bands[bandLen-2:], true) && (percentFromUpper < 3.4 || midPercentFromUpper < 3.3) {
				ignoredReason = "up but not confident"
				return true
			}

			if isHasCrossUpper(midInterval.Bands[bandLen-2:], true) && !isHasCrossUpper(shortInterval.Bands[bandLen-5:], true) {
				ignoredReason = "up but not confident 2nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && isHasCrossSMA(result.Bands[bandLen-2:], false) {
		if midLastBand.Candle.Close > float32(midLastBand.SMA) && isHasCrossSMA(midInterval.Bands[bandLen-2:], false) && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) {
				ignoredReason = "warning on down, on cross lower"
				return true
			}
		}
	}

	if isHasCrossUpper(shortInterval.Bands[bandLen-2:], true) {
		if !isHasCrossUpper(midInterval.Bands[bandLen-2:], true) {
			if midPercentFromUpper < 3.3 || percentFromUpper < 3.4 {
				ignoredReason = "cross upper and percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossUpper(result.Bands[bandLen-3:], true) {
		if midInterval.AllTrend.SecondTrend != models.TREND_UP || midInterval.AllTrend.ShortTrend != models.TREND_UP {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.AllTrend.ShortTrend == models.TREND_UP {
				if !isHasCrossUpper(shortInterval.Bands[bandLen/2:], true) && isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) {
					ignoredReason = "short trend downm, short intertval cross sma"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend != models.TREND_UP {
		if midInterval.AllTrend.ShortTrend != models.TREND_UP && midInterval.Direction == BAND_DOWN {
			if isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) || shortInterval.AllTrend.ShortTrend != models.TREND_UP {
				if !isHasCrossLower(shortInterval.Bands[bandLen-1:], false) || !isHasCrossLower(midInterval.Bands[bandLen-1:], false) || !(result.Position == models.ABOVE_SMA && result.AllTrend.SecondTrend == models.TREND_UP) {
					ignoredReason = "starting down, mid band down and short cross sma"
					return true
				}
			}
		}
	}

	if isHasCrossUpper(result.Bands[bandLen-2:], true) && result.Direction == BAND_DOWN {
		if midInterval.AllTrend.SecondTrend != models.TREND_UP && midInterval.AllTrend.ShortTrend != models.TREND_UP && !isHasCrossLower(midInterval.Bands[bandLen-2:], false) {
			if shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "band down"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) || shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "band down 2nd"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend != models.TREND_UP && midInterval.Position == models.ABOVE_SMA && !isHasCrossSMA(midInterval.Bands[bandLen-2:], false) {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "starting down gan."
				return true
			}
		}
	}

	if result.AllTrend.Trend != models.TREND_UP {
		if countBelowLower(midInterval.Bands[bandLen-2:], false) > 0 || isHasCrossLower(midInterval.Bands[bandLen-2:], true) {
			if midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2 {
				ignoredReason = "have below lower but percent below 3"
				return true
			}
		}

		if isHasCrossSMA(shortInterval.Bands[bandLen-2:], true) {
			if shortInterval.AllTrend.ShortTrend == models.TREND_UP && shortInterval.PriceChanges > 3.5 {
				ignoredReason = "cross sma and price change more than 3.5"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_UP {
		if countBelowLower(result.Bands[bandLen-3:], false) > 0 && result.PriceChanges > 4 {
			if shortInterval.AllTrend.Trend == models.TREND_UP && isHasCrossUpper(shortInterval.Bands[bandLen/2:], true) && shortInterval.AllTrend.ShortTrend != models.TREND_UP {
				ignoredReason = "minor up but short interval alredy corss upper"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.Trend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if shortInterval.AllTrend.ShortTrend == models.TREND_DOWN && isHasCrossUpper(shortInterval.Bands[bandLen/2:], true) {
				ignoredReason = "trend down, minor up but short interval already cross upper"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.Trend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend != models.TREND_UP {
			if midInterval.Position == models.BELOW_SMA && !isHasCrossLower(midInterval.Bands[bandLen-2:], false) {
				if shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend down, contain below lower, minor up but short interval pecent from sma below 3"
					return true
				}
			}
		}

		if midInterval.AllTrend.Trend == models.TREND_DOWN && countBelowLower(midInterval.Bands[bandLen-6:], false) > 0 {
			if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "trend down, contain below lower, minor up but short interval pecent from sma below 3 2nd"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && (midInterval.AllTrend.ShortTrend != models.TREND_UP || midInterval.Direction == BAND_DOWN) {
			if shortInterval.AllTrend.Trend == models.TREND_DOWN && shortInterval.Position == models.BELOW_SMA && !isHasCrossLower(shortInterval.Bands[bandLen/2:], false) {
				if shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend down, short not cross lower and percent below 3"
					return true
				}
			}
		}

		if isHasCrossSMA(result.Bands[bandLen-4:], false) && !isHasCrossLower(result.Bands[bandLen-4:], false) {
			if midInterval.AllTrend.SecondTrend == models.TREND_DOWN {
				if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend down, not cross lower and percent below 3"
					return true
				}
			}
		}

		if isHasCrossLower(result.Bands[bandLen-4:], false) {
			if isHasCrossLower(midInterval.Bands[bandLen-2:], true) {
				if countBelowLower(shortInterval.Bands[bandLen-4:], false) > 0 || shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend down, start up but inconviencing"
					return true
				}
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_SIDEWAY && result.AllTrend.ShortTrend == models.TREND_UP && result.PriceChanges > 5 {
		if midInterval.AllTrend.Trend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY && isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
			if shortInterval.AllTrend.Trend == models.TREND_UP && shortInterval.Position == models.ABOVE_SMA && shortHFourthPercentFromUpper < 3.2 {
				ignoredReason = "trend up but start to down"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_UP && result.Position == models.ABOVE_SMA && result.Direction == BAND_DOWN {
		if midInterval.AllTrend.Trend == models.TREND_SIDEWAY && midInterval.AllTrend.ShortTrend == models.TREND_UP && midHFirstPercentFromUpper < 3.3 {
			if shortHSecondPercentFromUpper < 3.2 && (isHasCrossUpper(shortInterval.Bands[bandLen-6:], true) || isHasCrossUpper(midInterval.Bands[bandLen-4:], true)) {
				ignoredReason = "trend up but start to down 2nd"
				return true
			}
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_SIDEWAY && isHasCrossUpper(result.Bands[bandLen-4:], true) {
		if countAboveUpper(result.Bands[bandLen-4:]) > 0 || result.Direction == BAND_DOWN {
			if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) && countBelowSMA(midInterval.Bands[bandLen-4:], true) > 0 {
				if shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend up but signal down"
					return true
				}
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if isHasCrossUpper(result.Bands[bandLen/2:], true) || !isHasCrossLower(result.Bands[bandLen-2:], false) {
				ignoredReason = "sideway2"
				return true
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && isHasCrossUpper(midInterval.Bands[bandLen/2:], true) && midInterval.Position == models.ABOVE_SMA {
			if shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend up already cross upper and starting down"
				return true
			}
		}

		if midInterval.AllTrend.ShortTrend == models.TREND_UP && countDownBand(midInterval.Bands[bandLen-4:]) > 2 {
			if shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend up already cross upper and starting down 2nd"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if isHasCrossLower(midInterval.Bands[bandLen-4:], false) && isHasCrossUpper(midInterval.Bands[bandLen/2:], true) && midInterval.AllTrend.Trend == models.TREND_DOWN {
			if shortInterval.AllTrend.Trend == models.TREND_DOWN && !isHasCrossLower(shortInterval.Bands[bandLen-4:], false) && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "trend down and start minor up but percent below 3"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.Trend == models.TREND_DOWN && isHasCrossLower(midInterval.Bands[bandLen-3:], true) {
			if shortInterval.AllTrend.Trend == models.TREND_DOWN && !isHasCrossLower(shortInterval.Bands[bandLen-4:], false) && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "trend down and start minor up but percent 2nd"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_DOWN && (countBelowLower(result.Bands[bandLen-4:], false) > 0 || isHasCrossLower(result.Bands[bandLen-4:], true)) {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if shortInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) && isHasCrossUpper(shortInterval.Bands[bandLen/2:], true) {
				ignoredReason = "on lower, minor up but short interval already cross upper"
				return true
			}
		}
	}

	if result.Position == models.ABOVE_SMA && !isHasCrossUpper(result.Bands[bandLen-4:], true) && isHasCrossUpper(result.Bands[bandLen/2:], true) {
		if midInterval.Position == models.BELOW_SMA && !isHasCrossLower(midInterval.Bands[bandLen-4:], false) && midHFirstPercentFromSMA < 3.2 {
			if shortInterval.AllTrend.Trend == models.TREND_SIDEWAY && shortInterval.AllTrend.ShortTrend != models.TREND_UP && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "down from uppper and short interval sideway"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && result.Position == models.BELOW_SMA && isHasCrossUpper(result.Bands[bandLen/2:], true) && !isHasCrossLower(result.Bands[bandLen-2:], false) {
		if midInterval.AllTrend.FirstTrend == models.TREND_DOWN && midInterval.AllTrend.SecondTrend != models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_UP {
			if midHFirstPercentFromSMA < 3.2 && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "down from uppper and short up percent from upper below 3"
				return true
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN && countBelowLower(result.Bands[bandLen-4:], false) > 0 {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "down from uppper and short up percent from upper below 3 2nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend != models.TREND_UP && !isHasCrossLower(result.Bands[bandLen-2:], false) {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "down from uppper and short up percent from upper below 3 3nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && isHasCrossUpper(result.Bands[bandLen-2:], true) && isBandHeadDoubleBody(result.Bands[bandLen-2:]) {
		if isBandHeadDoubleBody(midInterval.Bands[bandLen-2:]) {
			if shortInterval.AllTrend.Trend != models.TREND_UP && shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 {
				ignoredReason = "down, head double body and percent below 3"
				return true
			}
		}
	}

	if isBandHeadDoubleBody(result.Bands[bandLen-2:]) || ((result.AllTrend.Trend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN) && result.AllTrend.ShortTrend == models.TREND_DOWN) {
		if midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY && isHasCrossUpper(midInterval.Bands[bandLen-4:], true) {
			if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortInterval.Position == models.ABOVE_SMA {
				ignoredReason = "already cross upper and start down"
				return true
			}
		}

		if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) {
			if isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) || (shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 && midHFirstPercentFromUpper < 3.3) {
				ignoredReason = "down, and cross sma"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && isHasCrossUpper(result.Bands[bandLen-2:], true) {
		if midInterval.AllTrend.Trend == models.TREND_UP && midInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(midInterval.Bands[bandLen-5:], true) {
			if shortInterval.AllTrend.Trend == models.TREND_DOWN && shortInterval.Position == models.BELOW_SMA && !isHasCrossLower(shortInterval.Bands[bandLen/2:], false) {
				if countAboveSMA(shortInterval.Bands[bandLen-5:]) == 0 {
					ignoredReason = "starting down after up"
					return true
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN || (result.AllTrend.FirstTrend == models.TREND_DOWN && result.AllTrend.SecondTrend == models.TREND_SIDEWAY) {
		if midInterval.AllTrend.Trend == models.TREND_UP && midInterval.Position == models.ABOVE_SMA && !isHasCrossUpper(midInterval.Bands[bandLen-5:], true) {
			if !isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "starting down after up 2nd"
				return true
			}
		}
	}

	if (result.AllTrend.Trend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN) && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.ShortTrend == models.TREND_UP && isBandHeadDoubleBody(midInterval.Bands[bandLen-2:]) {
			if (midInterval.Position == models.BELOW_SMA && midHFirstPercentFromSMA < 3.2) || (midInterval.Position == models.ABOVE_SMA && midHFirstPercentFromUpper < 3.3) {
				if shortInterval.AllTrend.ShortTrend != models.TREND_UP {
					ignoredReason = "another minor up and start down"
					return true
				}
			}
		}
	}

	if result.AllTrend.Trend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN {
		if result.Direction == BAND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
			if (midInterval.AllTrend.SecondTrend == models.TREND_DOWN || midInterval.AllTrend.Trend == models.TREND_DOWN) && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
				if shortInterval.AllTrend.Trend == models.TREND_DOWN || shortInterval.AllTrend.SecondTrend == models.TREND_DOWN {
					ignoredReason = "all trend down and long interval band is down"
					return true
				}
			}
		}

		if midInterval.AllTrend.ShortTrend != models.TREND_UP && midInterval.Position == models.BELOW_SMA {
			if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) && !isHasCrossLower(midInterval.Bands[bandLen-6:], false) {
				if countBelowSMA(midInterval.Bands[bandLen-7:bandLen-1], true) == 0 {
					if isHasCrossSMA(shortInterval.Bands[bandLen-2:], false) || midPercentFromUpper < 3.3 || !isHasCrossLower(shortInterval.Bands[bandLen-4:], false) {
						ignoredReason = "tren down, mionor up and already cross sma"
						return true
					}
				}
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.Direction == BAND_DOWN {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.Direction == BAND_DOWN {
			ignoredReason = "tren up but band down"
			return true
		}
	}

	if result.AllTrend.ShortTrend == models.TREND_UP && (result.Direction == BAND_DOWN || percentFromUpper < 3.4 || isHasCrossUpper(result.Bands[bandLen-2:], true)) {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP {
			if shortHSecondPercentFromUpper < 3.2 && midHFirstPercentFromUpper < 3.3 {
				ignoredReason = "all interval trend up and pecent from upper below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if result.Position == models.BELOW_SMA && !isHasCrossLower(result.Bands[bandLen-3:], false) {
			if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && (midInterval.AllTrend.ShortTrend == models.TREND_DOWN || (midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY && isHasCrossLower(midInterval.Bands[bandLen-2:], false))) {
				if shortInterval.AllTrend.SecondTrend == models.TREND_DOWN && shortHSecondPercentFromSMA < 3.1 {
					ignoredReason = "trend down below sma but not cross lower"
					return true
				}
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.AllTrend.ShortTrend == models.TREND_DOWN {
			if midInterval.Position == models.BELOW_SMA && !isHasCrossLower(midInterval.Bands[bandLen-3:], false) {
				if shortInterval.Position == models.BELOW_SMA && shortHSecondPercentFromSMA < 3.1 && shortInterval.PriceChanges > 1 {
					ignoredReason = "on down and just minor up"
					return true
				}
			}
		}

		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && countAboveSMA(midInterval.Bands[bandLen/2:]) == 0 {
			if midHFirstPercentFromSMA < 3.2 && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "on down and just minor up 2nd"
				return true
			}
		}
	}

	if isHasCrossUpper(result.Bands[bandLen-2:], true) {
		if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) {
			if !isHasCrossUpper(shortInterval.Bands[bandLen-2:], true) && shortHSecondPercentFromUpper < 3.2 {
				ignoredReason = "trend up already cross upper and percent below 3"
				return true
			}
		}
	}

	if isHasCrossSMA(result.Bands[bandLen-2:], false) {
		if isHasCrossSMA(midInterval.Bands[bandLen-2:], false) {
			if isHasCrossUpper(shortInterval.Bands[bandLen-4:], true) {
				if result.AllTrend.ShortTrend != models.TREND_UP || midInterval.AllTrend.ShortTrend != models.TREND_UP || shortInterval.AllTrend.ShortTrend != models.TREND_UP {
					ignoredReason = "trend up and cross upper"
					return true
				}
			}
		}
	}

	if result.Direction == BAND_DOWN || result.AllTrend.ShortTrend != models.TREND_UP {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY && isHasCrossUpper(midInterval.Bands[bandLen-3:], true) {
			if shortInterval.AllTrend.ShortTrend == models.TREND_DOWN {
				ignoredReason = "possibility down, skip 2nd"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.ShortTrend == models.TREND_SIDEWAY {
		if midInterval.AllTrend.SecondTrend == models.TREND_UP && midInterval.AllTrend.ShortTrend == models.TREND_SIDEWAY {
			if shortHFourthPercentFromUpper < 3.2 {
				ignoredReason = "trend sideway and percent from upper below 3"
				return true
			}
		}
	}

	if result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
		if midInterval.AllTrend.SecondTrend == models.TREND_DOWN && midInterval.Position == models.BELOW_SMA {
			if shortHSecondPercentFromUpper < 3.2 && midHFirstPercentFromSMA < 3.2 {
				ignoredReason = "trend down and minor up aja"
				return true
			}
		}
	}

	return false
}

func isOnDown(result *models.BandResult) bool {
	if result.AllTrend.FirstTrend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN {
		return checkDiffHightClose(result.Bands)
	}

	return false
}

func checkDiffHightClose(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	higestPrice := getHigestPrice(bands)
	if higestPrice != lastBand.Candle.Close {
		diff := (higestPrice - lastBand.Candle.Close) / lastBand.Candle.Close * 100
		return diff > 3
	}
	return false
}

func countHigestHightMoreThanClose(bands []models.Band) int {
	lastBand := bands[len(bands)-1]
	count := 0
	for _, band := range bands[:len(bands)-1] {
		if band.Candle.Hight > lastBand.Candle.Close {
			count++
		}
	}

	return count
}

func afterUpThenDown(bands []models.Band) bool {
	higest := getHighestIndex(bands)
	lowest := getLowestIndex(bands)
	lastBand := bands[len(bands)-1]

	percentFromMidUpper := (float32(((lastBand.Upper-lastBand.SMA)/2)+lastBand.SMA) - lastBand.Candle.Close) / lastBand.Candle.Close * 100
	if higest < len(bands)/2 && lowest < higest {
		trend := CalculateTrendsDetail(bands[higest:])
		if bands[higest].Candle.Hight > float32(bands[higest].Upper) && trend.FirstTrend == models.TREND_DOWN && trend.SecondTrend != models.TREND_DOWN {
			return percentFromMidUpper < 3
		}
	}

	return false
}

func reversal2nd(bands models.BandResult) bool {
	midLowest := getLowestIndex(bands.Bands)
	if midLowest >= len(bands.Bands)-6 && midLowest < len(bands.Bands)-2 && CalculateTrendShort(bands.Bands[midLowest:]) == models.TREND_UP {
		if bands.Bands[midLowest].Candle.Close < float32(bands.Bands[midLowest].Lower) || bands.Bands[midLowest].Candle.Open > float32(bands.Bands[midLowest].Lower) {
			if !(bands.Bands[midLowest].Candle.Close < float32(bands.Bands[midLowest].Lower) && bands.Bands[midLowest].Candle.Open > float32(bands.Bands[midLowest].Lower)) {
				return true
			}
		}
	}

	return false
}
