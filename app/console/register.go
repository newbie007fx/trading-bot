package console

import (
	"telebot-trading/app/console/commands"
	"telebot-trading/external/cli"
)

func init() {
	cli.RegisterCommand(commands.ServeCommand())
	cli.RegisterCommand(commands.MigrateCommand())
	cli.RegisterCommand(commands.RollbackCommand())
}
