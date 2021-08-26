package crypto

import (
	"fmt"
	"log"
	"strconv"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto/driver"
	"telebot-trading/utils"
	"time"
)

type CandleRequest struct {
	Symbol       string
	Limit        int
	EndDate      int64
	Resolution   string
	ResponseChan chan CandleResponse
}

type CandleResponse struct {
	CandleData []models.CandleData
	Err        error
}

var canldeRequest chan CandleRequest
var previousTimeCheck time.Time = time.Now()
var thresholdPerMinute int64 = 60
var counter int64 = 0

func DispatchRequestJob(request CandleRequest) {
	canldeRequest <- request
}

func RequestCandleService() {
	canldeRequest = make(chan CandleRequest, 100)

	crypto := driver.GetCrypto()
	for request := range canldeRequest {
		checkCounter()

		response := CandleResponse{}
		response.CandleData, response.Err = crypto.GetCandlesData(request.Symbol, request.Limit, request.EndDate, request.Resolution)

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

func Buy(config models.CurrencyNotifConfig, candleData *models.CandleData) error {
	balance := GetBalance()
	if candleData == nil {
		crypto := driver.GetCrypto()
		candlesData, err := crypto.GetCandlesData(config.Symbol, 1, 0, "15")
		if err != nil {
			return err
		}
		candleData = &candlesData[0]
	}

	totalCoin := balance / candleData.Close
	repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": totalCoin, "hold_price": candleData.Close})
	SetBalance(balance - (totalCoin * candleData.Close))

	return nil
}

func Sell(config models.CurrencyNotifConfig, candleData *models.CandleData) error {
	balance := GetBalance()
	if candleData == nil {
		crypto := driver.GetCrypto()
		candlesData, err := crypto.GetCandlesData(config.Symbol, 1, 0, "15")
		if err != nil {
			return err
		}
		candleData = &candlesData[0]
	}

	totalBalance := config.Balance * candleData.Close
	repositories.UpdateCurrencyNotifConfig(config.ID, map[string]interface{}{"balance": 0})
	SetBalance(balance + totalBalance)

	return nil
}

func GetBalance() float32 {
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
