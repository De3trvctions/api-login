package config

import (
	"api-login/utility"

	"github.com/beego/beego/v2/core/config/env"
	"github.com/beego/beego/v2/server/web"
)

var (
	HttpPort          int64
	TokenSalt         string
	TokenExpMinute    int64
	TokenMaxExpSecond int64
)

func init() {
	HttpPort = utility.StringToInt64(env.Get("HttpPort", web.AppConfig.DefaultString("HttpPort", "8080")))
	TokenSalt = env.Get("TokenSalt", "")
	TokenExpMinute = utility.StringToInt64(env.Get("TokenExpMinute", "1440"))
	TokenMaxExpSecond = utility.StringToInt64(env.Get("TokenMaxExpSecond", "86400"))
}
