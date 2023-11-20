package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/user"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

var _ user.Repository = (*UserRepository)(nil)

func (s *UserRepository) Create(ctx context.Context, usr user.User) (user.User, error) {
	query := `insert into mart.usr(login, password) values (@login, @password) returning usr.id`
	args := pgx.NamedArgs{
		"login":    usr.Login,
		"password": usr.HashedPassword,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&usr.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.UniqueViolation == pgErr.Code {
				return usr, fmt.Errorf("row.Scan, contraint %s: %w",
					pgErr.ConstraintName, user.ErrUserAlreadyExist)
			}
		}
		return usr, fmt.Errorf("row.Scan: %w", err)
	}
	return usr, nil
}

func (s *UserRepository) FindByLogin(ctx context.Context, login string) (user.User, error) {
	query := `select id, login, password from mart.usr where login=@login`
	args := pgx.NamedArgs{
		"login": login,
	}
	row := s.db.QueryRow(ctx, query, args)
	var usr user.User
	err := row.Scan(&usr.ID, &usr.Login, &usr.HashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.ErrUserNotFound
		}
		return user.User{}, fmt.Errorf("row.Scan: %w", err)
	}
	return usr, nil
}
