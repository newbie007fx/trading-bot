package analysis

import (
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""
var matchPattern string = ""

func GetIgnoredReason() string {
	return ignoredReason
}

func GetMatchPattern() string {
	return matchPattern
}

func lowestFromBand(band models.Band) float32 {
	if band.Candle.Open > band.Candle.Close {
		return band.Candle.Close
	}

	return band.Candle.Open
}

func CountUpBand(bands []models.Band) int {
	counter := 0
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			difference := band.Candle.Close - band.Candle.Open
			if (difference / band.Candle.Open * 100) > 0.1 {
				counter++
			}
		}
	}

	return counter
}

func isLastBandDoublePreviousHeigest(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	bandList := bands[len(bands)-5 : len(bands)-1]

	return isBandMultipleThanList(lastBand, bandList, 2)
}

func isLastBandMultiplePreviousHeigest(bands []models.Band, multiplier int) bool {
	lastBand := bands[len(bands)-1]
	bandList := bands[len(bands)-5 : len(bands)-1]

	return isBandMultipleThanList(lastBand, bandList, multiplier)
}

func isBandMultipleThanList(band models.Band, bands []models.Band, multiplier int) bool {
	bandBodyHeight := band.Candle.Close - band.Candle.Open

	var higestBody float32 = 0
	for _, banL := range bands {
		bodyHeight := banL.Candle.Close - banL.Candle.Open
		if banL.Candle.Close < banL.Candle.Open {
			bodyHeight = banL.Candle.Open - banL.Candle.Close
		}

		if bodyHeight > higestBody {
			higestBody = bodyHeight
		}
	}

	return higestBody*float32(multiplier) < bandBodyHeight
}

func bandPercent(band models.Band) float32 {
	return (band.Candle.Close - band.Candle.Open) / band.Candle.Open * 100
}

func bandPercentFromUpper(band models.Band) float32 {
	return (float32(band.Upper) - band.Candle.Close) / band.Candle.Close * 100
}

func isLastBandHeigestHeightBand(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	lastBandClose := lastBand.Candle.Close

	for _, band := range bands[:len(bands)-1] {
		if band.Candle.Hight > lastBandClose {
			return false
		}
	}

	return true
}

func isOpenCloseBelowSMA(band models.Band) bool {
	return band.Candle.Open < float32(band.SMA) && band.Candle.Close < float32(band.SMA)
}

func countOpenCloseBelowSMA(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isOpenCloseBelowSMA(band) {
			count++
		}
	}
	return count
}

func isHasBelowSMA(band models.Band) bool {
	return band.Candle.Close < float32(band.SMA) || band.Candle.Open < float32(band.SMA)
}

func countHasBelowSMA(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isHasBelowSMA(band) {
			count++
		}
	}
	return count
}

func isOpenOrCloseBelowLower(band models.Band) bool {
	return band.Candle.Open < float32(band.Lower) || band.Candle.Close < float32(band.Lower)
}

func countOpenOrCloseBelowLower(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isOpenOrCloseBelowLower(band) {
			count++
		}
	}
	return count
}

func isHasOpenOrCloseBelowLower(bands []models.Band) bool {
	for _, band := range bands {
		if isOpenOrCloseBelowLower(band) {
			return true
		}
	}
	return false
}

func isOpenCloseAboveUpper(band models.Band) bool {
	return band.Candle.Open > float32(band.Upper) && band.Candle.Close > float32(band.Upper)
}

func countOpenCloseAboveUpper(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isOpenCloseAboveUpper(band) {
			count++
		}
	}
	return count
}

func isTailMoreThanBody(band models.Band, percentCheck bool) bool {
	tail := band.Candle.Open - band.Candle.Low
	body := band.Candle.Close - band.Candle.Open

	if tail > body {
		if percentCheck {
			return tail/band.Candle.Low*100 >= 3
		} else {
			return true
		}
	}

	return false
}

func isHeadMoreThanBody(band models.Band, percentCheck bool) bool {
	head := band.Candle.Hight - band.Candle.Close
	body := band.Candle.Close - band.Candle.Open

	if head > body {
		if percentCheck {
			return head/band.Candle.Close*100 >= 3
		} else {
			return true
		}
	}

	return false
}

