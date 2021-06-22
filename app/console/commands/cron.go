package commands

import (
	"telebot-trading/external/scheduler"

	"github.com/spf13/cobra"
)

func CronRunCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "cron:run"

	cmd.Short = "Run Cron Job"

	cmd.Long = `Run Cron Job`

	cmd.Run = func(cmd *cobra.Command, args []string) {
		service := scheduler.GetCronService()
		service.Start()
	}

	return cmd
}
