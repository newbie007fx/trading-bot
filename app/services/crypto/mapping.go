package crypto

import (
	"fmt"
	"telebot-trading/app/models"
	"telebot-trading/app/services/crypto/analysis"
)

func GenerateMsg(coinResult models.BandResult) string {
	format := "Coin name: <b>%s</b> \nDirection: <b>%s</b> \nPrice: <b>%f</b> \nVolume: <b>%f</b> \nTrend: <b>%s</b> \nPrice Changes: <b>%.2f%%</b> \nVolume Average Changes: <b>%.2f%%</b> \nNotes: <b>%s</b> \nPosition: <b>%s</b> \n"
	msg := fmt.Sprintf(format, coinResult.Symbol, DirectionString(coinResult.Direction), coinResult.CurrentPrice, coinResult.CurrentVolume, TrendString(coinResult.Trend), coinResult.PriceChanges, coinResult.VolumeChanges, coinResult.Note, PositionString(coinResult.Position))
	return msg
}

func TrendString(trend int8) string {
	if trend == models.TREND_UP {
		return "trend up"
	} else if trend == models.TREND_DOWN {
		return "trend down"
	}

	return "trend sideway"
}

func DirectionString(direction int8) string {
	if direction == analysis.BAND_UP {
		return "UP"
	}

	return "DOWN"
}

func PositionString(position int8) string {
	if position == models.ABOVE_UPPER {
		return "above upper"
	} else if position == models.ABOVE_SMA {
		return "above sma"
	} else if position == models.BELOW_SMA {
		return "below sma"
	}

	return "below lower"
}