func isUpperHeadMoreThanUpperBody(band models.Band) bool {
	if band.Candle.Open > band.Candle.Close {
		return false
	}

	if band.Candle.Open > float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
		return true
	}

	if !(band.Candle.Open <= float32(band.Upper) && band.Candle.Close >= float32(band.Upper)) {
		return false
	}

	upperToHead := band.Candle.Close - float32(band.Upper)
	upperToBody := float32(band.Upper) - band.Candle.Open

	return upperToHead > upperToBody
}

func countUpperHeadMoreThanUpperBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isUpperHeadMoreThanUpperBody(band) {
			count++
		}
	}
	return count
}

func isAboveUpperAndOrUpperHeadMoreThanUpperBody(band models.Band, bands []models.Band) bool {
	if isOpenCloseAboveUpper(band) && countUpperHeadMoreThanUpperBody(bands) > 1 {
		return true
	}

	if isOpenCloseAboveUpper(band) && countOpenCloseAboveUpper(bands) > 1 {
		return true
	}

	if isOpenCloseAboveUpper(band) && countOpenCloseAboveUpper(bands) == 1 && countUpperHeadMoreThanUpperBody(bands) == 1 {
		return true
	}

	return false
}

func isCrossLowerOnBody(band models.Band) bool {
	return band.Candle.Close < float32(band.Lower) || band.Candle.Open < float32(band.Lower)
}

func countCrossLowerOnBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isCrossLowerOnBody(band) {
			count++
		}
	}

	return count
}

func isCrossLower(band models.Band) bool {
	return band.Candle.Low < float32(band.Lower)
}

func countCrossLower(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isCrossLower(band) {
			count++
		}
	}

	return count
}

func isCrossUpSMAOnBody(band models.Band) bool {
	return band.Candle.Open < float32(band.SMA) && band.Candle.Close > float32(band.SMA)
}

func countCrossUpSMAOnBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isCrossUpSMAOnBody(band) {
			count++
		}
	}

	return count
}

func isCrossUpUpperOnBody(band models.Band) bool {
	return band.Candle.Close > float32(band.Upper)
}

func countCrossUpUpperOnBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isCrossUpUpperOnBody(band) {
			count++
		}
	}

	return count
}

func isHightCrossUpper(band models.Band) bool {
	return band.Candle.Hight >= float32(band.Upper)
}

func countHightCrossUpper(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isHightCrossUpper(band) {
			count++
		}
	}

	return count
}

func isBadBand(band models.Band, percentCheck bool) bool {
	return isHeadMoreThanBody(band, percentCheck) || isBandDown(band) || isTailMoreThanBody(band, percentCheck) || band.Candle.Open == band.Candle.Close
}

func isBandDown(band models.Band) bool {
	return band.Candle.Open > band.Candle.Close
}

func countBadBandAndCrossUpper(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if (isBadBand(band, false) && isHightCrossUpper(band)) || isOpenCloseAboveUpper(band) {
			count++
		}
	}

	return count
}

func CountBadBand(bands []models.Band, percentCheck bool) int {
	count := 0
	for _, band := range bands {
		if isBadBand(band, percentCheck) {
			count++
		}
	}

	return count
}

func countOpenCloseAboveUpperOrUpperHeadMoreThanUpperBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isOpenCloseAboveUpper(band) || isUpperHeadMoreThanUpperBody(band) {
			count++
		}
	}

	return count
}

func isCloseAboveSMA(band models.Band) bool {
	return band.Candle.Close > float32(band.SMA) && band.Candle.Close < float32(band.Upper)
}

func countCloseAboveSMA(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isCloseAboveSMA(band) {
			count++
		}
	}

	return count
}

func isOpenCloseAboveSMA(band models.Band) bool {
	return band.Candle.Open > float32(band.SMA) && band.Candle.Close > float32(band.SMA)
}

func countOpenCloseAboveSMA(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isOpenCloseAboveSMA(band) {
			count++
		}
	}

	return count
}

