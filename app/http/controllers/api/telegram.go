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

	"telebot-trading/app/services/crypto"

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
			balance := crypto.GetBalance()
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
				balance := crypto.GetBalance()
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
	} else if cmd == "/balance" {
		responseMsg = services.GetBalance()
	} else if cmd == "/sync-balance" {
		crypto.SyncBalance()
		responseMsg = services.GetBalance()
	} else {
		responseMsg = "command gak valid lur"
	}

	err := services.SendToTelegram(req.Message.Chat.ID, responseMsg)
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

	return services.HoldCoin(*currencyConfig, nil)
}

func handlerReleaseCoin(symbol string) error {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid symbol lur")
	}

	return services.ReleaseCoin(*currencyConfig, nil)
}

func handlerStatusCoin(symbol string) (string, error) {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("invalid symbol lur")
	}

	return services.GetCurrencyStatus(*currencyConfig), nil
}
