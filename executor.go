package data

import "context"

type ExecutorManager[T any] interface {
	Executor() T
	TxExecutor(context.Context) (TxExecutor[T], error)
}

type TxExecutor[T any] interface {
	Executor() T
	Commit() error
	Rollback() error
}
