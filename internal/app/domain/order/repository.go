package order

import (
	"context"
	"errors"
)

var ErrOrderAlreadyExist = errors.New("order already exist")

type Repository interface {
	CreateOrder(ctx context.Context, orderNum, userID string) error
	FindOrderByNum(ctx context.Context, orderNum string) (Order, error)
	FindOrdersByUserID(ctx context.Context, userID string) ([]Order, error)
	CheckIsOrderExist(ctx context.Context, orderNum, userID string) (bool, error)
	StartOrderProcessing(ctx context.Context, limit int) ([]Order, error)
	UpdateOrderStatusByID(ctx context.Context, id string, acc float64, status string) error
}
