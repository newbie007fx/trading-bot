package db

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

var migrationList []*gormigrate.Migration
var migrationService *gormigrate.Gormigrate

func RegistMigration(migration *gormigrate.Migration) {
	migrationList = append(migrationList, migration)
}

func setupMigration() {
	config := gormigrate.DefaultOptions
	config.TableName = "admin_migrations"

	migrationService = gormigrate.New(GetDB(), config, migrationList)
}

func Migrate() error {
	return migrationService.Migrate()
}

func Rollback() error {
	return migrationService.RollbackLast()
}
