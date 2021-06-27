package httpserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"telebot-trading/routes"
	"telebot-trading/utils"

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

func httpErrorHanlder() func(error, echo.Context) {
	return func(err error, c echo.Context) {
		log.Println(err.Error())

		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		resourcePath := utils.Env("RESOURCES_PATH", "resources")
		errorPage := fmt.Sprintf("%s/views/errors/%d.html", resourcePath, code)
		content, _ := ioutil.ReadFile(errorPage)
		c.HTML(code, string(content))
	}
}
