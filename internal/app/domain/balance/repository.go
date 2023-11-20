package balance

import (
	"context"
	"errors"
)

var ErrCheckConstraint = errors.New("check constraint")

type Repository interface {
	Create(ctx context.Context, userID string) error
	FindByUserID(ctx context.Context, userID string) (Balance, error)
	RefillByUserID(ctx context.Context, sum float64, userID string) error
	WithdrawByUserID(ctx context.Context, sum float64, userID string) error
}
