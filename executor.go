package data

import "context"

type ExecutorManager[T any] interface {
	Executor(context.Context) T
	TxExecutor(context.Context) TxExecutor[T]
}

type TxExecutor[T any] interface {
	Executor() T
	Commit() error
	Rollback() error
}
