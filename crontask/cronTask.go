package crontask

import (
	"api-login/dailyreport"

	"github.com/robfig/cron/v3"
)

func InitCronTask() {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("1 * * * * *", dailyreport.RunDailyLoginLogReport) // RUn at 0015 everyday

	c.Start()
}
