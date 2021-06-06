package routes

import (
	"errors"
	"telebot-trading/app/http/controllers/telegram_hook"

	"github.com/labstack/echo/v4"
)

func RegisterWebRoute(e *echo.Echo) {

	e.Any("/", func(c echo.Context) error {
		return errors.New("invalid request")
	})

	e.POST("/tele-hook", telegram_hook.ProcessTeleWebhook)
}
