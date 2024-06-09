package system

import (
	"api-login/consts"
	"api-login/dailyreport"
	"api-login/utility"
	"api-login/validation"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type FixController struct {
	BaseController
}

// FixDailyLoginLogReport
//
//	@Title			修复登录日志日报
//	@Description	修复登录日志日报
//	@Success		200			{string}	"success"
//	@Param			BeginDate	query		string	true	Report	Begin	Date
//	@Param			EndDate		query		string	true	Report	End		Date
//	@router			/fixdailiyloginlogreport [get]
func (ctl *FixController) FixDailyLoginLogReport() {
	type reqFixLoginReport struct {
		BeginDate string `valid:"Required;"`
		EndDate   string `valid:"Required;"`
	}
	req := reqFixLoginReport{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[FixController][FixDailyLoginLogReport] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[FixController][FixDailyLoginLogReport] Form Validate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	var reportDateList []string
	nowDayTime, _ := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:00", req.BeginDate[0:10]))
	diff, _ := utility.TimeStringBetween(req.EndDate[0:10], req.BeginDate[0:10])
	for i := 0; i < diff; i++ {
		curTime := nowDayTime.Add(time.Duration(-(diff-i)*24) * time.Hour).Format(time.DateOnly)
		reportDateList = append(reportDateList, curTime)
	}

	loginReport := dailyreport.LoginLogReportTask{}

	for _, v := range reportDateList {
		loginReport.Repair(v)
	}

	ctl.Success("Success")
}
