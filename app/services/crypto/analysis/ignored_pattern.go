package analysis

import (
	"log"
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""

func lowestFromBand(band models.Band) float32 {
	if band.Candle.Open > band.Candle.Close {
		return band.Candle.Close
	}

	return band.Candle.Open
}

func GetIgnoredReason() string {
	return ignoredReason
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
	lastBandBodyHeight := lastBand.Candle.Close - lastBand.Candle.Open

	var higestBody float32 = 0
	for _, band := range bands[len(bands)-5 : len(bands)-1] {
		bodyHeight := band.Candle.Close - band.Candle.Open
		if band.Candle.Close < band.Candle.Open {
			bodyHeight = band.Candle.Open - band.Candle.Close
		}

		if bodyHeight > higestBody {
			higestBody = bodyHeight
		}
	}

	return higestBody*2 < lastBandBodyHeight
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

func isOpenCloseBelowLower(band models.Band) bool {
	return band.Candle.Open < float32(band.Lower) && band.Candle.Close < float32(band.Lower)
}

func countOpenCloseBelowLower(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isOpenCloseBelowLower(band) {
			count++
		}
	}
	return count
}

func isHasOpenCloseBelowLower(bands []models.Band) bool {
	for _, band := range bands {
		if isOpenCloseBelowLower(band) {
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
	head := band.Candle.Hight - band.Candle.Close
	body := band.Candle.Close - band.Candle.Open

	return head > body
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

	return upperToHead >= upperToBody
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

func isHightCrossUpper(band models.Band) bool {
	return band.Candle.Hight > float32(band.Upper)
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
	return isHeadMoreThanBody(band) || band.Candle.Open > band.Candle.Close || isTailMoreThanBody(band)
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

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time) bool {
	ignoredReason = ""

	bandLen := len(short.Bands)
	shortLastBand := short.Bands[bandLen-1]
	shortSecondLastBand := short.Bands[bandLen-2]
	midLastBand := mid.Bands[bandLen-1]
	midSecondLastBand := mid.Bands[bandLen-2]
	longLastBand := long.Bands[bandLen-1]

	if isUpperHeadMoreThanUpperBody(shortLastBand) || isOpenCloseAboveUpper(shortLastBand) {
		if isUpperHeadMoreThanUpperBody(midLastBand) || isOpenCloseAboveUpper(midLastBand) {
			log.Println("skipped1")
			return false
		}
	}

	if currentTime.Minute() < 15 && bandPercent(shortLastBand) < 1.5 {
		log.Println("skipped2")
		return false
	}

	if long.AllTrend.ShortTrend == models.TREND_DOWN && mid.Position == models.BELOW_SMA {
		log.Println("skipped3")
		return false
	}

	if longLastBand.Candle.Low < float32(longLastBand.Lower) {
		if short.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER {
			log.Println("skipped4")
			return false
		}
	}

	if countOpenCloseAboveUpper(short.Bands[bandLen-2:]) > 1 {
		log.Println("skipped5")
		return false
	}

	if isUpperHeadMoreThanUpperBody(shortLastBand) && (long.AllTrend.ShortTrend != models.TREND_UP || long.PriceChanges < 0) {
		log.Println("skipped6")
		return false
	}

	if long.AllTrend.SecondTrend == models.TREND_DOWN && long.Position == models.BELOW_SMA && countOpenCloseBelowLower(long.Bands[bandLen-4:]) > 0 {
		log.Println("skipped7")
		return false
	}

	if isHightCrossUpper(shortLastBand) && isHightCrossUpper(midLastBand) && isHightCrossUpper(longLastBand) {
		if CountBadBand(short.Bands[bandLen-4:]) > 2 && short.AllTrend.SecondTrend == models.TREND_UP {
			log.Println("skipped8")
			return false
		}
	}

	if isHasOpenCloseBelowLower(short.Bands[bandLen-4:]) || isHasOpenCloseBelowLower(mid.Bands[bandLen-4:]) {
		log.Println("skipped9")
		return false
	}

	if long.AllTrend.SecondTrend == models.TREND_DOWN && long.Position == models.BELOW_SMA {
		if mid.AllTrend.SecondTrend != models.TREND_UP && (mid.Position == models.ABOVE_SMA || (mid.Position == models.ABOVE_UPPER && midLastBand.Candle.Open < float32(midLastBand.SMA))) {
			log.Println("skipped10")
			return false
		}
	}

	if isShortBandComplete(currentTime) {
		if short.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER {
			if countBadBandAndCrossUpper(short.Bands[bandLen-4:]) > 0 && (countBadBandAndCrossUpper(mid.Bands[bandLen-4:]) > 0 || countBadBandAndCrossUpper(long.Bands[bandLen-4:]) > 0) {
				log.Println("band complete: skipped1")
				return false
			}

			if countBadBandAndCrossUpper(short.Bands[bandLen-4:]) > 1 && countHightCrossUpper(mid.Bands[bandLen-4:]) == 1 && countHightCrossUpper(long.Bands[bandLen-4:]) == 1 {
				log.Println("band complete: skipped1.1")
				return false
			}

			if countBadBandAndCrossUpper(long.Bands[bandLen-4:]) > 2 {
				if countBadBandAndCrossUpper(short.Bands[bandLen-4:]) == 1 && countBadBandAndCrossUpper(mid.Bands[bandLen-4:]) == 1 {
					log.Println("band complete: skipped1.2")
					return false
				}
			}
		}

		if short.Position == models.ABOVE_SMA && mid.AllTrend.SecondTrend == models.TREND_UP {
			if CountBadBand(mid.Bands[bandLen-4:]) > 2 {
				log.Println("band complete: skipped2")
				return false
			}
		}

		if short.Position == models.ABOVE_SMA && mid.Position == models.ABOVE_SMA {
			if short.AllTrend.Trend == models.TREND_DOWN && mid.AllTrend.Trend == models.TREND_DOWN {
				log.Println("band complete: skipped3")
				return false
			}
		}

		if long.AllTrend.SecondTrend == models.TREND_DOWN {
			if midLastBand.Candle.Open < float32(midLastBand.SMA) && countCrossUpSMAOnBody(mid.Bands[bandLen-5:bandLen-1]) == 0 {
				if mid.Position == models.ABOVE_UPPER && short.Position == models.ABOVE_UPPER {
					log.Println("band complete: skipped4")
					return false
				}
			}
		}

		if mid.Position == models.ABOVE_UPPER && countHightCrossUpper(mid.Bands[bandLen-4:]) == 1 && !isLastBandDoublePreviousHeigest(mid.Bands) {
			if currentTime.Minute() > 55 || currentTime.Minute() < 5 {
				log.Println("band complete: skipped5")
				return false
			}
		}

		if long.Position == models.ABOVE_UPPER && long.AllTrend.ShortTrend != models.TREND_UP {
			if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP && CountBadBand(mid.Bands[bandLen-4:]) > 2 {
				if currentTime.Minute() > 55 || currentTime.Minute() < 5 {
					log.Println("band complete: skipped6")
					return false
				}
			}
		}

		if (long.Position == models.ABOVE_UPPER || isHightCrossUpper(longLastBand)) && long.AllTrend.ShortTrend == models.TREND_UP && countHightCrossUpper(long.Bands[bandLen-4:]) == 1 {
			if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP && CountBadBand(mid.Bands[bandLen-4:]) > 2 {
				if isHightCrossUpper(midLastBand) {
					log.Println("band complete: skipped7")
					return false
				}
			}
		}

		if long.AllTrend.SecondTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP && CountBadBand(long.Bands[bandLen-4:]) > 2 {
			if mid.AllTrend.SecondTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP && CountBadBand(mid.Bands[bandLen-4:]) > 2 {
				if isHightCrossUpper(midLastBand) {
					log.Println("band complete: skipped8")
					return false
				}
			}
		}

		if long.Position == models.ABOVE_UPPER && mid.Position == models.ABOVE_UPPER && short.Position == models.ABOVE_UPPER {
			if isUpperHeadMoreThanUpperBody(shortLastBand) {
				if currentTime.Minute() > 55 || currentTime.Minute() < 5 {
					log.Println("band complete: skipped9")
					return false
				}
			}
		}

		if short.Position == models.ABOVE_UPPER && countHightCrossUpper(short.Bands[bandLen-4:]) == 1 {
			if mid.AllTrend.SecondTrend != models.TREND_UP || long.AllTrend.SecondTrend != models.TREND_UP {
				if currentTime.Minute() > 55 || currentTime.Minute() < 5 {
					log.Println("band complete: skipped10")
					return false
				}
			}

			if mid.Position == models.ABOVE_UPPER && long.Position == models.ABOVE_UPPER && CountBadBand(mid.Bands[bandLen-4:]) > 2 {
				log.Println("band complete: skipped11")
				return false
			}

			if CountBadBand(short.Bands[bandLen-4:]) > 2 {
				log.Println("band complete: skipped12")
				return false
			}
		}

		if shortSecondLastBand.Candle.Open > shortSecondLastBand.Candle.Close || !isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) {
			if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 2.6 {
				if !(isHeadMoreThanBody(midSecondLastBand) && isHightCrossUpper(midSecondLastBand)) && !isOpenCloseAboveUpper(midLastBand) {
					ignoredReason = "band complete: pattern 1"
					return true
				}
			}
		}

		return false
	}

	if shortSecondLastBand.Candle.Open > shortSecondLastBand.Candle.Close || !isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) {
		if isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5 {
			if !(isHeadMoreThanBody(midSecondLastBand) && isHightCrossUpper(midSecondLastBand)) && !isOpenCloseAboveUpper(midLastBand) {
				ignoredReason = "pattern 1"
				return true
			}
		}
	}

	if !(isLastBandDoublePreviousHeigest(short.Bands) && bandPercent(shortLastBand) > 1.5) && !(isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) && bandPercent(shortSecondLastBand) > 1.5) {
		if countBadBandAndCrossUpper(short.Bands[bandLen-3:]) <= 1 && !isOpenCloseAboveUpper(shortLastBand) {
			if midSecondLastBand.Candle.Open > midSecondLastBand.Candle.Close || !isLastBandDoublePreviousHeigest(mid.Bands[:bandLen-1]) {
				if isLastBandDoublePreviousHeigest(mid.Bands) && bandPercent(midLastBand) > 3 && isLastBandHeigestBand(short.Bands, 4) {
					if !isOpenCloseAboveUpper(midLastBand) && !isHeadMoreThanBody(midSecondLastBand) {
						if !isAboveUpperAndOrUpperHeadMoreThanUpperBody(shortLastBand, short.Bands[bandLen-3:bandLen-1]) {
							ignoredReason = "pattern 2"
							return true
						}
					}
				}
			}
		}
	}

	if bandPercent(shortLastBand) < bandPercent(shortSecondLastBand) {
		if bandPercent(shortSecondLastBand)-bandPercent(shortLastBand) > 3.1 {
			if shortSecondLastBand.Candle.Close > shortSecondLastBand.Candle.Open {
				if isLastBandDoublePreviousHeigest(short.Bands[:bandLen-1]) && bandPercent(shortSecondLastBand) > 1.5 {
					if !isOpenCloseAboveUpper(midLastBand) && !isHeadMoreThanBody(midSecondLastBand) {
						if !isOpenCloseAboveUpper(shortLastBand) && !isOpenCloseAboveUpper(longLastBand) {
							ignoredReason = "pattern 3"
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func isShortBandComplete(currentTime time.Time) bool {
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
