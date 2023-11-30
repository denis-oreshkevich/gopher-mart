package user

import (
	"context"
	"errors"
)

var ErrUserAlreadyExist = errors.New("user already exist")

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	CreateUser(ctx context.Context, usr User) (User, error)
	FindUserByLogin(ctx context.Context, login string) (User, error)
}
