package config

import (
	"api-login/redis"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

// TODO: Change to nacos
func InitLogs() {
	logs.SetLogger("console", `{"level":7,"color":true}`) // Set level to Trace for maximum verbosity
	logs.EnableFuncCallDepth(true)                        // Enable func call depth to display file and line numbers
	logs.SetLogFuncCallDepth(3)
	logs.Info("[InitLogs] Init Logs Success")
}

func InitRedis() {
	redis.InitRedis()
}

// TODO: Change to nacos
func InitLanguage() {
	langs, _ := web.AppConfig.String("langs")
	logs.Info(langs)
	langTypes := strings.Split(langs, "|")
	for _, lang := range langTypes {
		if lang != "" {
			logs.Info("[InitLanguage] Initialize language: ", lang)
			if err := i18n.SetMessage(lang, "conf/locale_"+lang+".ini"); err != nil {
				logs.Error("[InitLanguage] Fail to set message file:", err)
			}
		}
	}
	logs.Info("[InitLanguage] Init Language Success")
}
