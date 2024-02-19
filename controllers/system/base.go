package system

import (
	"api-login/consts"
	"fmt"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

type BaseController struct {
	web.Controller
	i18n.Locale
	Code     int64
	Result   interface{}
	Ip       string
	DeviceId string
	Msg      string
	Token    string
}

func (ctl *BaseController) Success(obj interface{}) {
	ctl.Code = consts.SUCCESS
	ctl.Result = obj
	ctl.TraceJson(false)
}

func (ctl *BaseController) Error(code int64, msg ...string) {
	if code < 1 {
		code = consts.ERROR
	}
	ctl.Code = code

	ctl.GetLanguage()

	if ctl.Code > 0 {
		ctl.Msg = ctl.Tr("error." + fmt.Sprintf("%d", ctl.Code))
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
