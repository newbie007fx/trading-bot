package auth

import (
	"net/http"
	"telebot-trading/app/helper"
	"telebot-trading/app/http/requests"
	"telebot-trading/app/services"

	"github.com/labstack/echo/v4"
)

func ShowLoginFrom(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/login.gohtml", nil)
}

func ProcessLogin(c echo.Context) (err error) {
	req := new(requests.LoginRequest)
	c.Bind(req)

	sess := helper.GetSession(c)

	err, data := helper.Validate(req)
	if err == nil {
		err, admin := services.Login(data["email"].(string), data["password"].(string))
		if err == nil {
			data := map[string]interface{}{}
			data["is_login"] = true
			data["id"] = admin.ID
			sess.Set(data)
			return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusAccepted, nil, "Login Berhasil"))
		}
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse(http.StatusUnauthorized, nil, "Login gagal, pastikan username dan password valid"))
	}

	return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
}
