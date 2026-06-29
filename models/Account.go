package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	CommStruct
	Username    string `orm:"description(用户名)"`
	Password    string `orm:"description(用户密码)"`
	Email       string `orm:"description(邮件)"`
	Phone       int    `orm:"description(手机号码)"`
	CountryCode int    `orm:"description(国家地区编码)"`
}

type AccountInfo struct {
	CommStruct
	Username    string
	Email       string
	Phone       int
	CountryCode int
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
	return "api_account"
}

func (acc *Account) SetHashPassword(password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	acc.Password = string(hashPassword)
	return nil
}

func (acc *Account) VerifyPassword(password string) bool {
	if acc.Password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)) == nil
}
