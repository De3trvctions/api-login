package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type Account struct {
	CommStruct
	Username string `orm:"description(用户名)"`
	Password string `orm:"description(用户密码)"`
}

func init() {
	orm.RegisterModel(new(Account))
}

func (acc *Account) TableName() string {
	return "cloud_data_account"
}
