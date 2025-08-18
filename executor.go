package data

type ExecutorManager[T any] interface {
	Executor() T
}

type TxExecutor[T any] interface {
	Executor() T
	Commit() error
	Rollback() error
}
