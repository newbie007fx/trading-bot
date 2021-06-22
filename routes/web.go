package routes

import (
	"net/http"
	"telebot-trading/app/http/controllers/api"
	"telebot-trading/app/http/controllers/auth"
	"telebot-trading/app/http/controllers/home"
	localMiddleware "telebot-trading/app/http/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterWebRoute(e *echo.Echo) {

	e.Any("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "login")
	})

	e.GET("login", auth.ShowLoginFrom)
	e.POST("login", auth.ProcessLogin)

	g := e.Group("admin/", localMiddleware.AuthSession)
	g.Use(session.Middleware(sessions.NewCookieStore([]byte("597598ca-f457-437f-b08e-5823ff94e0aa"))))
	g.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-CSRF-TOKEN",
	}))
	g.GET("dashboard", home.Dashboard)

	g.POST("logout", auth.Logout)

	a := e.Group("api/")
	a.POST("tele-hook", api.ProcessTeleWebhook)
}
