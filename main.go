package main

import (
	"api-login/bootstrap"
	"api-login/config"
	cron "api-login/crontask"
	_ "api-login/routers"
	"fmt"
	initilize "standard-library/initialize"
	"standard-library/validation"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	if err := bootstrap.InitApp(); err != nil {
		panic(err)
	}

	// Init required information
	initilize.InitLogs()
	initilize.InitES()
	initilize.InitNacosConfig()
	initilize.InitRedis()
	initilize.InitDB()
	initilize.InitLanguage()
	initilize.InitMail()
	initilize.InitGRPC()
	validation.Init()

	bootstrap.ConfigureWeb()

	// Cron Job
	cron.InitCronTask()

	web.Run(fmt.Sprintf(":%d", config.HttpPort))
}
