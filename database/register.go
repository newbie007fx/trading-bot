package database

import (
	"telebot-trading/database/migrations"
	"telebot-trading/external/db"
)

func init() {
	db.RegistMigration(migrations.CreateAdminTable())
	db.RegistMigration(migrations.CreateCurrencyNotifConfigTable())
	db.RegistMigration(migrations.CreateSiteConfigTable())
	db.RegistMigration(migrations.AddOnHoldColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddVolumeColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddBalanceColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddHoldPriceColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddHoldedAtColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddReachTargetProfitAtColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddPriceChangeColumnCurrencyNotifTable())
	db.RegistMigration(migrations.AddStatusAndConfigColumnCurrencyNotifTable())
	db.RegistMigration(migrations.DropCreatedAtAndUpdatedAtOnCurrentyConfig())
	db.RegistMigration(migrations.AddCreatedAtAndUpdatedAtOnCurrentyConfig())
}
