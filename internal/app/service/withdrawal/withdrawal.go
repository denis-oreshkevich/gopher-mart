package withdrawal

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository"
)

var ErrInvalidSum = errors.New("sum in order is negative")

var ErrInsufficientFunds = errors.New("insufficient funds")

var ErrDuplicateOrder = errors.New("duplicate order")

type Service struct {
	wRepo     withdrawal.Repository
	orderRepo order.Repository
	balRepo   balance.Repository
	tr        repository.Transactor
}

func NewService(wRepo withdrawal.Repository, orderRepo order.Repository,
	balRepo balance.Repository, tr repository.Transactor) *Service {
	return &Service{
		wRepo:     wRepo,
		orderRepo: orderRepo,
		balRepo:   balRepo,
		tr:        tr,
	}
}

func (s *Service) Withdraw(ctx context.Context, withdraw withdrawal.Withdrawal) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetUserID: %w", err)
	}
	num := withdraw.Order
	err = s.tr.InTransaction(ctx, func(ctx context.Context) error {
		err = s.orderRepo.CreateOrder(ctx, num, userID)
		if err != nil {
			if errors.Is(err, order.ErrOrderAlreadyExist) {
				return fmt.Errorf("orderRepo.CreateOrder: %w: %w", err, ErrDuplicateOrder)
			}
			return fmt.Errorf("orderRepo.CreateOrder: %w", err)
		}
		sum := withdraw.Sum
		if sum < 0 {
			return ErrInvalidSum
		}

		err = s.balRepo.WithdrawBalanceByUserID(ctx, sum, userID)
		if err != nil {
			if errors.Is(err, balance.ErrCheckConstraint) {
				return fmt.Errorf("%w: %w", ErrInsufficientFunds, err)
			}
			return fmt.Errorf("balRepo.WithdrawBalanceByUserID: %w", err)
		}
		err = s.wRepo.RegisterWithdrawal(ctx, withdraw)
		if err != nil {
			return fmt.Errorf("repo.Register: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("tr.InTransaction: %w", err)
	}

	return nil
}

func (s *Service) FindUserWithdrawals(ctx context.Context) (withdrawal.Withdrawals, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	wdraw, err := s.wRepo.FindWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("wRepo.FindWithdrawalsByUserID: %w", err)
	}
	var ws withdrawal.Withdrawals = wdraw
	return ws, nil
}
