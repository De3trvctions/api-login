package system

import (
	"api-login-proto/common"
	"context"
	"errors"
	"fmt"
	"standard-library/consts"
	"standard-library/grpc"
	"standard-library/json"
	"standard-library/jwt"
	"standard-library/nacos"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	ogrpc "google.golang.org/grpc"
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

	GrpcConn        grpc.Conn
	GrpcCtx         context.Context
	GrpcServiceName string
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

func (ctl *BaseController) ConnGRpc(serviceName string) {
	ctl.GrpcServiceName = serviceName
	var err error
	ctl.GrpcConn, err = grpc.Get(serviceName)
	if err != nil {
		logs.Error("[tsGrpc.GetGrpcClient]获取服务[%s]失败 %s", serviceName, err)
		ctl.Error(consts.SERVER_ERROR)
	}

	// 针对异常进行处理
	if ctl.GrpcConn == nil {
		logs.Error("[tsGrpc.GetGrpcClient]获取服务[%s]conn为空，重新获取服务", serviceName)
		if err == nil {
			err = errors.New("conn is nil")
		}
		ctl.GrpcConn, _ = grpc.Get(serviceName)
	}
}

func (ctl *BaseController) CommonJSONRequest(data interface{}, api func(context.Context, *common.Request, ...ogrpc.CallOption) (*common.Response, error)) {
	ctl.CommonRequestMicroService(json.String(data), api)
}

func (ctl *BaseController) CommonRequestMicroService(data string, api func(context.Context, *common.Request, ...ogrpc.CallOption) (*common.Response, error)) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	request := &common.Request{
		Data: data,
	}

	response, _ := api(ctx, request)
	logs.Error("阿四大大大大", response)
	if response.Code != consts.SUCCESS_REQUEST && response.Code != 0 {
		if response.Data != "" {
			logs.Error("132", response)
			ctl.Error(response.Code, response.Data)
		} else {
			logs.Error("7655", response)
			ctl.Error(response.Code)
		}
	}
	logs.Error("fasfasas")
	ctl.DirectWriteJson(response.Data)
}

// DirectWriteJson 接受服务的返回，直接输出到前端 | 其一：返回前端数据，需要加密
func (ctl *BaseController) DirectWriteJson(data string) {
	ctl.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	logs.Error("11111", data)
	ctl.Ctx.WriteString(data)
	ctl.StopRun()
}

// SuccessJson 接受服务的返回，包装并直接输出到前端 | 其一：返回前端数据，需要加密
func (ctl *BaseController) SuccessJson(data string) {
	ctl.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	if data == "" {
		data = fmt.Sprintf(`{"Code": 200, "Msg": "%s", "Data": null}`, ctl.Msg)
	} else {
		data = fmt.Sprintf(`{"Code": 200, "Msg": "%s", "Data": %s}`, ctl.Msg, data)
	}

	ctl.Ctx.WriteString(data)

	ctl.StopRun()
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
