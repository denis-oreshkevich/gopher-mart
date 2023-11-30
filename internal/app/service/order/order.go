package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
)

var ErrOrderCreatedByAnotherUser = errors.New("order created by another user")

type Service struct {
	repo order.Repository
}

func NewService(repo order.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, orderNum string) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetUserID: %w", err)
	}
	err = s.repo.CreateOrder(ctx, orderNum, userID)
	if err != nil {
		if errors.Is(err, order.ErrOrderAlreadyExist) {
			logger.Log.Debug(fmt.Sprintf("order already exist orderNum = %s", orderNum))
			ord, errF := s.repo.FindOrderByNum(ctx, orderNum)
			if errF != nil {
				return fmt.Errorf("st.FindOrderByNum: %w", err)
			}
			same := userID == ord.UserID
			if !same {
				return ErrOrderCreatedByAnotherUser
			}
			return err
		}
		return fmt.Errorf("repo.Create: %w", err)
	}
	return nil
}

func (s *Service) FindUserOrders(ctx context.Context) (order.Orders, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	orders, err := s.repo.FindOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo.FindOrdersByUserID: %w", err)
	}
	var ord order.Orders = orders
	return ord, nil
}
