package analysis

import (
	"telebot-trading/app/models"
	"time"
)

var ignoredReason string = ""

func isLastBandCrossUpperAndPreviousBandNot(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	if lastBand.Candle.Open < lastBand.Candle.Close {
		secondLastBand := bands[len(bands)-2]
		return secondLastBand.Candle.Open > secondLastBand.Candle.Close && secondLastBand.Candle.Low > float32(secondLastBand.Lower)
	}
	return false
}

func isHasCrossUpper(bands []models.Band, withHead bool) bool {
	for _, band := range bands {
		if band.Candle.Close < band.Candle.Open {
			continue
		}

		if withHead {
			if band.Candle.Open < float32(band.Upper) && band.Candle.Hight >= float32(band.Upper) {
				return true
			}
			if band.Candle.Open > band.Candle.Close {
				if band.Candle.Close < float32(band.Upper) && band.Candle.Hight >= float32(band.Upper) {
					return true
				}
			}
		} else {
			if band.Candle.Open <= float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
				return true
			}
		}
	}
	return false
}

func isHasBandDownFromUpper(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open > float32(band.Upper) && band.Candle.Close < band.Candle.Open {
			return true
		}
	}

	return false
}

func isHasCrossSMA(bands []models.Band, bodyOnly bool) bool {
	for _, band := range bands {
		if bodyOnly {
			if band.Candle.Open < float32(band.SMA) && band.Candle.Close > float32(band.SMA) {
				return true
			}
		} else {
			if band.Candle.Low < float32(band.SMA) && band.Candle.Hight > float32(band.SMA) {
				return true
			}
		}
	}
	return false
}

func isHasOpenCloseAboveUpper(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Open > float32(band.Upper) && band.Candle.Close > float32(band.Upper) {
			return true
		}
	}
	return false
}

func isHasCrossLower(bands []models.Band, bodyOnly bool) bool {
	crossLowerBand := false
	for _, data := range bands {
		if bodyOnly {
			if (data.Candle.Close < float32(data.Lower) && data.Candle.Open > float32(data.Lower)) || (data.Candle.Open < float32(data.Lower) && data.Candle.Close > float32(data.Lower)) {
				crossLowerBand = true
				break
			}
		} else {
			if data.Candle.Low < float32(data.Lower) && (data.Candle.Close > float32(data.Lower) || data.Candle.Open > float32(data.Lower)) {
				crossLowerBand = true
				break
			}
		}
	}

	return crossLowerBand
}

func countCrossUpper(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Open <= float32(data.Upper) && data.Candle.Hight > float32(data.Upper) && (data.Candle.Close > data.Candle.Open) {
			if !isHasBadBand([]models.Band{data}) {
				count++
			}
		}
	}

	return count
}

func countCrossUpperOnBody(bands []models.Band) int {
	count := 0
	for _, data := range bands {
		if data.Candle.Open <= float32(data.Upper) && data.Candle.Close >= float32(data.Upper) && (data.Candle.Close > data.Candle.Open) {
			if !isHasBadBand([]models.Band{data}) {
				count++
			}
		}
	}

	return count
}

func countBelowSMA(bands []models.Band, strict bool) int {
	count := 0
	for _, data := range bands {
		if strict {
			if data.Candle.Close < float32(data.SMA) && data.Candle.Open < float32(data.SMA) {
				count++
			}
		} else {
			if data.Candle.Close < float32(data.SMA) {
				count++
			}
		}
	}

	return count
}

func countBelowLower(bands []models.Band, strict bool) int {
	count := 0
	for _, data := range bands {
		if strict {
			if data.Candle.Hight < float32(data.Lower) && data.Candle.Low < float32(data.Lower) {
				count++
			}
		} else {
			if data.Candle.Close < float32(data.Lower) && data.Candle.Open < float32(data.Lower) {
				count++
			}
		}
	}

	return count
}

func getHighestIndex(bands []models.Band) int {
	hiIndex := 0
	for i, band := range bands {
		if bands[hiIndex].Candle.Close <= band.Candle.Close {
			hiIndex = i
		}
	}

	return hiIndex
}

func getHighestHightIndex(bands []models.Band) int {
	hiIndex := 0
	for i, band := range bands {
		if bands[hiIndex].Candle.Hight <= band.Candle.Hight {
			hiIndex = i
		}
	}

	return hiIndex
}

func lowestFromBand(band models.Band) float32 {
	if band.Candle.Open > band.Candle.Close {
		return band.Candle.Close
	}

	return band.Candle.Open
}

func getHigestIndexSecond(bands []models.Band) int {
	firstHight := getHighestIndex(bands)

	secondHight := -1
	for i, band := range bands {
		if i != firstHight {
			if secondHight < 0 {
				secondHight = i
			} else if bands[secondHight].Candle.Close < band.Candle.Close {
				secondHight = i
			}
		}
	}

	return secondHight
}

