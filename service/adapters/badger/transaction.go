package badger

import (
	"github.com/boreq/errors"
	"github.com/dgraph-io/badger/v3"
)

type AdaptersFactory[T any] func(tx *badger.Txn) (T, error)

type TransactionProvider[T any] struct {
	db      *badger.DB
	factory AdaptersFactory[T]
}

func NewTransactionProvider[T any](db *badger.DB, factory AdaptersFactory[T]) *TransactionProvider[T] {
	return &TransactionProvider[T]{db: db, factory: factory}
}

func (t TransactionProvider[T]) Update(f func(adapters T) error) error {
	return t.db.Update(func(tx *badger.Txn) error {
		adapters, err := t.factory(tx)
		if err != nil {
			return errors.Wrap(err, "failed to build adapters")
		}
		return f(adapters)
	})
}

func (t TransactionProvider[T]) View(f func(adapters T) error) error {
	return t.db.View(func(tx *badger.Txn) error {
		adapters, err := t.factory(tx)
		if err != nil {
			return errors.Wrap(err, "failed to build adapters")
		}
		return f(adapters)
	})
}
