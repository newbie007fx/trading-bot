package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"telebot-trading/app/helper"
	"telebot-trading/app/http/requests"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"
	"time"

	"telebot-trading/app/services/crypto"
	"telebot-trading/app/services/external"

	"github.com/labstack/echo/v4"
)

func ProcessTeleWebhook(c echo.Context) error {
	req := new(requests.TeleWebhookRequest)

	if err := c.Bind(req); err != nil {
		return err
	}

	responseMsg := ""

	msgData := strings.Split(strings.ToLower(req.Message.Text), " ")
	cmd := msgData[0]
	if cmd == "/register" {
		services.SaveConfig("chat_id", strconv.FormatInt(req.Message.Chat.ID, 10))
		responseMsg = "oke sukses lur"
	} else if cmd == "/hold" {
		responseMsg = "invalid format lur"
		if len(msgData) > 1 {
			balance := crypto.GetBalanceFromConfig()
			err := handlerHoldCoin(msgData[1])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = fmt.Sprintf("hold berhasil lur, saldo: %f", balance)
			}
		}
	} else if cmd == "/release" {
		responseMsg = "invalid format lur"
		if len(msgData) > 1 {
			err := handlerReleaseCoin(msgData[1])
			if err != nil {
				responseMsg = err.Error()
			} else {
				balance := crypto.GetBalanceFromConfig()
				responseMsg = fmt.Sprintf("release berhasil lur, saldo: %f", balance)
			}
		}
	} else if cmd == "/muted" {
		repositories.SetConfigByName("is-muted", strconv.FormatBool(true))
		responseMsg = "muted berhasil lur"
	} else if cmd == "/unmuted" {
		repositories.SetConfigByName("is-muted", strconv.FormatBool(false))
		responseMsg = "unmuted berhasil lur"
	} else if cmd == "/mode" {
		responseMsg = "invalid format lur"
		if len(msgData) > 1 {
			if msgData[1] == "manual" || msgData[1] == "simulation" || msgData[1] == "automatic" {
				err := repositories.SetConfigByName("mode", msgData[1])
				if err != nil {
					responseMsg = err.Error()
				} else {
					crypto.SetBalance(100)
					responseMsg = "mode berhasil diset lur"
					//jobs.ChangeStrategy()
				}
			}
		}
	} else if cmd == "/status" {
		responseMsg = "invalid format lur"
		if len(msgData) > 1 {
			msg, err := handlerStatusCoin(msgData[1])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = fmt.Sprintf("Berikut ini status untuk coin %s :\n", msgData[1])
				responseMsg += msg
			}
		}
	} else if cmd == "/status-log" {
		responseMsg = "invalid format lur"
		if len(msgData) > 2 {
			msg, err := handlerStatusLog(msgData[1], msgData[2], msgData[3])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = fmt.Sprintf("Berikut ini status untuk coin %s :\n", msgData[1])
				responseMsg += msg
			}
		}
	} else if cmd == "/balance" {
		responseMsg = crypto.GetBalances()
	} else if cmd == "/sync-balance" {
		crypto.SyncBalance()
		responseMsg = crypto.GetBalances()
	} else if cmd == "/max-hold" {
		responseMsg = "invalid format lur"
		if len(msgData) > 1 {
			err := handlerMaxHold(msgData[1])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = "max hold berhasil diubah lur"
			}
		}
	} else if cmd == "/weight-log" {
		responseMsg = "invalid format lur"
		if len(msgData) > 2 {
			msg, err := handlerWeightLog(msgData[1], msgData[2])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = msg
			}
		}
	} else {
		responseMsg = "command gak valid lur"
	}

	err := external.SendToTelegram(req.Message.Chat.ID, responseMsg)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusOK, nil, "success"))

}

func handlerHoldCoin(symbol string) error {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid symbol lur")
	}

	return crypto.HoldCoin(*currencyConfig, nil)
}

func handlerReleaseCoin(symbol string) error {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid symbol lur")
	}

	return crypto.ReleaseCoin(*currencyConfig, nil)
}

func handlerStatusCoin(symbol string) (string, error) {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("invalid symbol lur")
	}

	return crypto.GetCurrencyStatus(*currencyConfig, "15m", nil), nil
}

func handlerMaxHold(maxHold string) error {
	result, err := strconv.ParseInt(maxHold, 10, 64)
	if err != nil {
		return errors.New("invalid max hold value")
	}
	crypto.SetMaxHold(result)
	return nil
}

func handlerWeightLog(symbol, date string) (string, error) {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("invalid symbol lur")
	}
	i, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return "", errors.New("invalid log date value")
	}
	tm := time.Unix(i, 0)
	msg := crypto.GetWeightLog(*currencyConfig, tm)
	return msg, nil
}

func handlerStatusLog(symbol, date, interval string) (string, error) {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("invalid symbol lur")
	}
	i, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return "", errors.New("invalid log date value")
	}
	tm := time.Unix(i, 0)
	msg := crypto.GetCurrencyStatus(*currencyConfig, interval, &tm)
	return msg, nil
}
