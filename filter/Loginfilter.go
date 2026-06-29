package filter

import (
	"fmt"
	"standard-library/consts"
	"standard-library/jwt"
	"standard-library/nacos"
	"standard-library/redis"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

var LoginManager = func(ctx *context.Context) {
	// 返回值值
	var responseText = `{"Code":10008,"Msg":"Not Login","Data":null}`
	token := ctx.Input.Header("Token")
	requestUrl := ctx.Request.URL.String()
	charIndex := strings.Index(requestUrl, "?")
	if charIndex >= 0 {
		requestUrl = requestUrl[0:charIndex]
	}

	// if requestUrl == "/system/user/login" {
	// 	// 不做拦截处理
	// 	return
	// }

	if len(token) == 0 {
		logs.Error("[LoginFilter] unauthorized request_url=%s token_present=%t", requestUrl, token != "")
		goto NoLogin
	} else {
		tokenMap := jwt.Parse(token, nacos.TokenSalt)
		if tokenMap == nil {
			logs.Error("[LoginFilter] unauthorized request_url=%s token_present=%t", requestUrl, token != "")
			goto NoLogin
		}

		// 获取对应的用户信息
		accountId, exist := tokenMap["AccountId"]
		if !exist || accountId == nil {
			logs.Error("[LoginFilter] unauthorized request_url=%s token_present=%t", requestUrl, token != "")
			goto NoLogin
		}

		// 超过最长登陆时间，需要重新登陆
		username, err := redis.Get(fmt.Sprintf(consts.AccountLoginByToken, token))
		if err != nil || len(username) == 0 {
			logs.Error("[LoginFilter] unauthorized request_url=%s account_id=%v token_present=%t", requestUrl, accountId, token != "")
			goto NoLogin
		} else {
			redisToken, err := redis.Get(fmt.Sprintf(consts.AccountLoginByUsername, username))
			if err != nil || len(redisToken) == 0 {
				logs.Error("[LoginFilter] unauthorized request_url=%s username=%s token_present=%t", requestUrl, username, token != "")
				goto NoLogin
			} else {
				if redisToken != token {
					logs.Error("[LoginFilter] unauthorized request_url=%s username=%s token_present=%t", requestUrl, username, token != "")
					goto NoLogin
				}
			}
		}
	}

	return

NoLogin:
	response := ctx.ResponseWriter
	response.Header().Set("content-type", "application/json")
	_, _ = response.Write([]byte(responseText))
}
