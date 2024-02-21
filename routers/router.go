package routers

import (
	"api-login/controllers"
	"api-login/controllers/system"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	systemNs := web.NewNamespace("/system",
		web.NSNamespace("/user", web.NSInclude(&system.LoginController{})),
		web.NSNamespace("/register", web.NSInclude(&system.RegisterController{})),
	)
	web.AddNamespace(systemNs)

	web.Router("/", &controllers.MainController{})
}
