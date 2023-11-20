package order

import (
	"context"
	"errors"
)

var ErrOrderAlreadyExist = errors.New("order already exist")

type Repository interface {
	Create(ctx context.Context, orderNum, userID string) error
	FindByNum(ctx context.Context, orderNum string) (Order, error)
	FindByUserID(ctx context.Context, userID string) ([]Order, error)
	CheckIsExist(ctx context.Context, orderNum, userID string) (bool, error)
	StartProcessing(ctx context.Context, limit int) ([]Order, error)
	UpdateStatusByID(ctx context.Context, id string, acc float64, status string) error
}
