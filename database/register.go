package database

import (
	"telebot-trading/database/migrations"
	"telebot-trading/external/db"
)

func init() {
	db.RegistMigration(migrations.CreateSiteConfigTable())
}
