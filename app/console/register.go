package console

import (
	"telebot-trading/app/console/commands"
	"telebot-trading/external/cli"
)

func init() {
	cli.RegisterCommand(commands.ServeCommand())
	cli.RegisterCommand(commands.MigrateCommand())
	cli.RegisterCommand(commands.RollbackCommand())
	cli.RegisterCommand(commands.GenerateAdminCommand())
	cli.RegisterCommand(commands.AddNotifConfigCommand())
	cli.RegisterCommand(commands.CronRunCommand())
	cli.RegisterCommand(commands.WorkerRunCommand())
}
