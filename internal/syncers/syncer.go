package syncers

import "github.com/lordbritishix/integral/internal/model"

type Syncer interface {
	SyncAccount(account model.Account) error
	SyncerName() string
}
