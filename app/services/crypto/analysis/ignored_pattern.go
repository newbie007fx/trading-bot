package analysis

import (
	"log"
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

func isLastBandHeigestBand(bands []models.Band, count int) bool {
	lastBand := bands[len(bands)-1]
	lastBandClose := lastBand.Candle.Close

	for _, band := range bands[len(bands)-count : len(bands)-1] {
		if band.Candle.Close > lastBandClose {
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
	return isHeadMoreThanBody(band) || band.Candle.Open > band.Candle.Close || isTailMoreThanBody(band) || band.Candle.Open == band.Candle.Close
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

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time, modeChecking string) bool {
	ignoredReason = ""
	matchPattern = ""

	if isOnFirstCheck(currentTime) {
		return approvedPatternFirstCheck(short, mid, long, modeChecking)
	}

	if isOnBandCompleteCheck(currentTime) {
		return approvedPatternOnCompleteCheck(short, mid, long, modeChecking, currentTime)
	}

	return false
}

func approvedPatternFirstCheck(short, mid, long models.BandResult, modeChecking string) bool {
	bandLen := len(short.Bands)
	shortLastBand := short.Bands[bandLen-1]
	midLastBand := mid.Bands[bandLen-1]
	longLastBand := long.Bands[bandLen-1]
	longSecondLastBand := long.Bands[bandLen-2]

	if modeChecking == models.MODE_TREND_UP {
		if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
			matchPattern = "first check up: pattern 1"
			return true
		}
	} else {
		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) <= 1 || countCrossUpUpperOnBody(short.Bands[bandLen-4:]) <= 1 {
						if isLastBandMultiplePreviousHeigest(short.Bands, 4) && bandPercent(shortLastBand) > 1.5 {
							matchPattern = "first check not up: pattern 1: t base"
							return true
						}
					}

					if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(mid.Bands[bandLen-2:]) == 1 && countCrossUpUpperOnBody(long.Bands[bandLen-2:]) == 1 {
						log.Println("masuk")
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
							matchPattern = "first check not up: pattern 1: jasmy base"
							return true
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

		// skipped
		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if isUpperHeadMoreThanUpperBody(longLastBand) && isUpperHeadMoreThanUpperBody(longSecondLastBand) {
						if countBadBandAndCrossUpper(long.Bands[bandLen-4:]) > 1 {
							ignoredReason = "first check not up: skipped 1"
							return false
						}
					}
				}
			}
		}

		if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
			matchPattern = "first check not up: pattern 1"
			return true
		}
	}

	return false
}

