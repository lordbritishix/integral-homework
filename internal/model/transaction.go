package model

import (
	"math/big"
	"time"
)

type Transaction struct {
	TransactionId string
	Account       Account
	ToAddress     string
	FromAddress   string
	Type          string
	Amount        big.Int
	Decimal       int
	Symbol        string
	Timestamp     time.Time
	TxnHash       string
	IsSpam        bool
}

type Transactions []Transaction

func (t Transactions) Len() int {
	return len(t)
}

func (t Transactions) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Transactions) Less(i, j int) bool {
	return t[i].Timestamp.After(t[j].Timestamp)
}
