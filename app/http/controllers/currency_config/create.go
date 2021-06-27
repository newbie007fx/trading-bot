package currency_config

import (
	"net/http"
	"telebot-trading/app/helper"
	"telebot-trading/app/http/requests"
	"telebot-trading/app/repositories"

	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {
	return c.Render(http.StatusOK, "currency_config/create.gohtml", nil)
}

func Save(c echo.Context) (err error) {
	req := new(requests.CurrencyConfigRequest)
	c.Bind(req)

	err, data := helper.Validate(req)
	if err == nil {
		err := repositories.SaveCurrencyNotifConfig(data)
		if err == nil {
			return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusOK, nil, "Data Berhasil Disimpan"))
		}
	}

	return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
}
