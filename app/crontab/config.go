package crontab

import (
	"telebot-trading/app/jobs"

	"github.com/robfig/cron/v3"
)

func RegisterCron(c *cron.Cron) {
	c.AddFunc("5 * * * *", jobs.CheckCryptoPrice)
}
