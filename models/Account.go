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

type AccountInfo struct {
	CommStruct
	Username string
	Email    string
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

func (acc *Account) SetHashPassword(password string) {
	hash := md5.Sum([]byte(password))
	acc.Password = hex.EncodeToString(hash[:])
}

func (acc *Account) GetHashPassword(password string) (hashPassword string) {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (acc *Account) List(req Account) (accountList []AccountInfo, errCode int, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	db := utility.NewDB()

	qb.Select("*")
	qb.From(acc.TableName())
	qb.Where("1=1")
	var args []interface{}

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

func (acc *Account) Info(req dto.ReqAccountDetail) (account AccountInfo, errCode int64, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qbWhere, _ := orm.NewQueryBuilder("mysql")
	var args []interface{}
	db := utility.NewDB()

	qb.Select("*")
	qb.From(acc.TableName())
	qbWhere.Where("1=1")

	if req.AccountId > 0 {
		qbWhere.And("id = ?")
		args = append(args, req.AccountId)
	}

	if req.Username != "" {
		qbWhere.And("username = ?")
		args = append(args, req.Username)
	}

	if req.CreateTime > 0 {
		qbWhere.And("create_time > ?")
		args = append(args, req.CreateTime)
	}

	if req.Email != "" {
		qbWhere.And("email > ?")
		args = append(args, req.Email)
	}

	sql := qb.String() + " " + qbWhere.String()
	err = db.Raw(sql).SetArgs(args).QueryRow(&account)
	if err != nil {
		errCode = consts.DB_GET_FAILED
		logs.Error("[Account][SelfInfo] Query error:", sql, args, err)
	}

	return
}

func (acc *Account) SelfInfo() (account AccountInfo, errCode int64, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qbWhere, _ := orm.NewQueryBuilder("mysql")
	var args []interface{}
	db := utility.NewDB()

	qb.Select("*")
	qb.From(acc.TableName())
	qbWhere.Where("id = ?")
	args = append(args, acc.Id)

	sql := qb.String() + " " + qbWhere.String()
	err = db.Raw(sql).SetArgs(args).QueryRow(&account)
	if err != nil {
		errCode = consts.DB_GET_FAILED
		logs.Error("[Account][SelfInfo] Query error:", sql, args, err)
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

	// Assign Value
	acc.Username = req.Username
	acc.SetHashPassword(req.Password)
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

func (acc *Account) Edit(req dto.ReqEditAccount) (errCode int64, err error) {
	db := utility.NewDB()
	if acc.Id > 0 {
		err = db.Get(acc, "Id")
		if err != nil {
			logs.Error("[Account][Edit] Get account error", err)
			errCode = consts.DB_GET_FAILED
			return
		}
	} else {
		logs.Error("[Account][Edit] No Account Id")
		errCode = consts.ACCOUNT_ID_INVALID
		return
	}

	var updateField []string

	if req.Email != acc.Email {
		acc.Email = req.Email
		updateField = append(updateField, "Email")
	}

	hashPassword := acc.GetHashPassword(req.NewPassword)
	if hashPassword != acc.Password {
		acc.Password = hashPassword
		updateField = append(updateField, "Password")
	}

	_, err = db.Update(acc, updateField...)
	if err != nil {
		logs.Error("[Account][Edit] Update account error. %+v , error: %+v", acc, err)
		errCode = consts.DB_UPDATE_FAILED
		return
	}

	return
}
