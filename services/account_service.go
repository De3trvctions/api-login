package services

import (
	"errors"
	"fmt"

	"api-login/models"
	"api-login/models/dto"
	"api-login/repositories"
	"standard-library/consts"
	"standard-library/redis"
	"standard-library/utility"

	"github.com/beego/beego/v2/core/logs"
)

type AccountService struct {
	Repo *repositories.AccountRepository
}

func NewAccountService() *AccountService {
	return &AccountService{Repo: repositories.NewAccountRepository()}
}

func normalizeEditRequest(accountID int64, req dto.ReqEditAccount) dto.ReqEditAccount {
	req.AccountId = accountID
	return req
}

func (s *AccountService) Register(req dto.ReqRegister) (int64, error) {
	count, err := s.Repo.CountByUsername(req.Username)
	if err != nil {
		logs.Error("[AccountService][Register] count username error", err)
		return consts.DB_GET_FAILED, err
	}
	if count > 0 {
		return consts.USERNAME_DUP, errors.New("duplicated username")
	}

	account := models.Account{
		Username: req.Username,
		Email:    req.Email,
	}
	if err := account.SetHashPassword(req.Password); err != nil {
		logs.Error("[AccountService][Register] hash password error", err)
		return consts.OPERATION_FAILED, err
	}
	account.SetCreateTime()

	if err := s.Repo.Create(&account); err != nil {
		logs.Error("[AccountService][Register] create account error", err)
		return consts.DB_INSERT_FAILED, err
	}
	return 0, nil
}

func (s *AccountService) EditProfile(accountID int64, req dto.ReqEditAccount) (int64, error) {
	req = normalizeEditRequest(accountID, req)
	account, err := s.Repo.GetByID(req.AccountId)
	if err != nil {
		logs.Error("[AccountService][EditProfile] get account error", err)
		return consts.DB_GET_FAILED, err
	}

	fields := make([]string, 0, 4)
	if req.Email != "" && req.Email != account.Email {
		if errCode, err := validateEmailCode(consts.RegisterEmailValidCode, req.Email, req.ValidCode); err != nil || errCode != 0 {
			return errCode, err
		}
		account.Email = req.Email
		fields = append(fields, "Email")
	}

	if req.CountryCode != 0 && req.Phone != 0 && (req.CountryCode != account.CountryCode || req.Phone != account.Phone) {
		account.CountryCode = req.CountryCode
		account.Phone = req.Phone
		fields = append(fields, "CountryCode", "Phone")
	}

	if req.NewPassword != "" && !account.VerifyPassword(req.NewPassword) {
		if err := account.SetHashPassword(req.NewPassword); err != nil {
			logs.Error("[AccountService][EditProfile] hash password error", err)
			return consts.OPERATION_FAILED, err
		}
		fields = append(fields, "Password")
	}

	if len(fields) == 0 {
		return 0, nil
	}

	account.SetUpdateTime()
	fields = append(fields, "UpdateTime")
	if err := s.Repo.Update(&account, fields...); err != nil {
		logs.Error("[AccountService][EditProfile] update account error", err)
		return consts.DB_UPDATE_FAILED, err
	}
	return 0, nil
}

func (s *AccountService) ResetPassword(req dto.ReqForgetPasswordSetNew) (int64, error) {
	account, err := s.Repo.GetByUsername(req.Username)
	if err != nil {
		logs.Error("[AccountService][ResetPassword] get account error", err)
		return consts.USERNAME_NOT_FOUND, err
	}
	if req.Email != account.Email {
		return consts.FORGET_PASSWORD_EMAIL_NOT_MATCH, errors.New("email not match")
	}
	if errCode, err := validateEmailCode(consts.ForgetPasswordEmailValidCode, req.Email, req.ValidCode); err != nil || errCode != 0 {
		return errCode, err
	}
	if err := account.SetHashPassword(req.Password); err != nil {
		logs.Error("[AccountService][ResetPassword] hash password error", err)
		return consts.OPERATION_FAILED, err
	}
	account.SetUpdateTime()
	if err := s.Repo.Update(&account, "Password", "UpdateTime"); err != nil {
		logs.Error("[AccountService][ResetPassword] update account error", err)
		return consts.DB_UPDATE_FAILED, err
	}
	return 0, nil
}

func validateEmailCode(keyPattern, email, validCode string) (int64, error) {
	exists, _ := redis.Exists(fmt.Sprintf(keyPattern, email))
	if !exists {
		return consts.VALID_CODE_NOT_MATCH, errors.New("valid code not found")
	}
	defer utility.DelEmailValidCodeLock(email)

	cachedCode, _ := redis.Get(fmt.Sprintf(keyPattern, email))
	if cachedCode != validCode {
		return consts.VALID_CODE_NOT_MATCH, errors.New("valid code not match")
	}
	return 0, nil
}
