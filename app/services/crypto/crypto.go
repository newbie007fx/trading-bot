package crypto

import (
	"fmt"
	"log"
	"telebot-trading/app/models"
	"telebot-trading/app/services/crypto/driver"
	"telebot-trading/utils"
	"time"
)

type CandleRequest struct {
	Symbol       string
	Limit        int
	Resolution   string
	ResponseChan chan CandleResponse
}

type CandleResponse struct {
	CandleData []models.CandleData
	Err        error
}

var canldeRequest chan CandleRequest
var previousTimeCheck time.Time = time.Now()
var thresholdPerMinute int64 = 58
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
		response.CandleData, response.Err = crypto.GetCandlesData(request.Symbol, request.Limit, request.Resolution)

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
