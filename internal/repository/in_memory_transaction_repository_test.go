package repository

import (
	"github.com/lordbritishix/integral/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInMemoryTransactionRepository_GetTransactions(t *testing.T) {
	repository := NewInMemoryTransactionRepository()

	now := time.Now()

	account1 := model.Account{AccountId: "1"}
	account2 := model.Account{AccountId: "2"}

	tx1 := model.Transaction{
		TransactionId: "1",
		Account:       account1,
		Timestamp:     now,
	}

	tx2 := model.Transaction{
		TransactionId: "2",
		Account:       account1,
		Timestamp:     now.AddDate(0, 0, 1),
	}

	tx3 := model.Transaction{
		TransactionId: "3",
		Account:       account2,
		Timestamp:     now.AddDate(0, 0, 2),
	}

	err := repository.SaveTransaction(tx1)
	assert.NoError(t, err)

	err = repository.SaveTransaction(tx2)
	assert.NoError(t, err)

	err = repository.SaveTransaction(tx3)
	assert.NoError(t, err)

	// account1 should have 2 transactions
	transactions, err := repository.GetTransactions(account1.AccountId)
	assert.NoError(t, err)
	assert.Equal(t, model.Transactions{tx2, tx1}, transactions)

	// account2 should have 1 transaction
	transactions, err = repository.GetTransactions(account2.AccountId)
	assert.NoError(t, err)
	assert.Equal(t, model.Transactions{tx3}, transactions)

}
