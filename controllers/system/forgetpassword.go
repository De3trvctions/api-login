package system

import (
	"api-login/config"
	"api-login/consts"
	"api-login/models"
	"api-login/models/dto"
	"api-login/utility"
	"api-login/validation"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type ForgetPasswordController struct {
	BaseController
}

// ForgetPasswordGetCode
//
//	@Title			注册
//	@Description	注册
//	@Success		200	{object}	models.Account
//	@router			/forgetpasswordgetcode [get]
func (ctl *ForgetPasswordController) ForgetPasswordGetCode() {
	req := dto.ReqForgetPassword{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[ForgetPasswordController][ForgetPasswordGetCode] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[ForgetPasswordController][ForgetPasswordGetCode]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	acc := models.Account{}
	acc.Username = req.Username
	db := utility.NewDB()
	err := db.Get(&acc, "Username")
	if err != nil {
		logs.Error("[ForgetPasswordController][ForgetPasswordGetCode] Account not found", err)
		ctl.Error(consts.USERNAME_NOT_FOUND)
	}

	if req.Email != acc.Email {
		logs.Error("[ForgetPasswordController][ForgetPasswordGetCode] Email not match.")
		ctl.Error(consts.FORGET_PASSWORD_EMAIL_NOT_MATCH)
	}

	code, timeLeft, errCode, err := utility.SendMail(req.Email, "Forget Password Validation Code", "Please use this validation code for forget password", config.ValidCodeExpMinute)
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
