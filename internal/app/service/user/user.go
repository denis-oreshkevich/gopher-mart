package user

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/user"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/service/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/util/auth"
)

type Service struct {
	repo   user.Repository
	balSvc *balance.Service
}

func NewService(repo user.Repository, balSvc *balance.Service) *Service {
	return &Service{
		repo:   repo,
		balSvc: balSvc,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) (user.User, error) {
	hp, err := auth.EncryptPassword(password)
	if err != nil {
		return user.User{}, fmt.Errorf("util.EncryptPassword: %w", err)
	}
	usr := user.New(login, hp)
	usr, err = s.repo.Create(ctx, usr)
	if err != nil {
		return user.User{}, fmt.Errorf("repo.Create: %w", err)
	}
	err = s.balSvc.Create(ctx, usr.ID)
	if err != nil {
		return user.User{}, fmt.Errorf("balSvc.Create: %w", err)
	}
	return usr, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	us, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("repo.FindByLogin: %w", err)
	}
	logger.Log.Debug(fmt.Sprintf("1 = %s, 2=%s", us.HashedPassword, password))
	err = auth.ComparePasswords(us.HashedPassword, password)
	if err != nil {
		logger.Log.Debug("incorrect passwords")
		return "", err
	}
	token, err := auth.GenerateToken(us.ID)
	if err != nil {
		return "", fmt.Errorf("auth.GenerateToken: %w", err)
	}
	return token, nil
}
