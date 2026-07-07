package repository

import "context"

type TxFunc func(repoProvider interface{}) error

type Transaction interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
