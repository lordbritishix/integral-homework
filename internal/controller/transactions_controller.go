package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lordbritishix/integral/internal/dto"
	"github.com/lordbritishix/integral/internal/model"
	"github.com/lordbritishix/integral/internal/repository"
	"github.com/lordbritishix/integral/internal/syncers"
	"net/http"
)

type TransactionsController struct {
	txnRepository      repository.TransactionsRepository
	accountsRepository repository.AccountRepository
	syncManager        syncers.SyncerManager
}

func NewTransactionsController(txnRepository repository.TransactionsRepository, accountRepository repository.AccountRepository, syncerManager syncers.SyncerManager) *TransactionsController {
	return &TransactionsController{
		txnRepository:      txnRepository,
		accountsRepository: accountRepository,
		syncManager:        syncerManager,
	}
}

func (n *TransactionsController) PostSyncHandler(writer http.ResponseWriter, request *http.Request) {
	accountId := mux.Vars(request)["accountId"]

	account, err := n.accountsRepository.GetAccount(accountId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if account == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	err = n.syncManager.SubmitSyncJob(*account)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	return
}

func (n *TransactionsController) GetTransactionsHandler(writer http.ResponseWriter, request *http.Request) {
	accountId := mux.Vars(request)["accountId"]

	if len(accountId) <= 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := n.accountsRepository.Exists(accountId)

	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactions, err := n.txnRepository.GetTransactions(accountId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if transactions == nil {
		transactions = []model.Transaction{}
	}

	response := dto.NewGetTransactionsResponse(transactions)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(jsonResponse)
	}
}

func (n *TransactionsController) GetStethDeposits(writer http.ResponseWriter, request *http.Request) {

}

func (n *TransactionsController) GetStethStats(writer http.ResponseWriter, request *http.Request) {

}
