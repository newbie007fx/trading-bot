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

func isTailMoreThanBody(band models.Band) bool {
	tail := band.Candle.Open - band.Candle.Low
	body := band.Candle.Close - band.Candle.Open

	return tail > body
}

func isHeadMoreThanBody(band models.Band) bool {
	head := band.Candle.Hight - band.Candle.Close
	body := band.Candle.Close - band.Candle.Open

	return head > body
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

func isCrossUpSMAOnBody(band models.Band) bool {
	return band.Candle.Close > float32(band.SMA)
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

func isBadBand(band models.Band) bool {
	return isHeadMoreThanBody(band) || isBandDown(band) || isTailMoreThanBody(band) || band.Candle.Open == band.Candle.Close
}

func isBandDown(band models.Band) bool {
	return band.Candle.Open > band.Candle.Close
}

func countBadBandAndCrossUpper(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if (isBadBand(band) && isHightCrossUpper(band)) || isOpenCloseAboveUpper(band) {
			count++
		}
	}

	return count
}

func CountBadBand(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isBadBand(band) {
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

func isDoubleSignificanUp(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	secondLastBand := bands[len(bands)-2]
	listBands := bands[len(bands)-6 : len(bands)-2]
	if CountBadBand(bands[len(bands)-2:]) == 0 {
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

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time) bool {
	ignoredReason = ""
	matchPattern = ""

	if isOnFirstCheck(currentTime) {
		return approvedPatternFirstCheck(short, mid, long)
	}

	if isOnBandCompleteCheck(currentTime) {
		return approvedPatternOnCompleteCheck(short, mid, long, currentTime)
	}

	return false
}

func approvedPatternFirstCheck(short, mid, long models.BandResult) bool {
	bandLen := len(short.Bands)
	shortLastBand := short.Bands[bandLen-1]
	shortSecondLastBand := short.Bands[bandLen-2]
	midLastBand := mid.Bands[bandLen-1]
	longLastBand := long.Bands[bandLen-1]

	if !(isUpperHeadMoreThanUpperBody(longLastBand) || isOpenCloseAboveUpper(longLastBand)) {
		if !(isUpperHeadMoreThanUpperBody(midLastBand) || isOpenCloseAboveUpper(midLastBand)) {
			if countOpenCloseAboveUpper(short.Bands[bandLen-4:]) == 0 && countOpenCloseBelowSMA(short.Bands[bandLen-4:]) < 3 {
				if short.Position == models.ABOVE_UPPER || short.AllTrend.SecondTrend == models.TREND_UP {
					if mid.Position == models.ABOVE_UPPER || mid.AllTrend.SecondTrend == models.TREND_UP {
						if long.Position == models.ABOVE_UPPER || long.AllTrend.SecondTrend == models.TREND_UP {
							if !(isBandMultipleThanList(shortLastBand, short.Bands[bandLen-5:bandLen-1], 5) && bandPercent(shortLastBand) <= 4) {
								if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
									matchPattern = "first check up: general pattern 1"
									return true
								}
							}
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if isLastBandDoublePreviousHeigest(long.Bands) && bandPercent(longLastBand) > 1.5 && countOpenCloseAboveUpper(long.Bands[bandLen-2:]) == 0 {
					if isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 1.5 {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 && bandPercent(shortLastBand) < 7 && countOpenCloseBelowSMA(short.Bands[bandLen-4:]) < 3 {
							matchPattern = "first check up: pattern 1: apt base"
							return true
						}
					}
				}
			}
		}
	}

	if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
		//skipped
		if countOpenCloseBelowSMA(short.Bands[bandLen-4:]) == 3 {
			ignoredReason = "first check up: skipped1"
			return false
		}

		if long.Position == models.ABOVE_UPPER && isUpperHeadMoreThanUpperBody(longLastBand) {
			if mid.Position == models.ABOVE_UPPER && countOpenCloseAboveUpper(mid.Bands[bandLen-3:]) > 1 {
				if short.Position == models.ABOVE_UPPER && countHightCrossUpper(short.Bands[bandLen-4:]) == 1 {
					ignoredReason = "first check up: skipped2"
					return false
				}
			}
		}
	}

	if !isLastBandDoublePreviousHeigest(short.Bands) {
		if short.Position == models.ABOVE_UPPER || short.AllTrend.SecondTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER || mid.AllTrend.SecondTrend == models.TREND_UP {
				if countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-4:]) < 3 && countUpperHeadMoreThanUpperBody(short.Bands[bandLen-4:]) < 3 {
					if bandPercent(shortLastBand) > 1.5 && bandPercent(shortSecondLastBand) > 2.6 {
						matchPattern = "first check up: pattern x 2"
					}

				}
			}
		}
	}

	if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if CalculateTrendShort(long.Bands[bandLen-5:bandLen-1], true) == models.TREND_UP && countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-4:]) <= 1 {
					if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) <= 1 || countCrossUpUpperOnBody(short.Bands[bandLen-4:]) <= 1 {
						if isLastBandMultiplePreviousHeigest(short.Bands, 4) && bandPercent(shortLastBand) > 5 {
							matchPattern = "first check not up: pattern 1: t base"
							return true
						}
					}
				}

				if countOpenCloseBelowSMA(short.Bands[bandLen-4:]) == 0 && !(isUpperHeadMoreThanUpperBody(midLastBand) || isOpenCloseAboveUpper(midLastBand)) {
					if !isBandMultipleThanList(short.Bands[bandLen-2], short.Bands[bandLen-6:bandLen-2], 2) {
						if !isUpperHeadMoreThanUpperBody(shortLastBand) && !isUpperHeadMoreThanUpperBody(longLastBand) {
							if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(mid.Bands[bandLen-2:]) == 1 && countCrossUpUpperOnBody(long.Bands[bandLen-2:]) == 1 {
								if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
									matchPattern = "first check not up: pattern 1: jasmy base"
									return true
								}
							}
						}
					}
				}

				if countOpenCloseBelowSMA(long.Bands[bandLen-4:]) == 0 && countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 1 && !isUpperHeadMoreThanUpperBody(longLastBand) {
					if countOpenCloseBelowSMA(mid.Bands[bandLen-4:]) == 0 && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && !isUpperHeadMoreThanUpperBody(midLastBand) {
						if countOpenCloseBelowSMA(short.Bands[bandLen-3:]) == 0 && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && !isUpperHeadMoreThanUpperBody(midLastBand) {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
								matchPattern = "first check not up: pattern 1: apt base"
								return true
							}
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if !isOpenCloseAboveUpper(midLastBand) && !isUpperHeadMoreThanUpperBody(midLastBand) {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
							matchPattern = "first check not up: pattern 1: santos base"
							return true
						}
					}
				}
			}
		}
	}

	if (long.Position == models.ABOVE_SMA || long.Position == models.BELOW_SMA) && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if !isBandMultipleThanList(short.Bands[bandLen-2], short.Bands[bandLen-6:bandLen-2], 2) {
					if !(isUpperHeadMoreThanUpperBody(midLastBand) || isOpenCloseAboveUpper(midLastBand)) {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
							matchPattern = "first check not up: pattern 1: ctxc base"
							return true
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countOpenCloseBelowSMA(long.Bands[bandLen-4:]) == 0 {
					if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countOpenCloseBelowSMA(mid.Bands[bandLen-4:]) == 0 {
						if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && shortLastBand.Candle.Open < float32(shortLastBand.SMA) {
							matchPattern = "first check not up: pattern 1: enj base"
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func approvedPatternOnCompleteCheck(short, mid, long models.BandResult, currentTime time.Time) bool {
	bandLen := len(short.Bands)
	longLastBand := long.Bands[bandLen-1]
	longSecondLastBand := long.Bands[bandLen-2]
	midLastBand := mid.Bands[bandLen-1]
	midSecondLastBand := mid.Bands[bandLen-2]
	shortLastBand := short.Bands[bandLen-1]
	shortSecondLastBand := short.Bands[bandLen-2]

	//skipp
	if isUpperHeadMoreThanUpperBody(midLastBand) && isUpperHeadMoreThanUpperBody(longLastBand) {
		ignoredReason = "mid and long upper head more than upper body"
		return false
	}

	if CalculateShortTrendWithConclusion(mid.Bands[:bandLen-1]) == models.TREND_DOWN && CalculateShortTrendWithConclusion(long.Bands[:bandLen-1]) == models.TREND_DOWN {
		ignoredReason = "mid and long previous short trend down"
		return false
	}

	if short.Position == models.ABOVE_UPPER && !isUpperHeadMoreThanUpperBody(shortLastBand) {
		if mid.Position == models.ABOVE_UPPER || mid.AllTrend.SecondTrend == models.TREND_UP {
			if long.Position == models.ABOVE_UPPER && isBandLongTimeComplete(currentTime) && !isUpperHeadMoreThanUpperBody(longLastBand) {
				if isLastBandDoublePreviousHeigest(long.Bands) && !isLastBandDoublePreviousHeigest(long.Bands[:bandLen-1]) && isSolidBand(longLastBand) {
					if isSolidBand(shortLastBand) && isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
						matchPattern = "band complete up: avax base"
						return true
					}
				}
			}
		}
	}

	if currentTime.Minute() > 58 || currentTime.Minute() < 3 {
		if (currentTime.Minute() < 3 && currentTime.Hour()%4 == 0) || (currentTime.Minute() > 58 && currentTime.Hour()%4 == 3) {
			if isSolidBand(longLastBand) && !isBandDown(longSecondLastBand) {
				if isLastBandDoublePreviousHeigest(long.Bands) && bandPercent(longLastBand) > 3.6 {
					matchPattern = "band complete up: solid long band pattern"
					return true
				}
			}
		} else {
			if isSolidBand(midLastBand) && !isBandDown(midSecondLastBand) {
				if countUpperHeadMoreThanUpperBody(short.Bands[bandLen-2:]) <= 1 {
					if isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 3 {
						matchPattern = "band complete up: solid mid band pattern"
						return true
					}
				}
			}
		}
	} else {
		if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP && long.Position == models.ABOVE_UPPER && countHightCrossUpper(long.Bands[bandLen-4:]) > 1 {
			if isBadBand(longSecondLastBand) || CalculateShortTrendWithConclusion(mid.Bands[:bandLen-1]) != models.TREND_UP {
				if mid.Position == models.ABOVE_UPPER && countHightCrossUpper(mid.Bands[bandLen-4:]) == 1 {
					if countOpenCloseBelowSMA(short.Bands[bandLen-10:]) > 1 {
						ignoredReason = "cross upper mentok kayaknya"
						return false
					}
				}
			}
		}
		if isSolidBand(shortLastBand) && !isBandDown(shortSecondLastBand) {
			if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
				matchPattern = "band complete up: solid short band pattern"
				return true
			}
		}
	}

	if !isLastBandDoublePreviousHeigest(short.Bands) {
		if short.Position == models.ABOVE_UPPER || short.AllTrend.SecondTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER || mid.AllTrend.SecondTrend == models.TREND_UP {
				if countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-4:]) < 3 && countUpperHeadMoreThanUpperBody(short.Bands[bandLen-4:]) < 3 {
					if bandPercent(shortLastBand) > 2.6 && bandPercent(shortSecondLastBand) > 2.6 {
						matchPattern = "band complete up: pattern x 2"
						return true
					}
				}
			}
		}
	}

	if !isLastBandDoublePreviousHeigest(short.Bands) {
		if short.Position == models.ABOVE_UPPER || short.AllTrend.SecondTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER || mid.AllTrend.SecondTrend == models.TREND_UP {
				if countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-4:]) < 3 && countUpperHeadMoreThanUpperBody(short.Bands[bandLen-4:]) < 3 {
					if isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 3 {
						matchPattern = "band complete up: pattern x 3"
						return true
					}
				}
			}
		}
	}

	if long.Position == models.BELOW_SMA && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && !isUpperHeadMoreThanUpperBody(shortLastBand) {
					if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && !isUpperHeadMoreThanUpperBody(midLastBand) {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
							matchPattern = "band complete not up: pattern 1:perl base"
							return true
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if !isUpperHeadMoreThanUpperBody(shortLastBand) {
					if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
						matchPattern = "band complete not up: pattern 1:dent base"
						return true
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_SIDEWAY && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && isLastBandDoublePreviousHeigest(mid.Bands) {
					if !(isUpperHeadMoreThanUpperBody(midLastBand) || isOpenCloseAboveUpper(midLastBand)) {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
							matchPattern = "band complete not up: pattern 1:lsk base"
							return true
						}
					}
				}

			}
		}
	}

	if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countOpenCloseBelowSMA(long.Bands[bandLen-6:]) == 0 {
					if countOpenCloseAboveUpper(mid.Bands[bandLen-4:]) == 0 && countOpenCloseAboveUpper(short.Bands[bandLen-4:]) == 0 {
						if countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-4:bandLen-1]) == 0 && countUpperHeadMoreThanUpperBody(short.Bands[bandLen-4:bandLen-1]) == 0 {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:fet2 base"
								return true
							}
						}
					}
				}

				if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
					if !isOpenCloseAboveUpper(midLastBand) && !isUpperHeadMoreThanUpperBody(midLastBand) {
						if isDoubleSignificanUp(short.Bands) && bandPercent(shortLastBand) < 4 {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:magic base"
								return true
							}
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 0 {
					if countOpenCloseAboveUpper(mid.Bands[bandLen-4:]) == 0 && countUpperHeadMoreThanUpperBody(mid.Bands[bandLen-4:]) == 0 {
						if countOpenCloseBelowSMA(long.Bands[bandLen-4:]) == 0 && countHightCrossUpper(mid.Bands[bandLen-4:]) == 1 {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:fis base"
								return true
							}
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
		if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if !isUpperHeadMoreThanUpperBody(longLastBand) && !isOpenCloseAboveUpper(longLastBand) {
					if !isUpperHeadMoreThanUpperBody(midLastBand) && !isOpenCloseAboveUpper(midLastBand) {
						if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 && bandPercent(shortLastBand) < 7 {
								matchPattern = "band complete not up: pattern 1:fet3 base"
								return true
							}
						}
					}
				}

				if countCrossUpUpperOnBody(long.Bands[bandLen-4:]) > 1 && !isUpperHeadMoreThanUpperBody(longLastBand) && !isOpenCloseAboveUpper(longLastBand) {
					if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
						if countHightCrossUpper(long.Bands[bandLen-4:]) < 4 && !isUpperHeadMoreThanUpperBody(longLastBand) {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:sol base"
								return true
							}
						}
					}
				}
			}
		}
	}

	if long.Position == models.ABOVE_SMA && bandPercentFromUpper(longLastBand) > 3 {
		if isLastBandDoublePreviousHeigest(long.Bands) && countOpenCloseBelowSMA(long.Bands[bandLen-4:]) > 2 {
			if isLastBandDoublePreviousHeigest(mid.Bands) && !isLastBandDoublePreviousHeigest(mid.Bands[:bandLen-1]) {
				if isSolidBand(midLastBand) {
					if isLastBandDoublePreviousHeigest(short.Bands) && isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 && bandPercent(shortLastBand) < 5 {
							matchPattern = "band complete not up: pattern 1:fet1 base"
							return true
						}
					}
				}
			}
		}
	}

	if short.Position == models.ABOVE_SMA && bandPercentFromUpper(shortLastBand) > 4 {
		if countHightCrossUpper(mid.Bands[bandLen-4:]) > 1 && isBandDown(midSecondLastBand) {
			if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
				matchPattern = "band complete not up: pattern 1:cfx base1"
				return true
			}
		}
	}

	if long.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 1 {
		if mid.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
			if countOpenCloseBelowSMA(mid.Bands[bandLen-4:]) == 0 && countOpenCloseBelowSMA(long.Bands[bandLen-4:]) == 0 {
				if short.Position == models.ABOVE_UPPER && countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 {
					if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
						matchPattern = "band complete not up: pattern 1:cfx base3"
						return true
					}
				}
			}
		}
	}

	if currentTime.Minute() > 58 || currentTime.Minute() < 3 {
		if (currentTime.Minute() < 3 && currentTime.Hour()%4 == 0) || (currentTime.Minute() > 58 && currentTime.Hour()%4 == 3) {
			if isSolidBand(longLastBand) && !isBandDown(longSecondLastBand) {
				if isLastBandDoublePreviousHeigest(long.Bands) && bandPercent(longLastBand) > 3.6 {
					matchPattern = "band complete not up: solid long band pattern"
					return true
				}
			}
		} else {
			if isSolidBand(midLastBand) && !isBandDown(midSecondLastBand) {
				if isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 3 {
					matchPattern = "band complete not up: solid mid band pattern"
					return true
				}
			}
		}
	} else {
		if isSolidBand(shortLastBand) && !isBandDown(shortSecondLastBand) {
			if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
				matchPattern = "band complete not up: solid short band pattern"
				return true
			}
		}
	}

	return false
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

func isOnFirstCheck(currentTime time.Time) bool {
	minute := currentTime.Minute()

	if minute == 4 || minute == 5 {
		return true
	} else if minute == 19 || minute == 20 {
		return true
	} else if minute == 34 || minute == 35 {
		return true
	} else if minute == 49 || minute == 50 {
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
