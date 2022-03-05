package trading_strategy

import (
	"fmt"
	"log"
	"strings"
	"telebot-trading/app/helper"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

type TradingStrategy interface {
	Execute(currentTime time.Time)
	InitService()
	Shutdown()
}

func checkCryptoHoldCoinPrice(requestTime time.Time) []models.BandResult {
	log.Println("starting crypto check price hold coin worker")

	holdCoin := []models.BandResult{}

	endDate := GetEndDate(requestTime)

	responseChan := make(chan crypto.CandleResponse)

	condition := map[string]interface{}{"is_on_hold": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil, nil)

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Limit:        40,
			EndDate:      endDate,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := crypto.MakeCryptoRequest(data, request)
		if result == nil {
			continue
		}

		holdCoin = append(holdCoin, *result)
	}

	return holdCoin
}

func checkCryptoAltCoinPrice(baseTime time.Time) []models.BandResult {
	log.Println("starting crypto check price for alt coin worker")

	altCoin := []models.BandResult{}

	endDate := GetEndDate(baseTime)

	responseChan := make(chan crypto.CandleResponse)

	limit := 80
	condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
	ignoredCoins := getIgnoreCoin()
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, &limit, ignoredCoins)

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			EndDate:      endDate,
			Limit:        40,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := crypto.MakeCryptoRequest(data, request)
		if result == nil || result.Direction == analysis.BAND_DOWN {
			if (result.AllTrend.Trend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN) && result.AllTrend.ShortTrend == models.TREND_DOWN {
				ignoreCoin(result.Symbol)
			}

			continue
		}

		result.Weight = analysis.CalculateWeight(result)
		if !analysis.IsIgnored(result, baseTime) {
			altCoin = append(altCoin, *result)
		}

	}

	return altCoin
}

func GetEndDate(baseTime time.Time) int64 {
	unixTime := baseTime.Unix()

	if baseTime.Minute()%15 == 0 {
		unixTime = unixTime - (int64(baseTime.Second()) % 60)
	}

	unixTime = (unixTime - 1) * 1000

	return unixTime
}

func ignoreCoin(coinSymbol string) {
	st := helper.GetSimpleStore()
	coinString := st.Get("ignore_coins")
	if coinString != nil {
		coinSymbol = fmt.Sprintf("%s,%s", *coinString, coinSymbol)
	}

	st.Set("ignore_coins", coinSymbol)
}

func getIgnoreCoin() []string {
	st := helper.GetSimpleStore()
	coinString := st.Get("ignore_coins")
	if coinString != nil {
		return nil
	}

	return strings.Split(*coinString, ",")
}
