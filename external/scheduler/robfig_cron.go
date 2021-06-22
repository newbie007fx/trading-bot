package scheduler

import (
	"telebot-trading/app/crontab"
	"time"

	"github.com/robfig/cron/v3"
)

var cronService *CronService

type CronService struct {
	server *cron.Cron
}

func (rc *CronService) Setup() error {
	rc.server = cron.New()

	timeLocation, _ := time.LoadLocation("Asia/Jakarta")
	cron.WithLocation(timeLocation)

	crontab.RegisterCron(rc.server)

	return nil
}

func (rc *CronService) Start() {
	rc.server.Start()
}

func (rc *CronService) Shutdown() {
	rc.server.Stop()
}

func GetCronService() *CronService {
	if cronService == nil {
		cronService = new(CronService)
	}

	return cronService
}
