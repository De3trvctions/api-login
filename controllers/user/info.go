package user

import (
	"api-login/consts"
	"api-login/controllers/system"
	"api-login/models"
	"api-login/models/dto"
	"api-login/utility"
	"api-login/validation"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type InfoController struct {
	system.BaseController
}

func (ctl *InfoController) Prepare() {
	ctl.BaseController.Prepare()
}

// Detail
//
//	@Title			注册
//	@Description	注册
//	@Success		200	{object}	web.M
//	@Param			AccountId	formData	string	true	登陆名
//	@router			/detail [get]
func (ctl *InfoController) Detail() {
	req := dto.ReqAccountDetail{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[Login] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[RegisterController][Register]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	acc := models.Account{}
	acc.Id = req.AccountId
	db := utility.NewDB()

	tx, err := db.Begin()
	if err != nil {
		logs.Error(err)
	}
	defer tx.Commit()

	err = db.Get(&acc, "Id")
	if err != nil {
		logs.Error("[InfoController][Detail]")
	}

	acc.SetUpdateTime()
	_, err = tx.Update(&acc, "UpdateTime")
	if err != nil {
		logs.Error("[InfoController][Detail]Update Error:", err)
	}

	ctl.Success(web.M{"Info": acc})
}
