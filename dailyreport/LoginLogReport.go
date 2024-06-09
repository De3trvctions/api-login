package dailyreport

import (
	"github.com/beego/beego/v2/client/orm"
)

type LoginLogDailyReport struct {
	Id            int64  `orm:"description(主键)"`
	ReportDay     string `orm:"description(创建时间)"`
	ReportDayTime uint64 `orm:"description(修改时间)"`
	Username      string `orm:"description(用户ID)"`
	LoginCount    int    `orm:"description(登录次数)"`
}

func init() {
	orm.RegisterModel(new(LoginLogDailyReport))
}

func (lg *LoginLogDailyReport) TableName() string {
	return "api_login_log_daily_report"
}
