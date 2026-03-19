package datastore

import (
	"context"

	"github.com/jackc/pgx/v5"

	storyW "codebase/internal/models"
)

type PgxDatastore struct {
	Pool        PGXPool
	bobExecutor BobExecutor
}

var _ storyW.Datastore = (*PgxDatastore)(nil)

func NewPgxDatastore(pool PGXPool) (*PgxDatastore, error) {
	return &PgxDatastore{pool, &BobExecutorPgx{pool}}, nil
}

func (ds *PgxDatastore) Begin(ctx context.Context) (storyW.TxDatastore, error) {
	tx, err := ds.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgxTxDatastore{
		PgxDatastore: &PgxDatastore{
			Pool:        tx,
			bobExecutor: &BobExecutorPgx{tx},
		},
		tx: tx,
	}, nil
}

type PgxTxDatastore struct {
	*PgxDatastore
	tx pgx.Tx
}

func (ds *PgxTxDatastore) Commit(ctx context.Context) error {
	return ds.tx.Commit(ctx)
}

func (ds *PgxTxDatastore) Rollback(ctx context.Context) error {
	return ds.tx.Rollback(ctx)
}
