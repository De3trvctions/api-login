package main

import (
	"api-login/config"
	"api-login/filter"
	initilize "api-login/initialize"
	_ "api-login/routers"
	"api-login/validation"
	"fmt"

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
	validation.Init()

	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Token", "Language"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}), web.WithReturnOnOutput(false))

	web.InsertFilter("/user/*", web.BeforeRouter, filter.LoginManager)

	// Cron Job
	initilize.InitCron()

	web.Run(fmt.Sprintf(":%d", config.HttpPort))
}
