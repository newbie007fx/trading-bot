package currency_config

import (
	"net/http"
	"strconv"
	"telebot-trading/app/helper"
	"telebot-trading/app/repositories"

	"github.com/labstack/echo/v4"
)

func Delete(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := repositories.DeleteCurrencyNotifConfig(uint(id))
	if err == nil {
		return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusOK, nil, "Data berhasil dihapus"))
	}

	return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
}
