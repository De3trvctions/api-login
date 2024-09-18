package system

import (
	"api-login/models"
	"api-login/models/dto"
	"fmt"
	"standard-library/consts"
	"standard-library/redis"
	"standard-library/utility"
	"standard-library/validation"

	"github.com/beego/beego/v2/core/logs"
)

type RegisterController struct {
	BaseController
}

// Register
//
//	@Title			注册
//	@Description	注册
//	@Param			Username	formData	string	true	账号
//	@Param			Password	formData	string	true	密码
//	@Param			Email		formData	string	true	邮箱
//	@Param			ValidCode	formData	string	true	验证码
//	@Success		200			string		"success"
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

	// Check Valid Code
	ex, _ := redis.Exists(fmt.Sprintf(consts.RegisterEmailValidCode, req.Email))
	if ex {
		defer utility.DelEmailValidCodeLock(req.Email)
		validCode, _ := redis.Get(fmt.Sprintf(consts.RegisterEmailValidCode, req.Email))
		if validCode != req.ValidCode {
			logs.Error("[RegisterController][Register] Valid Code not match. Please try again")
			ctl.Error(consts.VALID_CODE_NOT_MATCH)
		}
	} else {
		ctl.Error(consts.VALID_CODE_NOT_MATCH)
	}

	acc := models.Account{}
	errCode, err := acc.Register(req)
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[RegisterController][Register]Db error:", err)
		ctl.Error(errCode)
	}

	ctl.Success("success")
}
