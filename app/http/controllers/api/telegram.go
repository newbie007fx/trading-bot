package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"telebot-trading/app/helper"
	"telebot-trading/app/http/requests"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services"

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
		if len(msgData) > 2 {
			err := handlerHoldCoin(msgData[1])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = "hold berhasil lur"
			}
		}
	} else if cmd == "/release" {
		responseMsg = "invalid format lur"
		if len(msgData) > 2 {
			err := handlerReleaseCoin(msgData[1])
			if err != nil {
				responseMsg = err.Error()
			} else {
				responseMsg = "release berhasil lur"
			}
		}
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

	if !currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": true,
		}
		err = repositories.UpdateCurrencyNotifConfig(currencyConfig.ID, data)
		if err != nil {
			log.Println(err.Error())
			return errors.New("error waktu update lur")
		}
	}

	return nil
}

func handlerReleaseCoin(symbol string) error {
	currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid symbol lur")
	}

	if currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": false,
		}
		err = repositories.UpdateCurrencyNotifConfig(currencyConfig.ID, data)
		if err != nil {
			log.Println(err.Error())
			return errors.New("error waktu update lur")
		}
	}

	return nil
}
