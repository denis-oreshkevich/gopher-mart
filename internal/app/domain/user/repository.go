package user

import (
	"context"
	"errors"
)

var ErrUserAlreadyExist = errors.New("user already exist")

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	Create(ctx context.Context, usr User) (User, error)
	FindByLogin(ctx context.Context, login string) (User, error)
}
