package user

import (
	account "api-login-proto/account"
	"api-login/controllers/system"
	"api-login/models/dto"
	"api-login/repositories"
	"api-login/services"
	"standard-library/consts"
	"standard-library/validation"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type InfoController struct {
	system.PermissionController
	GRPCClient account.UserAccountServiceClient
}

func (ctl *InfoController) Prepare() {
	ctl.PermissionController.Prepare()
	ctl.ConnGRpc("service-account")
	ctl.GRPCClient = account.NewUserAccountServiceClient(ctl.GrpcConn.Conn())
	if ctl.GRPCClient == nil {
		ctl.Error(consts.SERVER_ERROR)
	}
}

// Detail
//
//	@Title			注册
//	@Description	注册
//	@Success		200			{object}	web.M
//	@Param			AccountId	query		string	true	登陆名
//	@router			/detail [get]
func (ctl *InfoController) Detail() {
	defer logs.Info("[InfoController][Detail]: Enter URL %+v", ctl.RequestUrl)
	req := dto.ReqAccountDetail{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[Login] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[RegisterController][Register]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	req.AccountId = ctl.AccountId
	repo := repositories.NewAccountRepository()
	accountInfo, err := repo.Detail(req)
	if err != nil {
		logs.Error("[InfoController][Detail] query account detail error:", err)
		ctl.Error(consts.DB_GET_FAILED)
		return
	}

	ctl.Success(web.M{"Info": accountInfo})
}

// Detail
//
//	@Title			编辑
//	@Description	编辑
//	@Success		200			{object}	web.M
//	@Param			AccountId	formData	string	true	登陆ID
//	@Param			Email		formData	string	true	登陆名
//	@Param			Password	formData	string	true	登陆名
//	@Param			NewPassword	formData	string	true	登陆名
//	@router			/ [put]
func (ctl *InfoController) Edit() {
	req := dto.ReqEditAccount{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[InfoController][Edit] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	req.AccountId = ctl.AccountId
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[InfoController][Edit]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	service := services.NewAccountService()
	errCode, err := service.EditProfile(ctl.AccountId, req)
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[InfoController][Edit]Db error:", err)
		ctl.Error(errCode)
	}

	ctl.Success("success")
}

// Info
//
//	@Title			Account Info
//	@Description	Account Info Detail
//	@Success		200			{object}	web.M
//	@Param			AccountId	query		int64	false	"AccountID"
//	@Param			Username	query		string	false	"Account Username"
//	@Param			Email		query		string	false	"Account Email"
//	@Param			CreateTime	query		int64	false	"Account create time"
//	@Param			Page		query		int64	false	"Page"
//	@Param			PageSize	query		int64	false	"Page Size"
//	@router			/info [get]
func (ctl *InfoController) Info() {
	req := dto.ReqAccountDetail{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[InfoController][Info] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[InfoController][Info] FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	ctl.CommonJSONRequest(req, ctl.GRPCClient.Info)
}
