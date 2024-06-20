package transactor

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type QueryEngine interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx OR pool
}

type TransactionManager struct {
	Pool *pgxpool.Pool
}

const key = "tx"

func (tm *TransactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	
	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		return tx.Rollback(ctx)
	}

	if err := tx.Commit(ctx); err != nil {
		return tx.Rollback(ctx)
	}

	return nil
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.Pool
}
