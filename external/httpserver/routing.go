package httpserver

import (
	"telebot-trading/routes"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerRouting(e *echo.Echo) {
	g := e.Group("/")
	g.Use(session.Middleware(sessions.NewCookieStore([]byte("597598ca-f457-437f-b08e-5823ff94e0aa"))))
	g.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-CSRF-TOKEN",
	}))
	routes.RegisterWebRoute(g)

	a := e.Group("api/")
	routes.RegisterApiRoute(a)
}
