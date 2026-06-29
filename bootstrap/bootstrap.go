package bootstrap

import (
	"api-login/config"
	"api-login/filter"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func InitApp() error {
	return config.Load()
}

func ConfigureWeb() {
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Token", "Language"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: false,
	}), web.WithReturnOnOutput(true))

	web.InsertFilter("/user/*", web.BeforeRouter, filter.LoginManager, web.WithReturnOnOutput(true))
	web.BConfig.CopyRequestBody = true
}
