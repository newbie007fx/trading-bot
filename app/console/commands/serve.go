package commands

import (
	"telebot-trading/external/httpserver"
	"telebot-trading/utils"
	"strconv"

	"github.com/spf13/cobra"
)

func ServeCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "serve"

	cmd.Short = "Run Telebot Trading App"

	cmd.Long = `Run Telebot Trading App`

	port := utils.Env("PORT", "8080")

	deafultPort, _ := strconv.Atoi(port)

	cmd.Flags().IntP("port", "p", deafultPort, "Set the port")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		httpserver.Run(port)
	}

	return cmd
}
