package data

import (
	"context"
	"database/sql"
)

type SqlExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SqlExecutorManager struct {
	db *sql.DB
}

func (s *SqlExecutorManager) Executor() SqlExecutor {
	return s.db
}

func (s *SqlExecutorManager) TxExecutor(ctx context.Context) (TxExecutor[SqlExecutor], error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &TxSqlExecutor{tx: tx}, nil
}

type TxSqlExecutor struct {
	tx *sql.Tx
}

func (t *TxSqlExecutor) Commit() error {
	return t.tx.Commit()
}

func (t *TxSqlExecutor) Executor() SqlExecutor {
	return t.tx
}

func (t *TxSqlExecutor) Rollback() error {
	return t.tx.Rollback()
}

func NewSqlExecutorManager(db *sql.DB) ExecutorManager[SqlExecutor] {
	return &SqlExecutorManager{
		db: db,
	}
}