func getIndexHigestCrossUpper(bands []models.Band) int {
	higestIndex := -1
	lastBand := bands[len(bands)-1]
	for i, band := range bands {
		if band.Candle.Close > lastBand.Candle.Close {
			if higestIndex != -1 {
				if bands[higestIndex].Candle.Close < band.Candle.Close {
					higestIndex = i
				}
			} else {
				higestIndex = i
			}
		}
	}

	return higestIndex
}

func countDownBand(bands []models.Band) int {
	counter := 0
	for _, band := range bands {
		if isBandDown(band) {
			counter++
		}
	}

	return counter
}

func isBandDown(band models.Band) bool {
	return band.Candle.Open > band.Candle.Close
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

func isContainHeadMoreThanBody(bands []models.Band) bool {
	for _, band := range bands {
		if band.Candle.Close > band.Candle.Open {
			head := band.Candle.Hight - band.Candle.Close
			body := band.Candle.Close - band.Candle.Open

			if head > body {
				return true
			}
		}
	}

	return false
}

func countHeadMoreThanBody(bands []models.Band) int {
	count := 0
	for _, band := range bands {
		if isHeadMoreThanBody(band) {
			count++
		}
	}

	return count
}

func isHeadMoreThanBody(band models.Band) bool {
	head := band.Candle.Hight - band.Candle.Close
	body := band.Candle.Close - band.Candle.Open
	if band.Candle.Open > band.Candle.Close {
		return true
	}

	return head > body
}

func isHasBadBand(bands []models.Band) bool {
	if countDownBand(bands) > 0 || countAboveUpper(bands) > 0 || isTailMoreThan(bands[len(bands)-1], 50) || countHeadMoreThanBody(bands) > 0 {
		return true
	}

	return false
}

func isUpperHeadMoreThanUpperBody(band models.Band) bool {
	allBody := band.Candle.Close - band.Candle.Open
	head := band.Candle.Close - float32(band.Upper)
	percent := head / allBody * 100
	return percent > 51
}

func countBandPercentChangesMoreThan(bands []models.Band, percent float32) int {
	count := 0
	for _, band := range bands {
		if band.Candle.Open < band.Candle.Close {
			percentChanges := (band.Candle.Close - band.Candle.Open) / band.Candle.Open * 100
			if percentChanges >= percent {
				count++
			}
		}
	}

	return count
}

func isTailMoreThan(band models.Band, percent float32) bool {
	bodyTail := band.Candle.Close - band.Candle.Low
	tail := band.Candle.Open - band.Candle.Low
	tailPercent := tail / bodyTail * 100

	return tailPercent > percent
}

func isHasUpperHeadMoreThanUpperBody(bands []models.Band) bool {
	for _, band := range bands {
		if isUpperHeadMoreThanUpperBody(band) {
			return true
		}
	}

	return false
}

func lastBandDoublePreviousHeigest(bands []models.Band) bool {
	lastBand := bands[len(bands)-1]
	lastBandBodyHeight := lastBand.Candle.Close - lastBand.Candle.Open

	var higestBody float32 = 0
	for _, band := range bands[len(bands)-5 : len(bands)-1] {
		bodyHeight := band.Candle.Close - band.Candle.Open
		if bodyHeight > higestBody {
			higestBody = bodyHeight
		}
	}

	return higestBody*2 < lastBandBodyHeight
}

func bandPercent(band models.Band) float32 {
	return (band.Candle.Close - band.Candle.Open) / band.Candle.Open * 100
}

func isCrossUpperNotBadBand(bands []models.Band) bool {
	var isCrossUpper bool = false

	for _, band := range bands[len(bands)-4:] {
		if band.Candle.Hight > float32(band.Upper) {
			isCrossUpper = true
			if isBandDown(band) || isHeadMoreThanBody(band) {
				return false
			}
		}
	}

	return isCrossUpper
}

func ApprovedPattern(short, mid, long models.BandResult, currentTime time.Time) bool {
	ignoredReason = ""

	bandLen := len(short.Bands)
	shortLastBand := short.Bands[bandLen-1]
	midLastBand := mid.Bands[bandLen-1]

	if lastBandDoublePreviousHeigest(short.Bands) || lastBandDoublePreviousHeigest(mid.Bands) {
		if bandPercent(shortLastBand) > 3 || bandPercent(midLastBand) > 3 {

			if short.AllTrend.ShortTrend == models.TREND_UP && mid.AllTrend.ShortTrend == models.TREND_UP && long.AllTrend.ShortTrend == models.TREND_UP {
				if isCrossUpperNotBadBand(mid.Bands) && isCrossUpperNotBadBand(long.Bands) {
					ignoredReason = "pattern 1"
					return true
				}
			}

		}
	}

	return false
}
