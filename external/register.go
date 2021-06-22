package external

import (
	"telebot-trading/bootstrap"
	"telebot-trading/external/cli"
	"telebot-trading/external/configuration"
	"telebot-trading/external/db"
	"telebot-trading/external/httpserver"
	"telebot-trading/external/scheduler"
)

func init() {
	bootstrap.SetMainService(&cli.ConsoleService{})
	bootstrap.RegisterService(&configuration.ConfigService{})
	bootstrap.RegisterService(httpserver.GetRouteService())
	bootstrap.RegisterService(db.GetDatabaseService())
	bootstrap.RegisterService(scheduler.GetCronService())
}
