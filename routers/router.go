package routers

import (
	"api-login/controllers"
	"api-login/controllers/system"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	systemNs := web.NewNamespace("/system",
		/*********************系统路由********************************/
		// 登录
		web.NSNamespace("/user", web.NSInclude(&system.LoginController{})),
	)
	web.AddNamespace(systemNs)

	web.Router("/", &controllers.MainController{})
}
