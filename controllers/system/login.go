package system

import (
	"api-login/consts"
	"api-login/models/dto"

	"github.com/beego/beego/v2/core/logs"
)

type LoginController struct {
	BaseController
}

// Login
//
//	@Title			登陆
//	@Description	登陆
//	@Success		200	"success"
//	@router			/login [get]
func (ctl *LoginController) Login() {
	req := dto.ReqLogin{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[Login] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}

	ctl.Error(consts.SERVER_ERROR)
}
