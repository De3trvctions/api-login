package system

import (
	"api-login/consts"
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
	ctl.Error(consts.ERROR)
}
