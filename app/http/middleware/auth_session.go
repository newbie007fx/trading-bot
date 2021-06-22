package middleware

import (
	"net/http"
	"telebot-trading/app/helper"

	"github.com/labstack/echo/v4"
)

func AuthSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := helper.GetSession(c)
		isLogin := sess.Get("is_login")
		if isLogin == nil {
			sess.AddFlashMessage("message", "Anda belum login, silahkan login terlebih dahulu")
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}
		return next(c)
	}
}
