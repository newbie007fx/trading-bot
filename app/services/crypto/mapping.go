package crypto

import (
	"fmt"
	"telebot-trading/app/models"
	"telebot-trading/app/services/crypto/analysis"
)

func GenerateMsg(coinResult models.BandResult) string {
	format := "Coin name: <b>%s</b> \nDirection: <b>%s</b> \nPrice: <b>%f</b> \nVolume: <b>%f</b> \nTrend: <b>%s (%s(%.2f)+%s(%.2f))</b> \nShort Trend: <b>%s</b> \nPrice Changes: <b>%.2f%%</b> \nVolume Average Changes: <b>%.2f%%</b> \nPosition: <b>%s</b> \n"
	msg := fmt.Sprintf(format, coinResult.Symbol, DirectionString(coinResult.Direction), coinResult.CurrentPrice, coinResult.CurrentVolume, TrendString(coinResult.AllTrend.Trend), TrendString(coinResult.AllTrend.FirstTrend), coinResult.AllTrend.FirstTrendPercent, TrendString(coinResult.AllTrend.SecondTrend), coinResult.AllTrend.SecondTrendPercent, TrendString(coinResult.AllTrend.ShortTrend), coinResult.PriceChanges, coinResult.VolumeChanges, PositionString(coinResult.Position))
	return msg
}

func HoldCoinMessage(config models.CurrencyNotifConfig, result *models.BandResult) string {
	var changes float32

	if config.HoldPrice < result.CurrentPrice {
		changes = (result.CurrentPrice - config.HoldPrice) / config.HoldPrice * 100
	} else {
		changes = (config.HoldPrice - result.CurrentPrice) / config.HoldPrice * 100
	}

	format := "Hold status: \nHold price: <b>%f</b> \nBalance: <b>%f</b> \nCurrent price: <b>%f</b> \nChanges: <b>%.2f%%</b> \nEstimation in USDT: <b>%f</b> \n"
	msg := fmt.Sprintf(format, config.HoldPrice, config.Balance, result.CurrentPrice, changes, (result.CurrentPrice * config.Balance))

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
