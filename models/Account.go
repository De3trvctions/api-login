package models

import (
	"api-login/consts"
	"api-login/models/dto"
	"api-login/utility"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type Account struct {
	CommStruct
	Username string `orm:"description(用户名)"`
	Password string `orm:"description(用户密码)"`
	Email    string `orm:"description(邮件)"`
}

func (acc *Account) SetUpdateTime() {
	acc.UpdateTime = uint64(time.Now().Unix())
}

func (acc *Account) SetCreateTime() {
	acc.CreateTime = uint64(time.Now().Unix())
}

func init() {
	orm.RegisterModel(new(Account))
}

func (acc *Account) TableName() string {
	return "cloud_data_account"
}

func (acc *Account) List(req Account) (accountList []Account, errCode int, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	db := utility.NewDB()

	qb.Select("*")
	qb.From(acc.TableName())
	qb.Where("1=1")
	var args []interface{}

	req.Id = 1
	if req.Id > 0 {
		qb.And("id = ?")
		args = append(args, req.Id)
	}

	if req.Username != "" {
		qb.And("username = ?")
		args = append(args, req.Username)
	}

	if req.CreateTime > 0 {
		qb.And("create_time > ?")
		args = append(args, req.CreateTime)
	}

	if req.Email != "" {
		qb.And("email > ?")
		args = append(args, req.Email)
	}

	sql := qb.String()
	_, err = db.Raw(sql).SetArgs(args).QueryRows(&accountList)
	if err != nil {
		logs.Error("[Account][List] Query error:", sql, args, err)
	}

	return
}

func (acc *Account) Register(req dto.ReqRegister) (errCode int64, err error) {
	// Check if username Exist
	db := utility.NewDB()
	acc.Username = req.Username

	count, err := db.Count(acc, "username", acc.Username)
	if err != nil {
		errCode = consts.DB_GET_FAILED
		logs.Error("[Account][Register] Check Username Error, ", err)
		return
	}
	if count > 0 {
		errCode = consts.USERNAME_DUP
		err = errors.New("duplicated username")
		return
	}

	// Hashing Password
	hash := md5.Sum([]byte(req.Password))
	hashPassword := hex.EncodeToString(hash[:])

	// Assign Value
	acc.Username = req.Username
	acc.Password = hashPassword
	acc.Email = req.Email
	acc.SetCreateTime()

	// Insert To DB
	_, err = db.Insert(acc)
	if err != nil {
		errCode = consts.DB_INSERT_FAILED
		logs.Error("[Account][Register] Insert Account error", err)
		return
	}

	return
}
