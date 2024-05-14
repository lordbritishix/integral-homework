package repository

import (
	"github.com/lordbritishix/integral/internal/model"
	"sort"
	"sync"
)

type InMemoryTransactionRepository struct {
	transactions sync.Map
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{}
}

func (i *InMemoryTransactionRepository) SaveTransaction(transaction model.Transaction) error {
	transactions, ok := i.transactions.Load(transaction.Account.AccountId)

	if !ok {
		i.transactions.Store(transaction.Account.AccountId, model.Transactions{transaction})
	} else {
		i.transactions.Store(transaction.Account.AccountId, append(transactions.(model.Transactions), transaction))
	}

	return nil
}

func (i *InMemoryTransactionRepository) GetTransactions(accountId string) (model.Transactions, error) {
	transactions, ok := i.transactions.Load(accountId)

	if !ok {
		return model.Transactions{}, nil
	}

	sort.Sort(transactions.(model.Transactions))

	return transactions.(model.Transactions), nil
}
