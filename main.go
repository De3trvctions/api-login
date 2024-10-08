package main

import (
	cron "api-login/crontask"
	"api-login/filter"
	_ "api-login/routers"
	"fmt"
	"standard-library/config"
	initilize "standard-library/initialize"
	"standard-library/validation"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	// Init required information
	initilize.InitLogs()
	initilize.InitNacosConfig()
	initilize.InitRedis()
	initilize.InitDB()
	initilize.InitLanguage()
	initilize.InitMail()
	initilize.InitGRPC()
	validation.Init()

	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Token", "Language"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}), web.WithReturnOnOutput(true))

	web.InsertFilter("/user/*", web.BeforeRouter, filter.LoginManager)

	// Cron Job
	cron.InitCronTask()

	web.Run(fmt.Sprintf(":%d", config.HttpPort))
}
