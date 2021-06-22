package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"telebot-trading/app/helper"
	"telebot-trading/app/http/requests"
	"telebot-trading/app/services"

	"github.com/labstack/echo/v4"
)

func ProcessTeleWebhook(c echo.Context) error {
	req := new(requests.TeleWebhookRequest)

	if err := c.Bind(req); err != nil {
		return err
	}

	if strings.ToLower(req.Message.Text) == "/register" {
		services.SaveConfig("chat_id", strconv.FormatInt(req.Message.Chat.ID, 10))
		err := services.SendToTelegram(req.Message.Chat.ID, "oke sukses lur")
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := services.SendToTelegram(req.Message.Chat.ID, "command gak valid lur")
		if err != nil {
			log.Panic(err)
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusOK, nil, "success"))

}
