package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/order"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

var _ order.Repository = (*OrderRepository)(nil)

func (s *OrderRepository) Create(ctx context.Context, orderNum, userID string) error {
	query := "insert into mart.ordr (num, user_id) values (@num, @user_id)"
	args := pgx.NamedArgs{
		"num":     orderNum,
		"user_id": userID,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.UniqueViolation == pgErr.Code {
				return fmt.Errorf("row.Scan, contraint %s: %w: %w",
					pgErr.ConstraintName, order.ErrOrderAlreadyExist, err)
			}
		}
		return fmt.Errorf("row.Scan: %w", err)
	}
	return nil
}

func (s *OrderRepository) FindByNum(ctx context.Context, orderNum string) (order.Order, error) {
	query := "select num, status, accrual, user_id, uploaded_at " +
		"from mart.ordr where num = @num"
	args := pgx.NamedArgs{
		"num": orderNum,
	}
	row := s.db.QueryRow(ctx, query, args)
	var ord order.Order
	nullFloat := sql.NullFloat64{}
	err := row.Scan(&ord.Number, &ord.Status, &nullFloat, &ord.UserID, &ord.UploadedAt)
	if err != nil {
		return order.Order{}, fmt.Errorf("cannot scan value. %w", err)
	}
	if nullFloat.Valid {
		ord.Accrual = nullFloat.Float64

	}
	return ord, nil
}

func (s *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]order.Order, error) {
	query := "select num, status, accrual, uploaded_at " +
		"from mart.ordr where user_id = @user_id order by uploaded_at asc"
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]order.Order, 0)
	for rows.Next() {
		var ord order.Order
		nullFloat := sql.NullFloat64{}
		err := rows.Scan(&ord.Number, &ord.Status, &nullFloat, &ord.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		if nullFloat.Valid {
			ord.Accrual = nullFloat.Float64
		}
		res = append(res, ord)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}

func (s *OrderRepository) CheckIsExist(ctx context.Context,
	orderNum, userID string) (bool, error) {
	query := "select exists(select 1 from mart.ordr " +
		"where ordr.num = @num and user_id = @user_id)"
	args := pgx.NamedArgs{
		"num":     orderNum,
		"user_id": userID,
	}
	row := s.db.QueryRow(ctx, query, args)
	var ex bool
	err := row.Scan(&ex)
	if err != nil {
		return ex, fmt.Errorf("row.Scan: %w", err)
	}
	return ex, nil
}

func (s *OrderRepository) StartProcessing(ctx context.Context, limit int) ([]order.Order, error) {
	query := "update mart.ordr set status = 'PROCESSING' from " +
		"(select id, num, status, user_id, uploaded_at from mart.ordr where status = 'NEW' " +
		"order by uploaded_at asc limit @lim) as sq where ordr.id = sq.id " +
		"returning sq.id, sq.num, sq.user_id, sq.uploaded_at"
	args := pgx.NamedArgs{
		"lim": limit,
	}
	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]order.Order, 0, limit)
	for rows.Next() {
		var ord order.Order
		err := rows.Scan(&ord.ID, &ord.Number, &ord.UserID, &ord.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		res = append(res, ord)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}

func (s *OrderRepository) UpdateStatusByID(ctx context.Context, id, status string) error {
	query := "update mart.ordr set status = @status where id = @id"
	args := pgx.NamedArgs{
		"id":     id,
		"status": status,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
