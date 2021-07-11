package currency_config

import (
	"net/http"
	"telebot-trading/app/repositories"

	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	data := repositories.GetCurrencyNotifConfigs(nil)
	return c.Render(http.StatusOK, "currency_config/list.gohtml", map[string]interface{}{
		"currency_configs": data,
	})
}
