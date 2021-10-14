package commands

import (
	"telebot-trading/app/jobs"

	"github.com/spf13/cobra"
)

func WorkerRunCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "worker:run"

	cmd.Short = "Run Worker"

	cmd.Long = `Run Worker`

	cmd.Run = func(cmd *cobra.Command, args []string) {
		coba := make(chan bool)
		go jobs.StartCryptoWorker()
		<-coba
	}

	return cmd
}
