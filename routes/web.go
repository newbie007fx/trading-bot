package routes

import (
	"net/http"
	"telebot-trading/app/http/controllers/api"
	"telebot-trading/app/http/controllers/auth"
	"telebot-trading/app/http/controllers/currency_config"
	"telebot-trading/app/http/controllers/home"
	"telebot-trading/app/http/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterWebRoute(e *echo.Group) {

	e.Any("", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "login")
	})

	e.GET("login", auth.ShowLoginFrom)
	e.POST("login", auth.ProcessLogin)

	g := e.Group("admin/", middleware.AuthSession)
	g.GET("dashboard", home.Dashboard)

	g.GET("currency-config", currency_config.List)
	g.GET("currency-config/create", currency_config.Create)
	g.POST("currency-config/create", currency_config.Save)
	g.GET("currency-config/edit/:id", currency_config.Edit)
	g.PUT("currency-config/edit/:id", currency_config.Update)
	g.PUT("currency-config/hold/:id", currency_config.Hold)
	g.PUT("currency-config/release/:id", currency_config.Release)
	g.PUT("currency-config/set-master/:id", currency_config.SetMaster)
	g.DELETE("currency-config/delete/:id", currency_config.Delete)

	g.POST("logout", auth.Logout)
}

func RegisterApiRoute(e *echo.Group) {
	e.POST("tele-hook", api.ProcessTeleWebhook)
}
