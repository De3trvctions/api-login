package consts

const (
	SUCCESS_REQUEST = 200
	LOGIN_LOCK      = 1000

	FAILED_REQUEST                  = 10000
	PARAM_ERROR                     = 10001
	USERNAME_DUP                    = 10002
	USERNAME_NOT_FOUND              = 10003
	PASSWORD_NOT_MATCH              = 10004
	ACCOUNT_ID_INVALID              = 10005
	VALID_CODE_EXIST                = 10006
	VALID_CODE_COOL_DOWN            = 10007
	VALID_CODE_NOT_MATCH            = 10008
	FORGET_PASSWORD_EMAIL_NOT_MATCH = 10010

	DB_GET_FAILED    = 20000
	DB_INSERT_FAILED = 20001
	DB_UPDATE_FAILED = 20002
	DB_DELETE_FAILED = 20003

	SERVER_ERROR     = 40000
	OPERATION_FAILED = 40001
	USER_NOT_LOGIN   = 40002
)
