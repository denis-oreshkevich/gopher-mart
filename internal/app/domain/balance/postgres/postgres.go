package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepository struct {
	db *pgxpool.Pool
}

func NewBalanceRepository(db *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{
		db: db,
	}
}

var _ balance.Repository = (*BalanceRepository)(nil)

func (s *BalanceRepository) Create(ctx context.Context, userID string) error {
	query := "insert into mart.balance (user_id) values (@user_id) on conflict do nothing"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (s *BalanceRepository) FindByUserID(ctx context.Context, userID string) (balance.Balance, error) {
	query := "select cur, withdrawn, user_id from mart.balance where user_id=@user_id"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	row := s.db.QueryRow(ctx, query, args)
	var bal balance.Balance
	err := row.Scan(&bal.Current, &bal.Withdrawn, &bal.UserID)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("row.Scan: %w", err)
	}
	return bal, nil
}

func (s *BalanceRepository) RefillByUserID(ctx context.Context, sum float64, userID string) error {
	query := "update mart.balance set cur = cur + @amount where user_id=@user_id"
	args := pgx.NamedArgs{
		"user_id": userID,
		"amount":  sum,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (s *BalanceRepository) WithdrawByUserID(ctx context.Context, sum float64, userID string) error {
	query := "update mart.balance set cur = cur - @amount, " +
		"withdrawn = withdrawn + @amount  where user_id=@user_id"
	args := pgx.NamedArgs{
		"user_id": userID,
		"amount":  sum,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.CheckViolation == pgErr.Code {
				return fmt.Errorf("%w: %w", balance.ErrCheckConstraint, err)
			}
		}
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
