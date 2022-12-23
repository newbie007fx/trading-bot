package currency_config

import (
	"net/http"
	"strconv"
	"telebot-trading/app/repositories"

	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	orderBy := "is_on_hold desc, price_changes desc"
	limit := 50
	pageString := c.QueryParam("page")
	page := 0
	offset := 0
	if pageString != "" {
		tmp, err := strconv.ParseInt(pageString, 10, 0)
		if err == nil {
			page = int(tmp)
			offset = page * limit
		}
	}
	data := repositories.GetCurrencyNotifConfigs(nil, &limit, &offset, &orderBy)
	return c.Render(http.StatusOK, "currency_config/list.gohtml", map[string]interface{}{
		"currency_configs": data,
		"page":             page,
	})
}
