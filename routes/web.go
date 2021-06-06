package routes

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func RegisterWebRoute(e *echo.Echo) {

	e.Any("/", func(c echo.Context) error {
		return errors.New("invalid request")
	})
}
