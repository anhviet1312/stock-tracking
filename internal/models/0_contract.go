package models

import (
	"context"

	"github.com/go-redis/redis_rate/v10"

	"codebase/internal/bob"
)

type ReadOnlyDatastore interface {
	FindStockBySymbol(ctx context.Context, symbol string) (*Stock, error)
}

type Datastore interface {
	ReadOnlyDatastore
	Begin(ctx context.Context) (TxDatastore, error)

	CreateManyStocks(ctx context.Context, params []*bob.StockSetter) error
	UpdateStock(ctx context.Context, symbol string, param *bob.StockSetter) (*Stock, error)
}

// TxDatastore a transactional Datastore. All functions are execute under a transaction.
// Call to functions after Commit or Rollback will return an error.
type TxDatastore interface {
	Datastore
	// Commit commits the transaction. Commit is safe to call multiple times. If the commit fails with a rollback status (e.g. the transaction was already
	// in a broken state) then an error where errors.Is(ErrTxCommitRollback) is true will be returned.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction. Rollback is safe to call multiple times. Hence, a defer tx.Rollback() is safe even if tx.Commit() will
	// be called first in a non-error condition. Any other failure of a real transaction will result in the datastore being closed.
	Rollback(ctx context.Context) error
}

type Limiter interface {
	Allow(ctx context.Context, key string, limit redis_rate.Limit) error
}

type Iterator[T any] interface {
	Next() bool
	Get() (T, error)
	Close() error
	Err() error
}
