package dto

type ReqRegister struct {
	Username string `valid:"IsUsername"`
	Password string
	Email    string `valid:"Email"`
}
