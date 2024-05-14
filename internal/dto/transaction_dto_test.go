package dto

import (
	"github.com/lordbritishix/integral/internal/model"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestNewTransactionDto(t *testing.T) {
	value := new(big.Int)
	value.SetInt64(111498504)

	transaction := model.Transaction{
		TransactionId: "1",
		Account: model.Account{
			AccountId: "1",
		},
		ToAddress:   "to",
		FromAddress: "from",
		Type:        model.WithdrawalType,
		Amount:      *value,
		Decimal:     6,
		Symbol:      "USDC",
		Timestamp:   time.Unix(1715664692, 0), //Tuesday, May 14, 2024 5:31:32 AM GMT
		TxnHash:     "abc",
	}

	transactionDto := NewTransactionDto(transaction)

	assert.Equal(t, transaction.TransactionId, transactionDto.TransactionId)
	assert.Equal(t, transaction.ToAddress, transactionDto.ToAddress)
	assert.Equal(t, transaction.FromAddress, transactionDto.FromAddress)
	assert.Equal(t, transaction.Type, transactionDto.Type)
	assert.Equal(t, "111.498504", transactionDto.Amount)
	assert.Equal(t, transaction.Decimal, transactionDto.Decimal)
	assert.Equal(t, "2024-05-14T05:31:32Z", transactionDto.Timestamp)
	assert.Equal(t, transaction.TxnHash, transactionDto.TxnHash)
}
