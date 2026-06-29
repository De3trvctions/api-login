package repositories

import (
	"strings"

	"api-login/models"
	"api-login/models/dto"
	"api-login/utility"
)

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func buildDetailQuery(req dto.ReqAccountDetail) (string, []interface{}) {
	sql := "SELECT id, create_time, update_time, deleted, username, email, phone, country_code FROM api_account WHERE 1=1"
	args := make([]interface{}, 0, 4)
	if req.AccountId > 0 {
		sql += " AND id = ?"
		args = append(args, req.AccountId)
	}
	if strings.TrimSpace(req.Username) != "" {
		sql += " AND username = ?"
		args = append(args, req.Username)
	}
	if strings.TrimSpace(req.Email) != "" {
		sql += " AND email = ?"
		args = append(args, req.Email)
	}
	if req.CreateTime > 0 {
		sql += " AND create_time > ?"
		args = append(args, req.CreateTime)
	}
	return sql, args
}

func (r *AccountRepository) Detail(req dto.ReqAccountDetail) (models.AccountInfo, error) {
	var out models.AccountInfo
	sql, args := buildDetailQuery(req)
	err := utility.NewDB().Raw(sql).SetArgs(args).QueryRow(&out)
	return out, err
}

func (r *AccountRepository) List(req dto.ReqAccountDetail) ([]models.AccountInfo, error) {
	var out []models.AccountInfo
	sql, args := buildDetailQuery(req)
	_, err := utility.NewDB().Raw(sql).SetArgs(args).QueryRows(&out)
	return out, err
}

func (r *AccountRepository) GetByID(id int64) (models.Account, error) {
	account := models.Account{}
	account.Id = id
	err := utility.NewDB().Get(&account, "Id")
	return account, err
}

func (r *AccountRepository) GetByUsername(username string) (models.Account, error) {
	account := models.Account{}
	account.Username = username
	err := utility.NewDB().Get(&account, "Username")
	return account, err
}

func (r *AccountRepository) Create(account *models.Account) error {
	_, err := utility.NewDB().Insert(account)
	return err
}

func (r *AccountRepository) Update(account *models.Account, fields ...string) error {
	_, err := utility.NewDB().Update(account, fields...)
	return err
}

func (r *AccountRepository) CountByUsername(username string) (int64, error) {
	account := models.Account{}
	return utility.NewDB().Count(&account, "username", username)
}
