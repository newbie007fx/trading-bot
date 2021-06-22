package auth

import (
	"net/http"
	"telebot-trading/app/helper"

	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) (err error) {
	sess := helper.GetSession(c)
	sess.Destroy()
	sess.AddFlashMessage("message", "anda telah logout")
	return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusOK, nil, "Logout Berhasil"))
}
