package system

import (
	"api-login/consts"
	"api-login/models/dto"
	"api-login/redis"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type LoginController struct {
	BaseController
}

// Login
//
//	@Title			登陆
//	@Description	登陆
//	@Success		200	"success"
//	@router			/login [get]
func (ctl *LoginController) Login() {
	req := dto.ReqLogin{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[Login] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}

	ableLogin, ableLoginRemaindingTime := ctl.getRedisLoginStatus(req.Username)

	if !ableLogin {
		leftSec := ableLoginRemaindingTime % 60
		leftMin := int(ableLoginRemaindingTime / 60)
		logs.Error("[LoginController][Login] 账号封锁中。剩余解封时间%d分钟%d秒", leftMin, leftSec)
		ctl.Error(consts.LOGIN_LOCK, fmt.Sprintf("请在%d分钟%d秒后再进行尝试", leftMin, leftSec))
	}

	ctl.Success("Success")
}

func (ctl *LoginController) getRedisLoginStatus(username string) (ableLogin bool, remaindingTime int) {
	ableLogin = true
	ex1, _ := redis.Exists(fmt.Sprintf(consts.FailLoginAccountLock, username))
	// ex2, _ := redis.Exists(fmt.Sprintf(consts.FailLoginAccountLockTime, username))

	if ex1 {
		ableLogin = false
		timeCache, _ := redis.Get(fmt.Sprintf(consts.FailLoginAccountLockTime, username))
		redisTime, _ := strconv.ParseInt(timeCache, 10, 64)
		timeLeft := time.Unix(redisTime, 0)
		remaindingTime = int(time.Until(timeLeft).Seconds())
	}
	return
}

func (ctl *LoginController) setRedisLoginFail(username string) (err error) {
	failCount := 0
	ex, _ := redis.Exists(fmt.Sprintf(consts.FailLoginCount, username))
	if ex {
		count, _ := redis.Get(fmt.Sprintf(consts.FailLoginCount, username))
		failCount, _ = strconv.Atoi(count)
		_, _ = redis.Del(fmt.Sprintf(consts.FailLoginCount, username))
	}

	if failCount >= 5 {
		_ = redis.Set(fmt.Sprintf(consts.FailLoginCount, username), failCount, time.Duration(15)*time.Minute)
		_ = redis.Set(fmt.Sprintf(consts.FailLoginAccountLock, username), 1, time.Duration(15)*time.Minute)
		_ = redis.Set(fmt.Sprintf(consts.FailLoginAccountLockTime, username), time.Now().Add(time.Duration(15)*time.Minute).Unix(), time.Duration(15)*time.Minute)
	} else {
		_ = redis.Set(fmt.Sprintf(consts.FailLoginCount, username), failCount)
	}
	return
}

func (ctl *LoginController) delRedisLoginFail(username string) {
	_, _ = redis.Del(fmt.Sprintf(consts.FailLoginCount, username))
	_, _ = redis.Del(fmt.Sprintf(consts.FailLoginAccountLock, username))
	_, _ = redis.Del(fmt.Sprintf(consts.FailLoginAccountLockTime, username))
}
