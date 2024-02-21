package system

import (
	"api-login/consts"
	"api-login/models"
	"api-login/models/dto"
	"api-login/validation"

	"github.com/beego/beego/v2/core/logs"
)

type RegisterController struct {
	BaseController
}

// Register
//
//	@Title			注册
//	@Description	注册
//	@Success		200	{object}	models.Account
//	@router			/ [post]
func (ctl *RegisterController) Register() {
	req := dto.ReqRegister{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[RegisterController][Register] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[RegisterController][Register]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	acc := models.Account{}
	_, _, _ = acc.List(acc)

	errCode, err := acc.Register(req)
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[RegisterController][Register]Db error:", err)
		ctl.Error(errCode)
	}

	ctl.Success("success")
	//ctl.Success(web.M{"Items": data, "Pagination": pagination})
}
