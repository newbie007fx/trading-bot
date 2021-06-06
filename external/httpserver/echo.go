package httpserver

import (
	"fmt"
	"telebot-trading/utils"

	"github.com/labstack/echo/v4"
)

var routeService *RouterService

type RouterService struct {
	server *echo.Echo
}

func (rs *RouterService) Setup() error {
	rs.server = echo.New()

	debug := utils.Env("debug", "false")

	rs.server.Debug = (debug == "true")

	registerRouting(rs.server)

	return nil
}

func (rs *RouterService) Start(port int) {
	address := fmt.Sprintf(":%d", port)
	rs.server.Logger.Fatal(rs.server.Start(address))
}

func (RouterService) Shutdown() {}

func GetRouteService() *RouterService {
	if routeService == nil {
		routeService = &RouterService{}
	}

	return routeService
}

func Run(port int) {
	service := GetRouteService()
	service.Start(port)
}
