package httpserver

import (
	"telebot-trading/routes"

	"github.com/labstack/echo/v4"
)

func registerRouting(e *echo.Echo) {

	routes.RegisterWebRoute(e)
}
