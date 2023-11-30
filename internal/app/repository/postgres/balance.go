package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/balance"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ balance.Repository = (*Repository)(nil)

func (r *Repository) CreateBalance(ctx context.Context, userID string) error {
	query := "insert into mart.balance (user_id) values (@user_id) on conflict do nothing"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (r *Repository) FindBalanceByUserID(ctx context.Context, userID string) (balance.Balance, error) {
	query := "select cur, withdrawn, user_id from mart.balance where user_id=@user_id"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	row := r.db.QueryRow(ctx, query, args)
	var bal balance.Balance
	err := row.Scan(&bal.Current, &bal.Withdrawn, &bal.UserID)
	if err != nil {
		return balance.Balance{}, fmt.Errorf("row.Scan: %w", err)
	}
	return bal, nil
}

func (r *Repository) RefillBalanceByUserID(ctx context.Context, sum float64, userID string) error {
	query := "update mart.balance set cur = cur + @amount where user_id=@user_id"
	args := pgx.NamedArgs{
		"user_id": userID,
		"amount":  sum,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (r *Repository) WithdrawBalanceByUserID(ctx context.Context, sum float64, userID string) error {
	query := "update mart.balance set cur = cur - @amount, " +
		"withdrawn = withdrawn + @amount  where user_id=@user_id"
	args := pgx.NamedArgs{
		"user_id": userID,
		"amount":  sum,
	}
	_, err := r.db.Exec(ctx, query, args)
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
