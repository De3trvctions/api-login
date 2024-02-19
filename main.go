package main

import (
	"api-login/config"
	_ "api-login/routers"
	"fmt"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	// Init required information
	config.InitLogs()
	config.InitRedis()
	config.InitLanguage()

	// redisKey := "mykey"
	// _ = redis.Set(redisKey, "myvalue", 0)
	// aaa, _ := redis.Get(redisKey)
	// logs.Error(aaa)

	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Token", "token", "RandCloud", "HelpLink", "Redirect", "Language"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}), web.WithReturnOnOutput(false))
	web.Run(fmt.Sprintf(":%d", config.HttpPort))
}
