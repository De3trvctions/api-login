package system

import (
	"api-login/consts"
	"api-login/models"
	"api-login/models/dto"
	"api-login/nacos"
	"api-login/redis"
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
//	@Success		200			string	"success"
//	@Param			Username	query	string	true	账号
//	@Param			Email		query	string	true	邮箱
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

	code, timeLeft, errCode, err := utility.SendMail(consts.ForgetPasswordEmailValidCode, consts.ForgetPasswordEmailValidCodeLock, req.Email, "Forget Password Validation Code", "Please use this validation code for forget password", nacos.ValidCodeExpMinute)
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

// ForgetPassword
//
//	@Title			注册
//	@Description	注册
//	@Success		200	string	"success"
//	@router			/forgetpassword [get]
func (ctl *ForgetPasswordController) ForgetPassword() {
	req := dto.ReqForgetPasswordSetNew{}

	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[ForgetPasswordController][ForgetPassword] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[ForgetPasswordController][ForgetPassword]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	db := utility.NewDB()

	acc := models.Account{}
	acc.Username = req.Username
	acc.Email = req.Email
	err := db.Get(&acc, "Username")
	if err != nil {
		logs.Error("[ForgetPasswordController][ForgetPassword] Account not found", err)
		ctl.Error(consts.USERNAME_NOT_FOUND)
	}

	if req.Email != acc.Email {
		logs.Error("[ForgetPasswordController][ForgetPassword] Email not match.")
		ctl.Error(consts.FORGET_PASSWORD_EMAIL_NOT_MATCH)
	}

	// Check Valid Code
	ex, _ := redis.Exists(fmt.Sprintf(consts.ForgetPasswordEmailValidCode, req.Email))
	if ex {
		defer utility.DelEmailValidCodeLock(req.Email)
		validCode, _ := redis.Get(fmt.Sprintf(consts.ForgetPasswordEmailValidCode, req.Email))
		if validCode != req.ValidCode {
			logs.Error("[ForgetPasswordController][ForgetPassword] Valid Code not match. Please try again")
			ctl.Error(consts.VALID_CODE_NOT_MATCH)
		}
	} else {
		ctl.Error(consts.VALID_CODE_NOT_MATCH)
	}

	acc.SetHashPassword(req.Password)
	acc.SetUpdateTime()

	_, err = db.Update(&acc, "Password", "UpdateTime")
	if err != nil {
		logs.Error("[ForgetPasswordController][ForgetPassword] Update account error. %+v , error: %+v", acc, err)
		ctl.Error(consts.DB_UPDATE_FAILED)
	}

	ctl.Success("success")
}
