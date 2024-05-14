package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	config2 "github.com/lordbritishix/integral/internal/config"
	txcontroller "github.com/lordbritishix/integral/internal/controller"
	"github.com/lordbritishix/integral/internal/model"
	"github.com/lordbritishix/integral/internal/repository"
	"github.com/lordbritishix/integral/internal/syncers"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	// setup logger
	logger, _ := zap.NewProduction()
	logger.Info("Running api server")

	// setup config
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Not loading env file")
	}

	config := config2.NewConfig()

	// setup repositories
	accountRepository := repository.NewInMemoryAccountRepository()
	txnRepository := repository.NewInMemoryTransactionRepository()

	// initialize syncers
	ethSyncer := syncers.NewEthereumSyncer(txnRepository, logger, config)
	syncerList := map[string]syncers.Syncer{}
	syncerList[model.EthereumNetwork] = ethSyncer

	syncManager := syncers.NewSyncManager(syncerList, logger)

	go func() {
		err := syncManager.Start()
		if err != nil {
			logger.Error("Unable to start eth syncer", zap.Error(err))
			return
		}
	}()

	// hard-code accounts for now
	err = accountRepository.SaveAccount(model.Account{
		AccountId:   "account_abc",
		AccountName: "Jim Quitevis",
		Wallet: model.Wallet{
			WalletId:   "wallet_abc",
			WalletName: "Jim's Eth Wallet",
			Address:    "0x912fd21d7a69678227fe6d08c64222db41477ba0",
			Network:    model.EthereumNetwork,
		},
	})

	if err != nil {
		panic("Unable to create a stub account")
	}

	// setup http server
	r := mux.NewRouter()

	controller := txcontroller.NewTransactionsController(txnRepository, accountRepository, syncManager)

	// Creates an account with the wallet id if it does not exist - and starts the wallet syncing process
	r.HandleFunc("/accounts/{accountId}/sync", controller.PostSyncHandler).Methods("POST")

	// Returns the sync'd transactions for the particular account id
	r.HandleFunc("/accounts/{accountId}/transactions", controller.GetTransactionsHandler).Methods("GET")

	// Returns the "total pooled ETH" and the "total shares" from the stETH
	r.HandleFunc("/steth/stats", controller.GetStethStats).Methods("GET")

	// Returns the last 5 addresses that have deposited into the stETH pool.
	r.HandleFunc("/steth/last-deposits", controller.GetStethDeposits).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Unable to start server", zap.Error(err))
	}
}
