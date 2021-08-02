package crypto

import (
	"telebot-trading/app/models"
	"telebot-trading/app/services/crypto/driver"
	"time"
)

type CandleRequest struct {
	Symbol       string
	Start        int64
	End          int64
	Resolution   string
	ResponseChan chan CandleResponse
}

type CandleResponse struct {
	CandleData []models.CandleData
	Err        error
}

var canldeRequest chan CandleRequest
var previousTimeCheck int64 = time.Now().Unix()
var thresholdPerMinute int64 = 60
var counter int64 = 0

func DispatchRequestJob(request CandleRequest) {
	canldeRequest <- request
}

func RequestCandleService() {
	canldeRequest = make(chan CandleRequest, 100)

	for request := range canldeRequest {
		checkCounter()

		response := CandleResponse{}
		crypto := driver.GetCrypto()
		response.CandleData, response.Err = crypto.GetCandlesData(request.Symbol, request.Start, request.End, request.Resolution)

		request.ResponseChan <- response
	}

	defer func() {
		close(canldeRequest)
	}()
}

func checkCounter() {
	currentTime := time.Now().Unix()
	sleep := 0
	timeLeft := currentTime - previousTimeCheck
	if counter == thresholdPerMinute {
		if timeLeft > 0 {
			sleep = int(timeLeft) + 1
		}
		counter = 0
	}

	if timeLeft > 59 {
		counter = 0
		sleep = 1
	}

	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
		previousTimeCheck = time.Now().Unix()
	}

	counter++
}
