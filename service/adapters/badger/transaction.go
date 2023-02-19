package badger

import (
	"github.com/boreq/errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/planetary-social/scuttlego-pub/service/app/commands"
)

type CommandsAdaptersFactory func(tx *badger.Txn) (commands.Adapters, error)

type CommandsTransactionProvider struct {
	db      *badger.DB
	factory CommandsAdaptersFactory
}

func NewCommandsTransactionProvider(db *badger.DB, factory CommandsAdaptersFactory) *CommandsTransactionProvider {
	return &CommandsTransactionProvider{db: db, factory: factory}
}

func (t CommandsTransactionProvider) Transact(f func(adapters commands.Adapters) error) error {
	return t.db.Update(func(tx *badger.Txn) error {
		adapters, err := t.factory(tx)
		if err != nil {
			return errors.Wrap(err, "failed to build adapters")
		}

		return f(adapters)
	})
}
