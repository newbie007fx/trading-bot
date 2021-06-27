package currency_config

import (
	"net/http"
	"strconv"
	"telebot-trading/app/helper"
	"telebot-trading/app/http/requests"
	"telebot-trading/app/repositories"

	"github.com/labstack/echo/v4"
)

func Edit(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currencyConfig, err := repositories.GetCurrencyNotifConfig(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Halaman tidak ditemukan")
	}

	return c.Render(http.StatusOK, "currency_config/edit.gohtml", map[string]interface{}{
		"currencyConfig": currencyConfig,
	})
}

func Update(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	req := new(requests.CurrencyConfigRequest)
	c.Bind(req)

	err, data := helper.Validate(req)
	if err == nil {
		err = repositories.UpdateCurrencyNotifConfig(uint(id), data)
		if err == nil {
			return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusCreated, nil, "Data berhasil diupdate"))
		}
	}

	return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
}

func Hold(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currencyConfig, err := repositories.GetCurrencyNotifConfig(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(http.StatusUnprocessableEntity, nil, "Data tidak ditemukan"))
	}

	if !currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": true,
		}
		err = repositories.UpdateCurrencyNotifConfig(uint(id), data)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
		}
	}
	return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusCreated, nil, "Data berhasil diupdate"))

}

func Release(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currencyConfig, err := repositories.GetCurrencyNotifConfig(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(http.StatusUnprocessableEntity, nil, "Data tidak ditemukan"))
	}

	if currencyConfig.IsOnHold {
		data := map[string]interface{}{
			"is_on_hold": false,
		}
		err := repositories.UpdateCurrencyNotifConfig(uint(id), data)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
		}
	}
	return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusCreated, nil, "Data berhasil diupdate"))

}

func SetMaster(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currencyConfig, err := repositories.GetCurrencyNotifConfig(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(http.StatusUnprocessableEntity, nil, "Data tidak ditemukan"))
	}

	if !currencyConfig.IsMaster {
		err := repositories.SetMaster(uint(id))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, helper.ErrorResponse(http.StatusUnprocessableEntity, err, ""))
		}
	}
	return c.JSON(http.StatusOK, helper.SuccessResponse(http.StatusCreated, nil, "Data berhasil diupdate"))
}
