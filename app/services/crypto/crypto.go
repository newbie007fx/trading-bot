package crypto

import (
	"fmt"
	"log"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/analysis"
	"telebot-trading/app/services/crypto/driver"
	"telebot-trading/utils"
	"time"
)

type CandleRequest struct {
	Symbol       string
	Limit        int
	EndDate      int64
	StartDate    int64
	Resolution   string
	ResponseChan chan CandleResponse
}

type CandleResponse struct {
	CandleData []models.CandleData
	Err        error
}

var canldeRequest chan CandleRequest
var previousTimeCheck time.Time = time.Now()
var thresholdPerMinute int64 = 140
var counter int64 = 0
var CandleLimit int64 = 40

func DispatchRequestJob(request CandleRequest) {
	canldeRequest <- request
}

func RequestCandleService() {
	canldeRequest = make(chan CandleRequest, 100)

	crypto := driver.GetCrypto()
	for request := range canldeRequest {
		checkCounter()

		response := CandleResponse{}
		response.CandleData, response.Err = crypto.GetCandlesData(request.Symbol, request.Limit, request.StartDate, request.EndDate, request.Resolution)

		request.ResponseChan <- response
	}

	defer func() {
		close(canldeRequest)
	}()
}

func checkCounter() {
	currentTime := time.Now()
	sleep := 0
	timeLeft := 59 - currentTime.Second()
	if counter == thresholdPerMinute {
		if timeLeft > 0 {
			sleep = int(timeLeft) + 2
		}
		counter = 0
	}

	if timeLeft == 0 || currentTime.Minute() != previousTimeCheck.Minute() {
		counter = 0
		sleep = 2
	}

	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
		previousTimeCheck = time.Now()
	}

	debug := utils.Env("debug", "false")
	if debug == "true" {
		log.Println(fmt.Sprintf("time: %d:%d counter: %d", currentTime.Minute(), currentTime.Second(), counter))
	}

	counter++
}

func GetWalletBalance() []map[string]interface{} {
	data := []map[string]interface{}{}
	condition := map[string]interface{}{"is_on_hold": true}
	currency_configs := repositories.GetCurrencyNotifConfigs(&condition, nil)
	if len(*currency_configs) > 0 {
		currentTime := time.Now()
		timeInMili := currentTime.Unix() * 1000

		for _, config := range *currency_configs {
			crypto := driver.GetCrypto()
			candlesData, err := crypto.GetCandlesData(config.Symbol, 1, 0, timeInMili, "15m")
			if err != nil {
				continue
			}

			tmp := map[string]interface{}{
				"symbol":          config.Symbol,
				"balance":         config.Balance,
				"estimation_usdt": candlesData[0].Close * config.Balance,
			}

			data = append(data, tmp)
		}
	}

	return data
}

func GetBalanceFromConfig() float32 {
	var balance float32 = 0

	result := repositories.GetConfigValueByName("balance")
	if result != nil {
		resultFloat, err := strconv.ParseFloat(*result, 32)
		if err == nil {
			balance = float32(resultFloat)
		}
	}

	return balance
}

func SetBalance(balance float32) error {
	s := fmt.Sprintf("%f", balance)
	return repositories.SetConfigByName("balance", s)
}

func GetMaxHold() int64 {
	var maxHold int64 = 1

	result := repositories.GetConfigValueByName("max_hold")
	if result != nil {
		result, err := strconv.ParseInt(*result, 10, 64)
		if err == nil {
			maxHold = result
		}
	}

	return maxHold
}

func SetMaxHold(maxHold int64) error {
	s := fmt.Sprintf("%d", maxHold)
	return repositories.SetConfigByName("max_hold", s)
}

func GetMode() string {
	mode := "manual"

	result := repositories.GetConfigValueByName("mode")
	if result != nil {
		mode = *result
	}

	return mode
}

func GetOnLongIntervalWeight(coin models.BandResult, masterCoinLocal models.BandResult, startDate, endDate int64) float32 {
	responseChan := make(chan CandleResponse)

	data, err := repositories.GetCurrencyNotifConfigBySymbol(coin.Symbol)
	if err != nil {
		return 0
	}

	request := CandleRequest{
		Symbol:       data.Symbol,
		StartDate:    startDate,
		EndDate:      endDate,
		Limit:        int(CandleLimit),
		Resolution:   "1h",
		ResponseChan: responseChan,
	}

	result := MakeCryptoRequest(*data, request)
	if result == nil {
		return 0
	}

	trendChecking := true
	if masterCoinLocal.Trend == models.TREND_DOWN || (masterCoinLocal.Trend == models.TREND_SIDEWAY && masterCoinLocal.Direction == analysis.BAND_DOWN) {
		trendChecking = result.Trend == models.TREND_UP || coin.Trend == models.TREND_UP
	}
	weight := analysis.CalculateWeightLongInterval(result, masterCoinLocal.Trend)
	if analysis.IsIgnored(result) || result.Direction == analysis.BAND_DOWN || !trendChecking {
		return 0
	}

	return weight
}
