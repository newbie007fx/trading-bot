package commands

import (
	"log"
	"strconv"
	"telebot-trading/app/repositories"
	"telebot-trading/app/services/crypto"
	"time"

	"github.com/spf13/cobra"
)

func TestingWeightCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "testing:weight"

	cmd.Short = "Testing weight log calculation"

	cmd.Long = `Testing weight log calculation`

	cmd.Flags().StringP("symbol", "s", "", "Set the coin symbol")

	cmd.Flags().StringP("time", "t", "", "Set the epoch time")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		go crypto.RequestCandleService()
		go crypto.StartSyncBalanceService()
		symbol, _ := cmd.Flags().GetString("symbol")
		date, _ := cmd.Flags().GetString("time")
		currencyConfig, err := repositories.GetCurrencyNotifConfigBySymbol(symbol)
		if err != nil {
			log.Println(err.Error())
			return
		}

		i, err := strconv.ParseInt(date, 10, 64)
		if err != nil {
			log.Println("invalid log date value")
			return
		}
		tm := time.Unix(i, 0)
		msg := crypto.GetWeightLog(*currencyConfig, tm)
		log.Println(msg)
	}

	return cmd
}
