package system

import (
	"api-login/config"
	"api-login/consts"
	"api-login/utility"
	"api-login/validation"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type CodeController struct {
	BaseController
}

// GetValidCode
//
//	@Title			忘记密码获取邮件验证码
//	@Description	忘记密码获取邮件验证码
//	@Success		200		{string}	"success"
//	@Param			Email	formData	string	true	Email
//	@router			/getvalidcode [get]
func (ctl *CodeController) GetValidCode() {
	type reqGetCode struct {
		Email string `valid:"Required;Email;"`
	}
	req := reqGetCode{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[CodeController][GetCode] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[CodeController][GetCode]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	code, timeLeft, errCode, err := utility.SendMail(consts.RegisterEmailValidCode, consts.RegisterEmailValidCodeLock, req.Email, "Email Validation Code", "", config.ValidCodeExpMinute)
	if err != nil || errCode != 0 || timeLeft > 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		if err != nil {
			logs.Error("[CodeController][GetCode]Db error:", err)
		}
		if timeLeft > 0 {
			a := time.Unix(timeLeft, 0)
			b := time.Now()
			c := int(a.Sub(b).Seconds())
			ctl.Error(errCode, fmt.Sprintf("Please try again in %v seconds", c))
		} else {
			ctl.Error(errCode)
		}
	}

	if web.BConfig.RunMode == web.DEV {
		ctl.Success(web.M{
			"Code": code,
		})
	} else {
		ctl.Success("success")
	}
}
