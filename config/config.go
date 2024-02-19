package config

import (
	"strconv"

	"github.com/beego/beego/v2/core/config/env"
	"github.com/beego/beego/v2/server/web"
)

var (
	HttpPort int64
)

func init() {
	num, _ := strconv.Atoi(env.Get("HttpPort", web.AppConfig.DefaultString("HttpPort", "8080")))
	HttpPort = int64(num)
}
