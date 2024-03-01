package dto

type ReqForgetPassword struct {
	Username string `valid:"Required;IsUsername"`
	Email    string `valid:"Required;Email"`
}
