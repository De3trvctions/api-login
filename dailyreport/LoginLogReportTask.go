package dailyreport

import (
	"api-login/consts"
	"api-login/nacos"
	"api-login/redis"
	"api-login/utility"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type LoginLogReportTask struct{}

func RunDailyLoginLogReport() {
	task := LoginLogReportTask{}
	task.RunDailyReport()
}

func (dailyReport *LoginLogReportTask) RunDailyReport() {
	isOk, _ := redis.SetNx(consts.LoginLogTaskLock, true, 3600*time.Second)
	defer redis.Del(consts.LoginLogTaskLock)
	if !isOk {
		logs.Error("[DailyReport][LoginLogReportTask] Report ececuting...")
		return
	}
	cacheKey := fmt.Sprintf(consts.LoginLogPrefix)
	lastDayTime, _ := redis.Get(cacheKey)
	logs.Info("[DailyReport][LoginLogReportTask] Last Execute Time: %s", lastDayTime)

	nowDayTime, _ := utility.TimeStartOfDay()
	runTime := nowDayTime.Add(-1 * 24 * time.Hour).Format(time.DateTime)
	if lastDayTime != "" && lastDayTime == runTime {
		logs.Info("[DailyReport][LoginLogReportTask] Not Time Yet!!!!!!", lastDayTime, runTime)
		return
	}

	var dayTimeList []string
	if lastDayTime == "" {
		dayTimeList = append(dayTimeList, runTime)
	} else {
		diff, _ := utility.TimeStringBetween(lastDayTime, runTime)
		for i := 0; i < diff; i++ {
			curTime := nowDayTime.Add(time.Duration(-(diff-i)*24) * time.Hour).Format(time.DateOnly)
			dayTimeList = append(dayTimeList, curTime)
		}
	}

	for _, t := range dayTimeList {
		dailyReport.dailyEventGiftCodeReport(t)
		// 写入Redis
		_ = redis.Set(cacheKey, t, -1)
		logs.Info("[DailyReport][LoginLogReportTask] Done Login Log Daily report, Date: %s", t)
	}
}

func (dailyReport *LoginLogReportTask) Repair(nowDayTime string) {
	isOk, _ := redis.SetNx(consts.LoginLogRepairTaskLock, true, 3600)
	if !isOk {
		logs.Error("[DailyReport][LoginLogReportTask] Repair Login Log Executing...")
		return
	}
	defer redis.Del(consts.LoginLogRepairTaskLock)
	dailyReport.dailyEventGiftCodeReport(nowDayTime)
}

func (dailyReport *LoginLogReportTask) dailyEventGiftCodeReport(nowDayTime string) {
	beginTime := fmt.Sprintf("%s 00:00:00", nowDayTime[0:10])
	endTime := fmt.Sprintf("%s 23:59:59", nowDayTime[0:10])
	db := utility.NewDB()
	insertData, err := dailyReport.getLoginLogSummaryData(beginTime, endTime)
	if err != nil {
		logs.Error("[DailyReport][LoginLogReportTask] Report Generate error 1", err)
		return
	}
	for i := range insertData {
		insertData[i].ReportDay = nowDayTime[0:10]
		insertData[i].ReportDayTime = utility.TimeParseWithoutError(fmt.Sprintf("%s 00:00:00", beginTime[0:10]), time.DateTime)
	}

	//判断批量插入有没有数据，如果存在 则执行批量插入
	if lenI := len(insertData); lenI > 0 {
		_, err := db.Raw("DELETE FROM api_login_log_daily_report WHERE report_day = ?").SetArgs(beginTime[:10]).Exec()
		if err != nil {
			logs.Error("[DailyReport][LoginLogReportTask] Delete Error", err)
			return
		}
		tx, _ := db.Begin()
		defer tx.Commit()
		_, err = tx.InsertMulti(60, &insertData)
		if err != nil {
			tx.Rollback()
			logs.Error("[DailyReport][LoginLogReportTask] Insert Error:", err)
			return
		}
	}
	logs.Info("[DailyReport][LoginLogReportTask] Done Report")
}

func (dailyReport *LoginLogReportTask) getLoginLogSummaryData(startTime, endTime string) (summaryData []LoginLogDailyReport, err error) {
	beginTimeInt, _ := time.Parse(startTime, time.DateOnly)
	endTimeInt, _ := time.Parse(endTime, time.DateOnly)
	qbSelect, _ := orm.NewQueryBuilder(nacos.DBDriver)
	qbFrom, _ := orm.NewQueryBuilder(nacos.DBDriver)
	qbWhere, _ := orm.NewQueryBuilder(nacos.DBDriver)
	db := utility.NewDB()
	var args []interface{}

	qbSelect.Select("B.username, COUNT(A.id) AS login_count")
	qbFrom.From("api_login_log AS A")
	qbFrom.InnerJoin("api_account AS B").On("A.user_id = B.id")
	qbWhere.Where("A.create_time >= ?").And("A.create_time <= ?")
	args = append(args, beginTimeInt, endTimeInt)

	sql := qbSelect.String() + " " + qbFrom.String() + " " + qbWhere.String()
	_, err = db.Raw(sql).SetArgs(args).QueryRows(&summaryData)
	if err != nil {
		logs.Error("[DailyReport][LoginLogReportTask] Query err", sql, args, err)
		return
	}

	return
}
