package commands

import (
	"log"
	"telebot-trading/external/db"

	"github.com/spf13/cobra"
)

func RollbackCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "rollback"

	cmd.Short = "Rollback Last Database Migration"

	cmd.Long = `Rollback Last Database Migration`

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := db.Rollback(); err != nil {
			log.Fatalf("Could not rollback: %v", err)
		}
		log.Printf("Rollback did run successfully")
	}

	return cmd
}
