package postgres

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	"github.com/jackc/pgx/v5"
)

var _ withdrawal.Repository = (*Repository)(nil)

func (r *Repository) RegisterWithdrawal(ctx context.Context, withdraw withdrawal.Withdrawal) error {
	query := "insert into mart.withdrawal (amount ,order_id) values " +
		"(@amount, (select ordr.id from mart.ordr where ordr.num = @order_num))"
	args := pgx.NamedArgs{
		"amount":    withdraw.Sum,
		"order_num": withdraw.Order,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (r *Repository) FindWithdrawalsByUserID(ctx context.Context,
	userID string) ([]withdrawal.Withdrawal, error) {
	query := "select amount, processed_at, mo.num from mart.withdrawal as mw " +
		"inner join mart.ordr as mo on mw.order_id = mo.id " +
		"where mo.user_id=@user_id order by processed_at asc"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]withdrawal.Withdrawal, 0)
	for rows.Next() {
		var w withdrawal.Withdrawal
		err := rows.Scan(&w.Sum, &w.ProcessedAt, &w.Order)
		if err != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		res = append(res, w)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}
