package trading_strategy

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

var checkOnTrendUpLimit int = 25

type TradingStrategy interface {
	Execute(currentTime time.Time)
	InitService()
	Shutdown()
}

func checkCryptoHoldCoinPrice(requestTime time.Time) []models.BandResult {
	holdCoin := []models.BandResult{}

	endDate := GetEndDate(requestTime, OPERATION_SELL)

	responseChan := make(chan crypto.CandleResponse)

	condition := map[string]interface{}{"is_on_hold = ?": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil, nil, nil)

	for _, data := range *currency_configs {
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			Limit:        40,
			EndDate:      endDate,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result := crypto.MakeCryptoRequest(request)
		if result == nil {
			continue
		}

		holdCoin = append(holdCoin, *result)
	}

	return holdCoin
}

func checkCoinOnTrendUp(baseTime time.Time) []models.BandResult {
	altCoin := []models.BandResult{}

	endDate := GetEndDate(baseTime, OPERATION_BUY)

	responseChan := make(chan crypto.CandleResponse)

	currencyConfigs := &[]models.CurrencyNotifConfig{}

	limit := checkOnTrendUpLimit
	var priceThreshold float32 = 1.5
	strictChecking := false

	lastUpdate := baseTime.Unix() - (60 * 5)
	condition := map[string]interface{}{
		"is_on_hold = ?":    false,
		"price_changes > ?": priceThreshold,
		"status = ?":        models.STATUS_ACTIVE,
		"updated_at > ?":    lastUpdate,
	}
	orderBy := "price_changes desc"

	ignoredCoins := services.GetIgnoredCurrencies()

	if checkOnTrendUpLimit < 13 {
		strictChecking = true
		coins := crypto.GetListCoinUp()
		if len(coins) == 0 {
			log.Println("skip process check on trend up limit is zero")
			return altCoin
		}
		if ignoredCoins != nil {
			coins = diffCoin(*ignoredCoins, coins)
			if len(coins) == 0 {
				log.Println("skip process check on trend up limit is zero2")
				return altCoin
			}
		}

		limit = len(coins)
		condition["symbol in ?"] = coins
		currencyConfigs = repositories.GetCurrencyNotifConfigs(&condition, &limit, nil, &orderBy)
		priceThreshold = 3.5
	} else {
		if ignoredCoins != nil {
			log.Println("ignored coins: ", *ignoredCoins)
		}

		currencyConfigs = repositories.GetCurrencyNotifConfigsIgnoredCoins(&condition, &limit, ignoredCoins, &orderBy)
		priceThreshold = 2
	}

	log.Println("found: ", len(*currencyConfigs))
	log.Println("mode checking: ", modeChecking)

	coinsString := ""
	for _, data := range *currencyConfigs {
		coinsString = coinsString + ", " + data.Symbol

		var result *models.BandResult
		request := crypto.CandleRequest{
			Symbol:       data.Symbol,
			EndDate:      endDate,
			Limit:        40,
			Resolution:   "15m",
			ResponseChan: responseChan,
		}

		result = crypto.MakeCryptoRequest(request)

		if result == nil || result.Direction == analysis.BAND_DOWN || result.AllTrend.ShortTrend != models.TREND_UP || result.PriceChanges < priceThreshold || (strictChecking && (result.AllTrend.ShortTrend == models.TREND_DOWN || result.AllTrend.Trend == models.TREND_DOWN || result.AllTrend.SecondTrend == models.TREND_DOWN)) {
			if result.Direction == analysis.BAND_DOWN && result.AllTrend.SecondTrend == models.TREND_DOWN && result.AllTrend.ShortTrend == models.TREND_DOWN {
				services.SetIgnoredCurrency(result.Symbol, 1)
			}

			continue
		}

		altCoin = append(altCoin, *result)
	}
	log.Println("listCoin: ", coinsString)

	return altCoin
}

const OPERATION_SELL = 1
const OPERATION_BUY = 2

func GetEndDate(baseTime time.Time, operation int64) int64 {
	unixTime := baseTime.Unix()

	if baseTime.Minute()%15 == 0 {
		unixTime = unixTime - (int64(baseTime.Second()) % 60) - 1
	}

	if operation == OPERATION_BUY {
		unixTime = (unixTime - 120)
	}

	return unixTime * 1000
}

func diffCoin(ignoredCoins, upCoins []string) []string {
	difference := make([]string, 0)
	for _, upCoin := range upCoins {
		found := false
		for _, ignoreCoin := range ignoredCoins {
			if upCoin == ignoreCoin {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, upCoin)
		}
	}
	return difference
}
