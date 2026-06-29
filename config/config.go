package config

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// Use when ever data is in app.conf
var (
	HttpPort         int64 = 8080
	NacosPort        int64
	NacosUrl         string
	NacosNamespaceId string
	NacosDataId      string
	NacosGroupId     string
)

func Load() error {
	HttpPort = getInt64Value("HttpPort", HttpPort)
	NacosPort = getInt64Value("NacosPort", NacosPort)
	NacosUrl = getValue("NacosUrl")
	NacosNamespaceId = getValue("NacosNamespaceId")
	NacosDataId = getValue("NacosDataId")
	NacosGroupId = getValue("NacosGroupId")
	return nil
}

func getInt64Value(key string, fallback int64) int64 {
	value := getValue(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		logs.Error("[Conf][getInt64Value]Error", err)
		return fallback
	}
	return parsed
}

func getValue(key string) string {
	value, err := web.AppConfig.String(key)
	if err != nil {
		logs.Error("[Conf][getValue]Error", err)
		return ""
	}
	return value
}
