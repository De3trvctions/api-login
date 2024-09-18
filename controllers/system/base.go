package system

import (
	"fmt"
	"standard-library/consts"
	"standard-library/jwt"
	"standard-library/nacos"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

type BaseController struct {
	web.Controller
	i18n.Locale
	Code       int64
	Result     interface{}
	Ip         string
	DeviceId   string
	Msg        string
	Token      string
	AccountId  int64
	Username   string
	RequestUrl string
}

func (ctl *BaseController) Prepare() {

}

func (ctl *BaseController) Success(obj interface{}) {
	ctl.Code = consts.SUCCESS_REQUEST
	ctl.Result = obj
	ctl.TraceJson(false)
}

func (ctl *BaseController) Error(code int64, msg ...string) {
	if code < 1 {
		code = consts.SERVER_ERROR
	}
	ctl.Code = code

	ctl.GetLanguage()

	if ctl.Code > 0 {
		ctl.Msg = ctl.Tr("error." + fmt.Sprintf("%d", ctl.Code))
	}

	if len(msg) > 0 {
		ctl.Msg += "(" + fmt.Sprint(msg) + ")"
	}

	ctl.TraceJson(true)
}

func (ctl *BaseController) TraceJson(errorEnd bool) {
	res := map[string]interface{}{"Code": ctl.Code, "Msg": ctl.Msg, "Data": ctl.Result}
	ctl.Data["json"] = &res
	ctl.ServeJSON()
	ctl.Controller.StopRun()
}

func (ctl *BaseController) GetLanguage() {
	ctl.Lang = ctl.Ctx.Request.Header.Get("Language")
	if len(ctl.Lang) == 0 {
		ctl.Lang = "zh-CN"
	}
}

type PermissionController struct {
	BaseController
}

func (ctl *PermissionController) Prepare() {
	ctl.CheckLogin()
}

func (ctl *PermissionController) CheckLogin() {
	loginToken := ctl.GetUserFromToken()

	ctl.AccountId = loginToken.AccountId
	ctl.Username = loginToken.UserName
	requestUrl := ctl.Ctx.Request.URL.String()
	charIndex := strings.Index(requestUrl, "?")
	if charIndex >= 0 {
		requestUrl = requestUrl[0:charIndex]
	}
	ctl.RequestUrl = requestUrl

}

func (ctl *PermissionController) GetUserFromToken() (loginToken LoginToken) {
	token := ctl.Ctx.Input.Header("Token")
	if len(token) > 0 {
		tokenMap := jwt.Parse(token, nacos.TokenSalt)
		accountId, _ := tokenMap["AccountId"].(float64)
		username, _ := tokenMap["UserName"].(string)
		ip, _ := tokenMap["Ip"].(string)
		loginToken = LoginToken{
			AccountId: int64(accountId),
			UserName:  username,
			Ip:        ip,
			Token:     token,
		}
	}

	return
}

type LoginToken struct {
	AccountId int64
	UserName  string
	Ip        string
	Token     string
}
