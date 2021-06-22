package commands

import (
	"log"
	"telebot-trading/app/models"
	"telebot-trading/external/db"
	"time"

	"github.com/spf13/cobra"
)

func AddNotifConfigCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "add:config"

	cmd.Short = "Run Database Migration"

	cmd.Long = `Run Database Migration`

	cmd.Flags().StringP("symbol", "s", "", "Set Symbol")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		symbol, _ := cmd.Flags().GetString("symbol")
		notifConfig := getConfigData(symbol)
		result := db.GetDB().Create(&notifConfig)
		if result.Error != nil {
			log.Print("Error generate symbol data")
			log.Print(result.Error.Error())
		}

		log.Printf("add Symbol data did run successfully")
	}

	return cmd
}

func getConfigData(symbol string) models.CurrencyNotifConfig {
	notifConfig := models.CurrencyNotifConfig{}
	notifConfig.Symbol = symbol
	notifConfig.CreatedAt = time.Now()
	notifConfig.UpdatedAt = time.Now()
	return notifConfig
}
