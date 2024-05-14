package repository

import "github.com/lordbritishix/integral/internal/model"

type TransactionsRepository interface {
	SaveTransaction(transaction model.Transaction) error
	GetTransactions(accountId string) (model.Transactions, error)
}
