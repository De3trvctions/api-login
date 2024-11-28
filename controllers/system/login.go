package system

import (
	"api-login-proto/login"
	"fmt"
	"standard-library/consts"
	"standard-library/models/dto"
	"standard-library/redis"
	"standard-library/validation"

	"github.com/beego/beego/v2/core/logs"
)

type LoginController struct {
	PermissionController
	GRPCClient login.UserLoginServiceClient
}

func (ctl *LoginController) Prepare() {
	ctl.PermissionController.Prepare()
	//初始化连接
	ctl.ConnGRpc("service-login")
	ctl.GRPCClient = login.NewUserLoginServiceClient(ctl.GrpcConn.Conn())
	if ctl.GRPCClient == nil {
		ctl.Error(consts.SERVER_ERROR)
	}
}

// Login
//
//	@Title			登陆
//	@Description	登陆
//	@Success		200			"success"
//	@Param			Username	formData	string	true	账号
//	@Param			Password	formData	string	true	用户密码
//	@Param			DeviceId	formData	string	false	设备号
//	@Param			IP			formData	string	false	IP地址
//	@router			/ [post]
func (ctl *LoginController) Login() {

	req := dto.ReqLogin{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[LoginController][Login] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[LoginController][Login]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	ctl.CommonJSONRequest(req, ctl.GRPCClient.Login)
}

// func (ctl *LoginController) getRedisLoginStatus(username string) (ableLogin bool, remaindingTime int) {
// 	ableLogin = true
// 	ex1, _ := redis.Exists(fmt.Sprintf(consts.FailLoginAccountLock, username))

// 	if ex1 {
// 		ableLogin = false
// 		timeCache, _ := redis.Get(fmt.Sprintf(consts.FailLoginAccountLockTime, username))
// 		redisTime, _ := strconv.ParseInt(timeCache, 10, 64)
// 		timeLeft := time.Unix(redisTime, 0)
// 		remaindingTime = int(time.Until(timeLeft).Seconds())
// 	}
// 	return
// }

// func (ctl *LoginController) setRedisLoginFail(username string) (err error) {
// 	failCount := 0
// 	ex, _ := redis.Exists(fmt.Sprintf(consts.FailLoginCount, username))
// 	if ex {
// 		count, _ := redis.Get(fmt.Sprintf(consts.FailLoginCount, username))
// 		failCount, _ = strconv.Atoi(count)
// 		_, _ = redis.Del(fmt.Sprintf(consts.FailLoginCount, username))
// 	}

// 	if failCount >= 5 {
// 		_ = redis.Set(fmt.Sprintf(consts.FailLoginCount, username), failCount, time.Duration(15)*time.Minute)
// 		_ = redis.Set(fmt.Sprintf(consts.FailLoginAccountLock, username), 1, time.Duration(15)*time.Minute)
// 		_ = redis.Set(fmt.Sprintf(consts.FailLoginAccountLockTime, username), time.Now().Add(time.Duration(15)*time.Minute).Unix(), time.Duration(15)*time.Minute)
// 	} else {
// 		_ = redis.Set(fmt.Sprintf(consts.FailLoginCount, username), failCount)
// 	}
// 	return
// }

// func (ctl *LoginController) delRedisLoginFail(username string) {
// 	_, _ = redis.Del(fmt.Sprintf(consts.FailLoginCount, username))
// 	_, _ = redis.Del(fmt.Sprintf(consts.FailLoginAccountLock, username))
// 	_, _ = redis.Del(fmt.Sprintf(consts.FailLoginAccountLockTime, username))
// }

// func getToken(req dto.ReqLogin, id int64) (token string) {
// 	// tokenSalt, _ := web.AppConfig.String("TokenSalt")
// 	// tokenExpMinute, _ := web.AppConfig.Int("TokenExpMinute")
// 	delToken(req.Username)
// 	token = jwt.Gen(web.M{
// 		"Username":  req.Username,
// 		"AccountId": id,
// 		"Ip":        req.IP,
// 	}, nacos.TokenSalt, time.Duration(nacos.TokenExpMinute)*time.Minute)
// 	setToken(token, req.Username)
// 	return
// }

// func setToken(token, username string) {
// 	// Set token expired as 1 days
// 	_ = redis.Set(fmt.Sprintf(consts.AccountLoginByToken, token), username, time.Duration(nacos.TokenExpMinute)*time.Minute)
// 	_ = redis.Set(fmt.Sprintf(consts.AccountLoginByUsername, username), token, time.Duration(nacos.TokenExpMinute)*time.Minute)
// }

func delToken(username string) {
	ex, _ := redis.Exists(fmt.Sprintf(consts.AccountLoginByUsername, username))
	if ex {
		token, _ := redis.Get(fmt.Sprintf(consts.AccountLoginByUsername, username))
		_, err1 := redis.Del(fmt.Sprintf(consts.AccountLoginByToken, token))
		_, err2 := redis.Del(fmt.Sprintf(consts.AccountLoginByUsername, username))
		if err1 != nil {
			logs.Error("[delToken] Error 1: ", err1)
		}
		if err2 != nil {
			logs.Error("[delToken] Error 2: ", err2)
		}
	}
}

// Logout
//
//	@Title			Logout
//	@Description	Logout
//	@Success		200			"success"
//	@router			/ [put]
func (ctl *LoginController) Logout() {
	defer logs.Info("[LoginController][Logout] Username:", ctl.Username)

	loginToken := ctl.GetUserFromToken()
	delToken(loginToken.UserName)

	ctl.Success("success")
}
