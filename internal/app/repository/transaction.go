package repository

import "context"

type Transactor interface {
	InTransaction(ctx context.Context, transact func(context.Context) error) error
}
