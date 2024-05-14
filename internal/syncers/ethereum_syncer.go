package syncers

import (
	"errors"
	"github.com/lordbritishix/integral/internal/config"
	"github.com/lordbritishix/integral/internal/dto/etherscan"
	"github.com/lordbritishix/integral/internal/model"
	"github.com/lordbritishix/integral/internal/providers"
	"github.com/lordbritishix/integral/internal/repository"
	"go.uber.org/zap"
	"math/big"
	"strconv"
	"time"
)

type EthereumSyncer struct {
	repository *repository.TransactionsRepository
	log        *zap.Logger
	etherscan  *providers.EtherscanClient
}

func NewEthereumSyncer(repository repository.TransactionsRepository, log *zap.Logger, config *config.Config) *EthereumSyncer {
	return &EthereumSyncer{
		repository: &repository,
		log:        log,
		etherscan:  providers.NewEtherscanClient(config),
	}
}

func (e *EthereumSyncer) SyncAccount(account model.Account) error {
	// todo: start block and end block should be looped, and start from the beginning, record watermarks
	startBlock := 18000000
	endBlock := 18100000

	// get ETH transfers
	ethTxs, err := e.processEthTransactions(account, startBlock, endBlock)
	if err != nil {
		return err
	}

	// get token transfers
	tokenTxs, err := e.processTokenTransactions(account, startBlock, endBlock)
	if err != nil {
		return err
	}

	e.log.Info("Synced ethereum wallet", zap.Int("count", len(ethTxs)+len(tokenTxs)))

	return nil
}

func (e *EthereumSyncer) processEthTransactions(account model.Account, startBlock int, endBlock int) ([]etherscan.Transaction, error) {
	transactions, err := e.etherscan.GetEthTransactions(account.Wallet.Address, startBlock, endBlock, 100)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		timestampUnix, err := strconv.ParseInt(transaction.TimeStamp, 10, 64)
		if err != nil {
			return nil, err
		}

		value := new(big.Int)
		_, ok := value.SetString(transaction.Value, 10)
		if !ok {
			return nil, errors.New("invalid value")
		}

		tx := model.Transaction{
			TransactionId: transaction.Hash,
			Account:       account,
			ToAddress:     transaction.To,
			FromAddress:   transaction.From,
			Amount:        *value,
			Decimal:       18,
			Symbol:        "ETH",
			Timestamp:     time.Unix(timestampUnix, 0),
			TxnHash:       transaction.Hash,
			Type:          getDepositType(transaction, account.Wallet.Address),
			IsSpam:        false,
		}

		// todo: batch save instead
		err = (*e.repository).SaveTransaction(tx)
		if err != nil {
			return nil, err
		}
	}

	// todo: also create a transaction for gas fees if from address is same as passed address?
	return transactions, nil
}

func (e *EthereumSyncer) processTokenTransactions(account model.Account, startBlock int, endBlock int) ([]etherscan.Transaction, error) {
	transactions, err := e.etherscan.GetTokenTransactions(account.Wallet.Address, startBlock, endBlock, 100)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		timestampUnix, err := strconv.ParseInt(transaction.TimeStamp, 10, 64)
		if err != nil {
			return nil, err
		}

		value := new(big.Int)
		_, ok := value.SetString(transaction.Value, 10)
		if !ok {
			return nil, errors.New("invalid value")
		}

		decimal, err := strconv.ParseInt(transaction.TimeStamp, 10, 32)
		if err != nil {
			return nil, err
		}

		isSpam := false

		if transaction.ContractAddress != nil {
			// todo: needs to be cached cause this will hit rate limiting hard
			isSpam, err = e.etherscan.IsTokenSpam(*transaction.ContractAddress)

			if err != nil {
				return nil, err
			}
		}

		tx := model.Transaction{
			TransactionId: transaction.Hash,
			Account:       account,
			ToAddress:     transaction.To,
			FromAddress:   transaction.From,
			Amount:        *value,
			Decimal:       int(decimal),
			Symbol:        *transaction.TokenSymbol,
			Timestamp:     time.Unix(timestampUnix, 0),
			TxnHash:       transaction.Hash,
			Type:          getDepositType(transaction, account.Wallet.Address),
			IsSpam:        isSpam,
		}

		// todo: batch save instead
		err = (*e.repository).SaveTransaction(tx)
		if err != nil {
			return nil, err
		}
	}

	// todo: also create a transaction for gas fees if from address is same as passed address?
	return transactions, nil
}

func getDepositType(transaction etherscan.Transaction, walletAddress string) string {
	depositType := model.WithdrawalType
	if transaction.From == walletAddress {
		depositType = model.WithdrawalType
	} else if transaction.To == walletAddress {
		depositType = model.DepositType
	}

	return depositType
}

func (e *EthereumSyncer) SyncerName() string {
	return "EthereumSyncer"
}
