package models

type CommStruct struct {
	Id            int64  `orm:"description(主键)"`
	AdminId       int64  `orm:"description(所有人)"`
	CreateAdminId int64  `orm:"description(创建人)"`
	CreateTime    uint64 `orm:"description(创建时间)"`
	UpdateTime    uint64 `orm:"description(修改时间)"`
	Deleted       int    `orm:"description(0:正常  1:删除)"`
}
