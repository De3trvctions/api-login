package routers

import (
	"api-login/controllers"
	"api-login/controllers/system"
	"api-login/controllers/user"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	userNs := web.NewNamespace("/user",
		web.NSNamespace("/info", web.NSInclude(&user.InfoController{})),
	)
	web.AddNamespace(userNs)

	systemNs := web.NewNamespace("/system",
		web.NSNamespace("/login", web.NSInclude(&system.LoginController{})),
		web.NSNamespace("/register", web.NSInclude(&system.RegisterController{})),
		web.NSNamespace("/code", web.NSInclude(&system.CodeController{})),
		web.NSNamespace("/forgetpassword", web.NSInclude(&system.ForgetPasswordController{})),
	)
	web.AddNamespace(systemNs)

	web.Router("/", &controllers.MainController{})
}
