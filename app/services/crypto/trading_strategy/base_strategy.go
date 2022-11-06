package trading_strategy

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/crypto/analysis"
	"time"
)

const LIMIT_COIN_CHECK int = 60

var countTrendUp int = 0
var checkOnTrendUpLimit int = 25

type TradingStrategy interface {
	Execute(currentTime time.Time)
	InitService()
	Shutdown()
}

func checkCryptoHoldCoinPrice(requestTime time.Time) []models.BandResult {
	log.Println("starting crypto check price hold coin worker")

	holdCoin := []models.BandResult{}

	endDate := GetEndDate(requestTime, OPERATION_SELL)

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

		result := crypto.MakeCryptoRequest(request)
		if result == nil {
			continue
		}

		holdCoin = append(holdCoin, *result)
	}

	return holdCoin
}

func checkCoinOnTrendUp(baseTime time.Time, previousResult map[string]*models.BandResult) []models.BandResult {
	altCoin := []models.BandResult{}

	endDate := GetEndDate(baseTime, OPERATION_BUY)

	responseChan := make(chan crypto.CandleResponse)

	if checkOnTrendUpLimit == 0 {
		log.Println("skip process check on trend up limit is zero")
		return altCoin
	}

	limit := checkOnTrendUpLimit + 1

	condition := map[string]interface{}{"is_master": false, "is_on_hold": false}
	orderBy := "price_changes desc"
	currencyConfigs := repositories.GetCurrencyNotifConfigs(&condition, &limit, &orderBy)

	for _, data := range *currencyConfigs {
		if data.PriceChanges < 1 {
			continue
		}

		var result *models.BandResult
		if resultLoc, ok := previousResult[data.Symbol]; ok {
			result = resultLoc
		} else {
			request := crypto.CandleRequest{
				Symbol:       data.Symbol,
				EndDate:      endDate,
				Limit:        40,
				Resolution:   "15m",
				ResponseChan: responseChan,
			}

			result = crypto.MakeCryptoRequest(request)
		}

		if result == nil || result.Direction == analysis.BAND_DOWN || result.AllTrend.ShortTrend != models.TREND_UP || result.AllTrend.SecondTrend != models.TREND_UP {
			continue
		}

		altCoin = append(altCoin, *result)
	}

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
