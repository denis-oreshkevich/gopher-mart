package user

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/auth"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/user"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/repository"
)

type Service struct {
	userRepo user.Repository
	balRepo  balance.Repository
	tr       repository.Transactor
}

func NewService(userRepo user.Repository,
	balRepo balance.Repository, tr repository.Transactor) *Service {
	return &Service{
		userRepo: userRepo,
		balRepo:  balRepo,
		tr:       tr,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) (user.User, error) {
	hp, err := auth.EncryptPassword(password)
	if err != nil {
		return user.User{}, fmt.Errorf("util.EncryptPassword: %w", err)
	}
	usr := user.New(login, hp)
	err = s.tr.InTransaction(ctx, func(ctx context.Context) error {
		u, err := s.userRepo.CreateUser(ctx, usr)
		if err != nil {
			return fmt.Errorf("userRepo.CreateUser: %w", err)
		}
		err = s.balRepo.CreateBalance(ctx, u.ID)
		if err != nil {
			return fmt.Errorf("balRepo.CreateBalance: %w", err)
		}
		usr = u
		return nil
	})
	if err != nil {
		return user.User{}, fmt.Errorf("tr.InTransaction: %w", err)
	}
	return usr, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	us, err := s.userRepo.FindUserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("userRepo.FindUserByLogin: %w", err)
	}
	err = auth.ComparePasswords(us.HashedPassword, password)
	if err != nil {
		return "", err
	}
	token, err := auth.GenerateToken(us.ID)
	if err != nil {
		return "", fmt.Errorf("auth.GenerateToken: %w", err)
	}
	return token, nil
}
