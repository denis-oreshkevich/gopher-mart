package balance

import (
	"context"
	"errors"
)

var ErrCheckConstraint = errors.New("check constraint")

type Repository interface {
	CreateBalance(ctx context.Context, userID string) error
	FindBalanceByUserID(ctx context.Context, userID string) (Balance, error)
	RefillBalanceByUserID(ctx context.Context, sum float64, userID string) error
	WithdrawBalanceByUserID(ctx context.Context, sum float64, userID string) error
}
