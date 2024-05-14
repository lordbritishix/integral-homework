package repository

import "github.com/lordbritishix/integral/internal/model"

type InMemoryAccountRepository struct {
	accounts map[string]model.Account
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[string]model.Account),
	}
}

func (r *InMemoryAccountRepository) SaveAccount(account model.Account) error {
	r.accounts[account.AccountId] = account
	return nil
}

func (r *InMemoryAccountRepository) GetAccount(accountId string) (*model.Account, error) {
	account := r.accounts[accountId]
	return &account, nil
}

func (r *InMemoryAccountRepository) Exists(accountId string) (bool, error) {
	_, ok := r.accounts[accountId]

	return ok, nil
}
