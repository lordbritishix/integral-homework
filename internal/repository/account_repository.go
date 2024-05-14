package repository

import "github.com/lordbritishix/integral/internal/model"

type AccountRepository interface {
	SaveAccount(account model.Account) error
	GetAccount(accountId string) (*model.Account, error)
	Exists(accountId string) (bool, error)
}