func isDoubleSignificanUp(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	listBands := bands[len(bands)-6 : len(bands)-2]
	if CountBadBand(bands[len(bands)-2:], false) == 0 {
		if isBandMultipleThanList(lastBand, listBands, 2) && isBandMultipleThanList(secondLastBand, listBands, 2) {
			return true
		}
	}

	return false
}

func isDownFromUpper(band models.Band) bool {
	return band.Candle.Open > float32(band.Upper) && band.Candle.Open > band.Candle.Close
}

func countDownFromUpper(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isDownFromUpper(band) {
			count++
		}
	}

	return count
}

func isSolidBand(band models.Band) bool {
	if band.Candle.Close < band.Candle.Open {
		return false
	}

	body := band.Candle.Close - band.Candle.Open
	head := band.Candle.Hight - band.Candle.Close
	tail := band.Candle.Open - band.Candle.Low

	headPercent := head / body * 100
	tailPercent := tail / body * 100
	headThreshold := 10
	if tailPercent == 0 {
		headThreshold = 15
	}

	return tailPercent < 10 && headPercent < float32(headThreshold)
}

func isContainNotTrendup(result models.BandResult) bool {
	return result.AllTrend.FirstTrend != models.TREND_UP || result.AllTrend.SecondTrend != models.TREND_UP
}

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time, isNoNeedDoubleCheck bool) bool {
	ignoredReason = ""
	matchPattern = ""

	if isOnBandCompleteCheck(currentTime) {
		return approvedPatternOnCompleteCheck(short, mid, long, currentTime, isNoNeedDoubleCheck)
	}

	return false
}

