package postgres

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithdrawalRepository struct {
	db *pgxpool.Pool
}

func NewWithdrawalRepository(db *pgxpool.Pool) *WithdrawalRepository {
	return &WithdrawalRepository{
		db: db,
	}
}

var _ withdrawal.Repository = (*WithdrawalRepository)(nil)

func (s *WithdrawalRepository) Register(ctx context.Context, withdraw withdrawal.Withdrawal) error {
	query := "insert into mart.withdrawal (amount ,order_id) values " +
		"(@amount, (select ordr.id from mart.ordr where ordr.num = @order_num))"
	args := pgx.NamedArgs{
		"amount":    withdraw.Sum,
		"order_num": withdraw.Order,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (s *WithdrawalRepository) FindByUserID(ctx context.Context,
	userID string) ([]withdrawal.Withdrawal, error) {
	query := "select amount, processed_at, mo.num from mart.withdrawal as mw" +
		"inner join mart.ordr as mo on mw.order_id = mo.id " +
		"where mo.user_id=@user_id order by processed_at asc"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := s.db.Query(ctx, query, args)
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
