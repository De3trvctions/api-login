package dto

type ReqAccountDetail struct {
	AccountId  int64 `valid:"Required"`
	Username   string
	Email      string
	CreateTime int64
}
