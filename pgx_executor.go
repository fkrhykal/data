package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPoolExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type PgxPoolExecutorManager struct {
	pool *pgxpool.Pool
}

func (p *PgxPoolExecutorManager) Executor() PgxPoolExecutor {
	return p.pool
}

func (p *PgxPoolExecutorManager) TxExecutor(ctx context.Context) (TxExecutor[PgxPoolExecutor], error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &TxPgxPoolExecutor{ctx: ctx, tx: tx}, nil
}

type TxPgxPoolExecutor struct {
	ctx context.Context
	tx  pgx.Tx
}

func (t *TxPgxPoolExecutor) Commit() error {
	return t.tx.Commit(t.ctx)
}

func (t *TxPgxPoolExecutor) Executor() PgxPoolExecutor {
	return t.tx
}

func (t *TxPgxPoolExecutor) Rollback() error {
	return t.tx.Rollback(t.ctx)
}

func NewPgxPoolExecutorManager(pool *pgxpool.Pool) ExecutorManager[PgxPoolExecutor] {
	return &PgxPoolExecutorManager{
		pool: pool,
	}
}
