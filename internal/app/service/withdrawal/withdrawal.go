package withdrawal

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	balsvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/service/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/util/auth"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrInvalidSum = errors.New("sum in order is negative")

var ErrInsufficientFunds = errors.New("insufficient funds")

type Service struct {
	repo       withdrawal.Repository
	orderSvc   *order.Service
	balanceSvc *balsvc.Service
}

func NewService(repo withdrawal.Repository, orderSvc *order.Service,
	balanceSvc *balsvc.Service) *Service {
	return &Service{
		repo:       repo,
		orderSvc:   orderSvc,
		balanceSvc: balanceSvc,
	}
}

func (s *Service) Withdraw(ctx context.Context, withdraw withdrawal.Withdrawal) error {
	//TODO tx
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetUserID: %w", err)
	}
	num := withdraw.Order
	exists, err := s.orderSvc.CheckIsExist(ctx, num, userID)
	if err != nil {
		return fmt.Errorf("orderSvc.CheckIsExist: %w", err)
	}
	if !exists {
		return ErrOrderNotFound
	}
	sum := withdraw.Sum
	if sum < 0 {
		return ErrInvalidSum
	}
	err = s.balanceSvc.WithdrawByUserID(ctx, sum, userID)
	if err != nil {
		if errors.Is(err, balance.ErrCheckConstraint) {
			return fmt.Errorf("%w: %w", ErrInsufficientFunds, err)
		}
		return fmt.Errorf("balanceSvc.WithdrawByUserID: %w", err)
	}
	err = s.repo.Register(ctx, withdraw)
	if err != nil {
		return fmt.Errorf("repo.Register: %w", err)
	}
	return nil
}

func (s *Service) FindUserWithdrawals(ctx context.Context) (withdrawal.Withdrawals, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	wdraw, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo.FindByUserID: %w", err)
	}
	var ws withdrawal.Withdrawals = wdraw
	return ws, nil
}
