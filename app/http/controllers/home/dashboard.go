package home

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Dashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard.gohtml", nil)
}
