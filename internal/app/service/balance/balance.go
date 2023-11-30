package balance

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
)

type Service struct {
	repo balance.Repository
}

func NewService(repo balance.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID string) error {
	err := s.repo.CreateBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo.CreateBalance: %w", err)
	}
	return nil
}

func (s *Service) FindUserBalance(ctx context.Context) (balance.Balance, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("auth.GetUserID: %w", err)
	}
	bal, err := s.repo.FindBalanceByUserID(ctx, userID)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("repo.FindBalanceByUserID: %w", err)
	}
	return bal, nil
}
