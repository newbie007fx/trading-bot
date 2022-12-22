package commands

import (
	"telebot-trading/app/services/crypto"

	"github.com/spf13/cobra"
)

func UpdateCurrencyCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "currency:update"

	cmd.Short = "Run Telebot Trading App"

	cmd.Long = `Run Telebot Trading App`

	cmd.Run = func(cmd *cobra.Command, args []string) {
		crypto.UpdateCurrency()
	}

	return cmd
}
