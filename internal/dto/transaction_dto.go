package dto

import (
	"github.com/lordbritishix/integral/internal/model"
	"math"
	"math/big"
	"time"
)

type GetTransactionsResponse struct {
	Data  []TransactionDto `json:"data"`
	Count int              `json:"count"`
}

func NewGetTransactionsResponse(transactions []model.Transaction) GetTransactionsResponse {
	var converted []TransactionDto

	for _, t := range transactions {
		converted = append(converted, NewTransactionDto(t))
	}

	return GetTransactionsResponse{
		Data:  converted,
		Count: len(transactions),
	}
}

type TransactionDto struct {
	TransactionId string `json:"id"`
	AccountId     string `json:"accountId"`
	ToAddress     string `json:"toAddress"`
	FromAddress   string `json:"fromAddress"`
	Type          string `json:"type"`
	Amount        string `json:"amount"`
	Symbol        string `json:"symbol"`
	Decimal       int    `json:"decimal"`
	Timestamp     string `json:"timestamp"`
	TxnHash       string `json:"txnHash"`
	IsSpam        bool   `json:"isSpam"`
}

func NewTransactionDto(transaction model.Transaction) TransactionDto {
	dividend := new(big.Int)
	dividend.SetInt64(int64(math.Pow10(transaction.Decimal)))
	result := new(big.Float)
	result.Quo(new(big.Float).SetInt(&transaction.Amount), new(big.Float).SetInt(dividend))

	return TransactionDto{
		TransactionId: transaction.TransactionId,
		AccountId:     transaction.Account.AccountId,
		ToAddress:     transaction.ToAddress,
		FromAddress:   transaction.FromAddress,
		Type:          transaction.Type,
		Amount:        result.String(),
		Symbol:        transaction.Symbol,
		Decimal:       transaction.Decimal,
		Timestamp:     transaction.Timestamp.UTC().Format(time.RFC3339),
		TxnHash:       transaction.TxnHash,
		IsSpam:        transaction.IsSpam,
	}
}