func approvedPatternOnCompleteCheck(short, mid, long models.BandResult, currentTime time.Time, isNoNeedDoubleCheck bool) bool {
	bandLen := len(short.Bands)
	longLastBand := long.Bands[bandLen-1]
	midLastBand := mid.Bands[bandLen-1]
	shortLastBand := short.Bands[bandLen-1]

	if (isSolidBand(shortLastBand) || isLastBandDoublePreviousHeigest(short.Bands) || isNoNeedDoubleCheck) && bandPercent(shortLastBand) >= 2 {
		if CountBadBand(short.Bands[bandLen-4:bandLen-1], false) < 3 {
			if isUpperHeadMoreThanUpperBody(midLastBand) && (isUpperHeadMoreThanUpperBody(longLastBand) || isBandMultipleThanList(longLastBand, long.Bands[bandLen-4:bandLen-1], 2)) {
				if short.Position == models.ABOVE_UPPER && (isLastBandDoublePreviousHeigest(short.Bands) || countOpenCloseAboveUpper(short.Bands[bandLen-4:]) > 0) {
					ignoredReason = "mid and long upper head more than body"
					return false
				}
			}

			if mid.Position == models.ABOVE_SMA && mid.AllTrend.Trend == models.TREND_DOWN && (countCrossUpSMAOnBody(mid.Bands[bandLen-10:]) == 1 || countCrossUpUpperOnBody(mid.Bands[bandLen-10:]) == 1) {
				if countOpenCloseBelowSMA(mid.Bands[bandLen-10:]) >= 7 && long.AllTrend.SecondTrend == models.TREND_DOWN && countOpenCloseBelowSMA(long.Bands[bandLen-4:]) < 2 {
					ignoredReason = "mid first up on down"
					return false
				}
			}

			if isBadBand(shortLastBand, false) && isBadBand(midLastBand, false) && isBadBand(longLastBand, false) {
				ignoredReason = "all bad band"
				return false
			}

			if long.Position == models.ABOVE_UPPER && countUpperHeadMoreThanUpperBody(long.Bands[bandLen-4:]) > 1 {
				if short.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(shortLastBand) {
					ignoredReason = "max upper"
					return false
				}
			}

			if shortLastBand.Candle.Close < float32(shortLastBand.Upper) {
				if CalculateShortTrendWithConclusion(short.Bands[:bandLen-1]) == models.TREND_DOWN {
					if countCrossLowerOnBody(short.Bands[bandLen-3:bandLen-1]) > 0 {
						ignoredReason = "short cross lower"
						return false
					}
				}

				if short.Position == models.ABOVE_SMA && (countOpenCloseBelowSMA(short.Bands[bandLen-10:]) >= 9 && countCrossLower(mid.Bands[bandLen-4:]) > 0) {
					if mid.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_DOWN {
						ignoredReason = "short first cross sma, mid last 4 cross lower"
						return false
					}
				}

				if long.AllTrend.ShortTrend == models.TREND_DOWN && long.AllTrend.SecondTrend == models.TREND_DOWN && mid.AllTrend.SecondTrend == models.TREND_DOWN {
					if short.AllTrend.FirstTrend == models.TREND_DOWN && countHightCrossUpper(short.Bands[bandLen-10:]) == 0 && bandPercentFromUpper(shortLastBand) < 1 {
						ignoredReason = "first up on down trend"
						return false
					}
				}

				if isUpperHeadMoreThanUpperBody(longLastBand) && CountBadBand(long.Bands[bandLen-2:], false) > 0 {
					if isLastBandDoublePreviousHeigest(mid.Bands) && CountBadBand(mid.Bands[bandLen-4:], false) > 2 {
						if isLastBandDoublePreviousHeigest(short.Bands) {
							ignoredReason = "up down on long upper head more than body"
							return false
						}
					}
				}

				if short.AllTrend.SecondTrend != models.TREND_UP && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend != models.TREND_UP {
					matchPattern = "below upper but down from upper"
					return true
				}

				if long.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER {
					if countUpperHeadMoreThanUpperBody(long.Bands[bandLen-2:]) > 0 && countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-2:]) > 0 {
						matchPattern = "below upper but significan up trend"
						return true
					}
				}

				matchPattern = "below upper"
				return true
			}

			if shortLastBand.Candle.Close > float32(shortLastBand.Upper) {
				if long.AllTrend.ShortTrend == models.TREND_DOWN {
					if countCrossLowerOnBody(long.Bands[bandLen-2:]) > 0 {
						ignoredReason = "short above upper and long cross lower"
						return false
					}
				}

				if midLastBand.Candle.Close > float32(midLastBand.SMA) {
					if countCrossLowerOnBody(mid.Bands[bandLen-5:]) > 0 {
						if short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
							if short.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(short.Bands[bandLen-10:]) == 1 {
								ignoredReason = "short above upper and first up after down"
								return false
							}
						}
					}
				}

				if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && isBandMultipleThanList(shortLastBand, short.Bands[bandLen-5:bandLen-1], 3) {
					if mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER {
						if countHightCrossUpper(long.Bands[bandLen-4:]) > 1 && isBandDown(long.Bands[bandLen-2]) {
							ignoredReason = "short above upper and all above upper but bad"
							return false
						}
					}
				}

				if isFirstCrossSMA(mid.Bands) && long.AllTrend.ShortTrend == models.TREND_DOWN {
					ignoredReason = "short above upper and mid first cross sma, long short trend down"
					return false
				}

				if isFirstCrossSMA(mid.Bands) && isLastBandDoublePreviousHeigest(mid.Bands) && isFirstCrossSMA(long.Bands) && isLastBandDoublePreviousHeigest(long.Bands) {
					ignoredReason = "short above upper and mid and long first cross sma"
					return false
				}

				if long.Position == models.ABOVE_UPPER && countCrossUpSMAOnBody(long.Bands[bandLen-10:]) == 1 {
					if countHasBelowSMA(mid.Bands[bandLen-15:]) > 10 && mid.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(mid.Bands[bandLen-15:]) == 1 {
						ignoredReason = "mid first cross upper after mostly below sma"
						return false
					}
				}

				if isLastBandDoublePreviousHeigest(short.Bands) && countUpperHeadMoreThanUpperBody(short.Bands[bandLen-3:]) > 1 {
					ignoredReason = "mid first cross upper after mostly below sma"
					return false
				}

				if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP && long.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 1 {
					if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP && isLastBandDoublePreviousHeigest(mid.Bands) && mid.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
						if isLastBandDoublePreviousHeigest(short.Bands) {
							ignoredReason = "double and first cross upper"
							return false
						}
					}
				}

				if long.Position == models.ABOVE_SMA && countOpenCloseAboveSMA(long.Bands[bandLen-10:]) <= 1 && long.AllTrend.Trend == models.TREND_DOWN {
					if isUpperHeadMoreThanUpperBody(shortLastBand) && isUpperHeadMoreThanUpperBody(midLastBand) {
						ignoredReason = "on down, short and mid upper head more than body"
						return false
					}
				}

				if short.Position == models.ABOVE_UPPER && short.AllTrend.Trend == models.TREND_UP {
					if mid.Position == models.ABOVE_UPPER && mid.AllTrend.Trend == models.TREND_UP {
						if long.Position == models.ABOVE_UPPER && long.AllTrend.Trend == models.TREND_UP {
							if long.AllTrend.FirstTrend == models.TREND_UP && long.AllTrend.SecondTrend == models.TREND_UP {
								if long.AllTrend.ShortTrend == models.TREND_UP && countHightCrossUpper(long.Bands[bandLen-5:]) > 3 {
									if isLastBandDoublePreviousHeigest(short.Bands) && isLastBandDoublePreviousHeigest(mid.Bands) {
										ignoredReason = "up up and long cross upper more than 2"
										return false
									}
								}
							}
						}
					}
				}

				if midLastBand.Candle.Close > float32(midLastBand.SMA) && countOpenCloseBelowSMA(mid.Bands[bandLen-6:bandLen-1]) == 5 {
					if mid.AllTrend.FirstTrend == models.TREND_DOWN || mid.AllTrend.SecondTrend == models.TREND_DOWN {
						if isUpperHeadMoreThanUpperBody(shortLastBand) || isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) {
							ignoredReason = "first up ate down"
							return false
						}
					}
				}

				if long.AllTrend.FirstTrend == models.TREND_DOWN && longLastBand.Candle.Close > float32(longLastBand.SMA) && countHasBelowSMA(long.Bands[bandLen-10:]) == 9 {
					if mid.Position == models.ABOVE_UPPER && short.Position == models.ABOVE_UPPER {
						if isLastBandDoublePreviousHeigest(mid.Bands) || isLastBandDoublePreviousHeigest(short.Bands) {
							ignoredReason = "first up on long trend down"
							return false
						}
					}
				}

				if long.AllTrend.SecondTrend == models.TREND_DOWN && (long.AllTrend.ShortTrend == models.TREND_DOWN || CalculateTrendShort(long.Bands[bandLen-6:bandLen-1], true) == models.TREND_DOWN) {
					if long.AllTrend.Trend == models.TREND_DOWN && mid.AllTrend.Trend == models.TREND_DOWN {
						if isLastBandDoublePreviousHeigest(mid.Bands) || isLastBandDoublePreviousHeigest(short.Bands) {
							ignoredReason = "on long trend down minor up"
							return false
						}
					}
				}

				if long.AllTrend.Trend == models.TREND_DOWN {
					if countCrossLowerOnBody(short.Bands[bandLen-7:]) > 0 && isLastBandDoublePreviousHeigest(short.Bands) {
						ignoredReason = "on long trend down, short up down"
						return false
					}
				}

				if long.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER {
					if long.AllTrend.SecondTrend == models.TREND_UP && countOpenCloseBelowSMA(long.Bands[bandLen-11:]) > 0 && countCrossUpUpperOnBody(long.Bands[bandLen-5:]) == 1 {
						if isLastBandDoublePreviousHeigest(short.Bands) {
							ignoredReason = "on long trend up klimaks"
							return false
						}
					}

					if isUpperHeadMoreThanUpperBody(longLastBand) && isLastBandDoublePreviousHeigest(short.Bands) {
						ignoredReason = "on long trend up klimaks2"
						return false
					}

					if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP && countCrossUpUpperOnBody(long.Bands[bandLen-4:]) > 1 {
						if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP && isLastBandDoublePreviousHeigest(mid.Bands) {
							if isLastBandDoublePreviousHeigest(short.Bands) && isUpperHeadMoreThanUpperBody(shortLastBand) {
								ignoredReason = "on long trend up klimaks3"
								return false
							}
						}
					}
				}

				if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend != models.TREND_UP {
					if isLastBandDoublePreviousHeigest(short.Bands) {
						ignoredReason = "on long trend up but starting down"
						return false
					}
				}

				if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP && CalculateTrendShort(long.Bands[bandLen-6:bandLen-1], true) == models.TREND_DOWN {
					if mid.Position == models.ABOVE_UPPER && mid.AllTrend.ShortTrend == models.TREND_UP && countOpenCloseAboveSMA(mid.Bands[bandLen-5:bandLen-1]) == 0 && isLastBandDoublePreviousHeigest(mid.Bands) {
						if short.AllTrend.ShortTrend == models.TREND_UP && countOpenCloseAboveSMA(short.Bands[bandLen-5:bandLen-1]) == 0 && isLastBandDoublePreviousHeigest(short.Bands) {
							ignoredReason = "on long trend up, after klimaks, minor up"
							return false
						}
					}
				}

				if long.AllTrend.ShortTrend == models.TREND_UP && long.Position == models.ABOVE_UPPER {
					if isDownAfterSignificanUp(mid) && isLastBandDoublePreviousHeigest(short.Bands) {
						ignoredReason = "after significan up"
						return false
					}
				}

				matchPattern = "short above upper"
				return true
			}
		}
	}

	return false
}