func approvedPatternOnCompleteCheck(short, mid, long models.BandResult, modeChecking string, currentTime time.Time) bool {
	bandLen := len(short.Bands)
	longLastBand := long.Bands[bandLen-1]
	midLastBand := mid.Bands[bandLen-1]
	shortLastBand := short.Bands[bandLen-1]
	shortSecondLastBand := short.Bands[bandLen-2]

	if modeChecking == models.MODE_TREND_UP {
		if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_SIDEWAY && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if !isOpenCloseAboveUpper(longLastBand) && !isUpperHeadMoreThanUpperBody(longLastBand) {
						if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
							matchPattern = "band complete up: pattern 1:ldo base"
							return true
						}
					}
				}
			}
		}

		if !isOpenCloseAboveUpper(longLastBand) && !isUpperHeadMoreThanUpperBody(longLastBand) {
			if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 && !isOpenCloseAboveUpper(longLastBand) {
				matchPattern = "band complete up: pattern 1"
				return true
			}
		}
	} else {
		if long.Position == models.BELOW_SMA && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
						matchPattern = "band complete not up: pattern 1:perl base"
						return true
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
					if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
						matchPattern = "band complete not up: pattern 1:lsk base"
						return true
					}

				}
			}
		}

		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_SIDEWAY && short.AllTrend.ShortTrend == models.TREND_UP {
					if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
						if countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 1 {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:fet1 base"
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
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:fis base"
								return true
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
								if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
									matchPattern = "band complete not up: pattern 1:fet3 base"
									return true
								}
							}
						}
					}

					if countCrossUpUpperOnBody(long.Bands[bandLen-4:]) > 1 && !isUpperHeadMoreThanUpperBody(longLastBand) && !isOpenCloseAboveUpper(longLastBand) {
						if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 {
							if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
								matchPattern = "band complete not up: pattern 1:sol base"
								return true
							}
						}
					}
				}
			}
		}

		//skipped
		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if countBadBandAndCrossUpper(mid.Bands[bandLen-4:]) > 0 && countBadBandAndCrossUpper(long.Bands[bandLen-4:]) > 0 {
						ignoredReason = "band complete not up: skipped1"
						return false
					}

					if countOpenCloseAboveUpper(mid.Bands[bandLen-4:]) > 0 {
						if CountBadBand(short.Bands[bandLen-4:]) > 2 {
							ignoredReason = "band complete not up: skipped1.1"
							return false
						}

						if CountBadBand(mid.Bands[bandLen-4:]) > 2 {
							ignoredReason = "band complete not up: skipped1.2"
							return false
						}
					}

					if countOpenCloseAboveUpper(short.Bands[bandLen-4:bandLen-1]) > 0 {
						if countCrossUpUpperOnBody(short.Bands[bandLen-4:bandLen-1]) > 1 {
							ignoredReason = "band complete not up: skipped1.3"
							return false
						}
					}

					if countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 1 {
						if CountBadBand(long.Bands[bandLen-4:]) > 2 {
							if isUpperHeadMoreThanUpperBody(midLastBand) {
								if countCrossUpUpperOnBody(short.Bands[bandLen-4:]) > 1 && isUpperHeadMoreThanUpperBody(shortLastBand) {
									ignoredReason = "band complete not up: skipped1.4"
									return false
								}
							}
						}
					}

					if isDoubleSignificanUp(long.Bands[:bandLen-1]) {
						ignoredReason = "band complete not up: skipped1.5"
						return false
					}

					if isUpperHeadMoreThanUpperBody(long.Bands[bandLen-1]) && isUpperHeadMoreThanUpperBody(mid.Bands[bandLen-1]) {
						if isBandMultipleThanList(shortLastBand, short.Bands[bandLen-6:bandLen-1], 5) {
							ignoredReason = "band complete not up: skipped1.6"
							return false
						}
					}

					if isLastBandDoublePreviousHeigest(long.Bands) && isLastBandDoublePreviousHeigest(mid.Bands) {
						if isLastBandDoublePreviousHeigest(short.Bands) {
							if currentTime.Minute() > 58 || currentTime.Minute() < 2 {
								ignoredReason = "band complete not up: skipped1.7"
								return false
							}
						}
					}

					if isDoubleSignificanUp(long.Bands) && isUpperHeadMoreThanUpperBody(longLastBand) {
						ignoredReason = "band complete not up: skipped1.8"
						return false
					}

					if isDoubleSignificanUp(mid.Bands) && isUpperHeadMoreThanUpperBody(midLastBand) {
						ignoredReason = "band complete not up: skipped1.9"
						return false
					}

					if countDownFromUpper(long.Bands[bandLen-4:]) > 0 && isUpperHeadMoreThanUpperBody(longLastBand) {
						ignoredReason = "band complete not up: skipped1.10"
						return false
					}
				}
			}
		}

		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_SMA && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if isDoubleSignificanUp(long.Bands[:bandLen-1]) {
						ignoredReason = "band complete not up: skipped2"
						return false
					}
				}
			}
		}

		if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_DOWN && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					if isUpperHeadMoreThanUpperBody(shortLastBand) {
						if currentTime.Minute() > 58 || currentTime.Minute() < 2 {
							ignoredReason = "band complete not up: skipped3"
							return false
						}
					}
				}
			}
		}

		if long.Position == models.ABOVE_SMA && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_DOWN {
			if mid.Position == models.ABOVE_SMA && mid.AllTrend.SecondTrend == models.TREND_DOWN && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
					ignoredReason = "band complete not up: skipped4"
					return false
				}
			}
		}

		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if mid.Position == models.ABOVE_UPPER && mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP {
				if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_SIDEWAY && short.AllTrend.ShortTrend == models.TREND_UP {
					if isDoubleSignificanUp(long.Bands) {
						if countDownFromUpper(mid.Bands[bandLen-4:]) > 0 {
							ignoredReason = "band complete not up: skipped5"
							return false
						}
					}
				}
			}
		}

		if long.Position == models.ABOVE_UPPER && long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
			if short.Position == models.ABOVE_UPPER && short.AllTrend.SecondTrend == models.TREND_UP && short.AllTrend.ShortTrend == models.TREND_UP {
				if countDownFromUpper(long.Bands[bandLen-4:]) > 0 {
					ignoredReason = "band complete not up: skipped6"
					return false
				}

				if mid.Position == models.ABOVE_UPPER && mid.AllTrend.ShortTrend == models.TREND_UP {
					if countCrossUpUpperOnBody(mid.Bands[bandLen-4:]) == 1 && countCrossUpUpperOnBody(long.Bands[bandLen-4:]) == 1 {
						if countOpenCloseBelowSMA(mid.Bands[bandLen-4:]) > 2 && countOpenCloseBelowSMA(long.Bands[bandLen-4:]) > 2 {
							if short.Position == models.ABOVE_UPPER && isBadBand(shortLastBand) {
								ignoredReason = "band complete not up: skipped6.1"
								return false
							}

						}
					}
				}
			}
		}

		if shortLastBand.Candle.Open != shortSecondLastBand.Candle.Close {
			ignoredReason = "band complete not up: skipped7"
			return false
		}

		if mid.AllTrend.ShortTrend == models.TREND_UP {
			if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
				matchPattern = "band complete not up: pattern 1"
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
