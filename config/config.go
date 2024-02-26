package config

import (
	"api-login/mail"
	"api-login/utility"
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

var (
	HttpPort          int64
	TokenSalt         string
	TokenExpMinute    int64
	TokenMaxExpSecond int64
	Mail              []*mail.Option
)

func init() {
	HttpPort = utility.StringToInt64(getValue("HttpPort"))
	TokenSalt = getValue("TokenSalt")
	TokenExpMinute = utility.StringToInt64(getValue("TokenExpMinute"))
	TokenMaxExpSecond = utility.StringToInt64(getValue("TokenMaxExpSecond"))
	err := json.Unmarshal([]byte(getValue("AppPassword")), &Mail)
	if err != nil {
		logs.Error("Error:", err)
		return
	}
}

func getValue(key string) string {
	value, err := web.AppConfig.String(key)
	if err != nil {
		logs.Error("[Conf][getValue]Error", err)
		return ""
	}
	return value
}
