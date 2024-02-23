package filter

import (
	"api-login/config"
	"api-login/consts"
	"api-login/jwt"
	"api-login/redis"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

var LoginManager = func(ctx *context.Context) {
	// 返回值值
	var responseText = `{"Code":10008,"Msg":"未登陆"}`
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
		logs.Error("[LoginFilter]Token Empty not login[Url](%s) [Token](%s)", requestUrl, token)
		goto NoLogin
	} else {
		tokenMap := jwt.Parse(token, config.TokenSalt)
		//logs.Debug("[LoginFilter]FromToken", token, config.TokenSalt, tokenMap)
		if tokenMap == nil {
			logs.Error("[LoginFilter]FromToken Failed Not Login[Url](%s) [Token](%s)", requestUrl, token)
			goto NoLogin
		}

		// 获取对应的用户信息
		accountId, exist := tokenMap["AccountId"]
		if !exist || accountId == nil {
			logs.Error("[LoginFilter]Get admin Id Failed Not Login[Url](%s) [Token](%s)", requestUrl, token)
			goto NoLogin
		}

		// 超过最长登陆时间，需要重新登陆
		username, err := redis.Get(fmt.Sprintf(consts.AccountLoginByToken, token))
		if err != nil || len(username) == 0 {
			logs.Error("[LoginFilter]Token Not exist %v, %v", requestUrl, accountId)
			goto NoLogin
		} else {
			redisToken, err := redis.Get(fmt.Sprintf(consts.AccountLoginByUsername, username))
			if err != nil || len(redisToken) == 0 {
				logs.Error("[LoginFilter]Username Not exist %v, %v", requestUrl, username)
				goto NoLogin
			} else {
				if redisToken != token {
					fmt.Print("[LoginFilter] Token not same")
					goto NoLogin
				}
			}
		}
	}

	return

NoLogin:
	response := ctx.ResponseWriter
	response.Header().Add("content-type", "application/json")
	_, _ = response.Write([]byte(responseText))
}
