package routes

import (
	"net/http"
	"telebot-trading/app/http/controllers/api"
	"telebot-trading/app/http/controllers/auth"
	"telebot-trading/app/http/controllers/home"
	"telebot-trading/app/http/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterWebRoute(e *echo.Echo) {

	e.Any("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "login")
	})

	e.GET("login", auth.ShowLoginFrom)
	e.POST("login", auth.ProcessLogin)

	g := e.Group("admin/", middleware.AuthSession)
	g.GET("dashboard", home.Dashboard)

	g.POST("logout", auth.Logout)

	a := e.Group("api/")
	a.POST("/tele-hook", api.ProcessTeleWebhook)
}
