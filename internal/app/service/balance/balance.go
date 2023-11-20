package balance

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/util/auth"
)

type Service struct {
	repo balance.Repository
}

func NewService(repo balance.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetUserID: %w", err)
	}
	err = s.repo.Create(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}
	return nil
}

func (s *Service) FindUserBalance(ctx context.Context) (balance.Balance, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("auth.GetUserID: %w", err)
	}
	bal, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("repo.FindByUserID: %w", err)
	}
	return bal, nil
}

func (s *Service) RefillByUserID(ctx context.Context, sum float64, userID string) error {
	return s.repo.RefillByUserID(ctx, sum, userID)
}
func (s *Service) WithdrawByUserID(ctx context.Context, sum float64, userID string) error {
	return s.repo.WithdrawByUserID(ctx, sum, userID)
}
