package db

import (
	"database/sql"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	dbDriver, _ := web.AppConfig.String("db.driver")
	dbUser, _ := web.AppConfig.String("db.user")
	dbPass, _ := web.AppConfig.String("db.passwd")
	dbHost, _ := web.AppConfig.String("db.host")
	dbPort, _ := web.AppConfig.String("db.port")
	dbName, _ := web.AppConfig.String("db.name")
	// Construct data source name (DSN)
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"

	// Register MySQL database driver
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// Register default database
	orm.RegisterDataBase("default", dbDriver, dsn)

	// Open a database connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logs.Error("[InitDB] Open DB fail")
		return
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		logs.Error("[InitDB] Ping DB fail")
		return
	}

	logs.Info("[InitDB] Init DB Success")
}

func GetDB() *sql.DB {
	return db
}

// func aaa() {
// 	aa, _ := orm.NewQueryBuilder("mysql")

// 	aa.Select()

// }
