package commands

import (
	"log"
	"telebot-trading/external/db"

	"github.com/spf13/cobra"
)

func MigrateCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "migrate"

	cmd.Short = "Run Database Migration"

	cmd.Long = `Run Database Migration`

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := db.Migrate(); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
		log.Printf("Migration did run successfully")
	}

	return cmd
}
