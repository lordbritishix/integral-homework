package syncers

import (
	"github.com/google/uuid"
	"github.com/lordbritishix/integral/internal/model"
	"go.uber.org/zap"
	"sync"
	"time"
)

type SyncerManager struct {
	syncers map[string]Syncer
	syncCh  chan model.SyncJob
	log     *zap.Logger
}

func NewSyncManager(syncers map[string]Syncer, log *zap.Logger) SyncerManager {
	return SyncerManager{
		syncers: syncers,
		syncCh:  make(chan model.SyncJob, 100),
		log:     log,
	}
}

func (e *SyncerManager) Start() error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	defer wg.Done()
	defer close(e.syncCh)

	for {
		select {
		case data := <-e.syncCh:
			{
				// todo: increase durability when encountering transient errors
				err := e.runSyncWallet(data.Account)
				if err != nil {
					e.log.Error("failed to sync wallet", zap.Error(err))
					continue
				}
			}
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (e *SyncerManager) SubmitSyncJob(account model.Account) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	e.syncCh <- model.SyncJob{
		Id:      id.String(),
		Account: account,
	}
	return nil
}

func (e *SyncerManager) runSyncWallet(account model.Account) error {
	e.log.Info("Start syncing account", zap.String("account", account.AccountId))

	syncer, ok := e.syncers[account.Wallet.Network]

	if !ok {
		e.log.Warn("no syncer capable of syncing the account", zap.String("account", account.AccountId))
		return nil
	}

	err := syncer.SyncAccount(account)
	if err != nil {
		e.log.Error("failed to sync account", zap.String("account", account.AccountId), zap.Error(err))
		return err
	}

	return nil
}
