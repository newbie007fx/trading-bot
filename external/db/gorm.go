package db

import (
	"fmt"
	"telebot-trading/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var databaseService *DatabaseService

type DatabaseService struct {
	DB *gorm.DB
}

func (DatabaseService) loadConnString() string {
	host := utils.Env("DB_HOST", "localhost")
	port := utils.Env("DB_PORT", "3306")
	user := utils.Env("DB_USER", "root")
	pass := utils.Env("DB_PASS", "root")
	database := utils.Env("DB_NAME", "todo")
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, database)
	return connstring
}

func (ds *DatabaseService) Setup() error {
	var err error
	ds.DB, err = gorm.Open(mysql.Open(ds.loadConnString()), &gorm.Config{
		Logger: logger.Default.LogMode(getLogMode()),
	})

	if err != nil {
		fmt.Println(err)
		panic("error connecting database")
	}

	setupMigration()
	return nil
}

func (ds *DatabaseService) Shutdown() {}

func GetDatabaseService() *DatabaseService {
	if databaseService == nil {
		databaseService = &DatabaseService{}
	}

	return databaseService
}

func GetDB() *gorm.DB {
	return GetDatabaseService().DB
}

func getLogMode() logger.LogLevel {
	debug := utils.Env("debug", "false")
	logMode := logger.Error
	if debug == "true" {
		logMode = logger.Info
	}
	return logMode
}
