package crontab

import (
	"github.com/robfig/cron/v3"
)

func RegisterCron(c *cron.Cron) {
	// c.AddFunc("*/5 * * * *", jobs.CheckCryptoPrice)
	// c.AddFunc("01 08,22 * * *", jobs.UpdateVolume)
}
