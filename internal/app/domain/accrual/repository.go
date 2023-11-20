package accrual

import (
	"context"
	"errors"
)

var ErrOrderNotRegistered = errors.New("order not registered")

var ErrTooManyRequests = errors.New("too many requests")

type Repository interface {
	FindByOrderNum(ctx context.Context, num string) (Accrual, error)
}