func isDownAfterSignificanUp(result models.BandResult) bool {
	bands := result.Bands
	if result.AllTrend.SecondTrend == models.TREND_UP && result.AllTrend.SecondTrendPercent < 10 {
		midLen := len(bands) / 2
		bandsSecond := bands[len(bands)-midLen:]
		lowest := getLowesIndexPrice(bandsSecond)
		higest := getHigestIndexPrice(bandsSecond)
		if lowest < higest && result.AllTrend.ShortTrend == models.TREND_UP {
			return higest < len(bandsSecond)-5
		}
	}

	return false
}

func getHigestIndexPrice(bands []models.Band) int {
	var index int = 0
	for i, band := range bands {
		if bands[index].Candle.Close < band.Candle.Close {
			index = i
		}
	}

	return index
}

func getLowesIndexPrice(bands []models.Band) int {
	var index int = 0
	for i, band := range bands {
		if bands[index].Candle.Close > band.Candle.Close {
			index = i
		}
	}

	return index
}

func isFirstCrossSMA(bands []models.Band) bool {
	if isCrossUpSMAOnBody(bands[len(bands)-1]) {
		return countOpenCloseBelowSMA(bands[len(bands)-5:len(bands)-1]) == 4
	}

	return false
}

func countTrendDown(short, mid, long models.BandResult) int {
	count := 0

	if short.AllTrend.FirstTrend != models.TREND_UP {
		count++
	}
	if short.AllTrend.SecondTrend != models.TREND_UP {
		count++
	}

	if mid.AllTrend.FirstTrend != models.TREND_UP {
		count++
	}
	if mid.AllTrend.SecondTrend != models.TREND_UP {
		count++
	}

	if long.AllTrend.FirstTrend != models.TREND_UP {
		count++
	}
	if long.AllTrend.SecondTrend != models.TREND_UP {
		count++
	}

	return count
}

func isOnBandCompleteCheck(currentTime time.Time) bool {
	minute := currentTime.Minute()

	if minute == 59 || minute == 0 {
		return true
	} else if minute == 14 || minute == 15 {
		return true
	} else if minute == 29 || minute == 30 {
		return true
	} else if minute == 44 || minute == 45 {
		return true
	}

	return false
}

func isBandLongTimeComplete(currentTime time.Time) bool {
	hour := currentTime.Hour()
	minute := currentTime.Minute()

	if (hour%4 == 0 || hour%4 == 3) && (minute > 57 || minute < 2) {
		return true
	}

	return false
}
