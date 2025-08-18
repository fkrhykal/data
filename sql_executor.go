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

type SqlTxManager struct {
	db *sql.DB
}

func (m *SqlTxManager) Executor() SqlExecutor {
	return m.db
}

func (m *SqlTxManager) TxExecutor(ctx context.Context) (TxExecutor[SqlExecutor], error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &sqlTxExecutor{tx: tx}, nil
}

func NewSqlExecutorManager(db *sql.DB) ExecutorManager[SqlExecutor] {
	return &SqlTxManager{db: db}
}

type sqlTxExecutor struct {
	tx *sql.Tx
}

func (t *sqlTxExecutor) Executor() SqlExecutor {
	return t.tx
}

func (t *sqlTxExecutor) Rollback() error {
	return t.tx.Rollback()
}

func (t *sqlTxExecutor) Commit() error {
	return t.tx.Commit()
}
