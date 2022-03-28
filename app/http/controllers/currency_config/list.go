package currency_config

import (
	"net/http"
	"telebot-trading/app/repositories"

	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	orderBy := "is_master desc, is_on_hold desc, price_changes desc"
	data := repositories.GetCurrencyNotifConfigs(nil, nil, &orderBy)
	return c.Render(http.StatusOK, "currency_config/list.gohtml", map[string]interface{}{
		"currency_configs": data,
	})
}
