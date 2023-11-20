package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/denis-oreshkevich/gopher-mart/migration"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
	"sync"
)

var (
	db     *pgxpool.Pool
	pgOnce sync.Once

	dbErr error
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pgOnce.Do(func() {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			dbErr = fmt.Errorf("pgxpool.New: %w", err)
			return
		}
		if err = pool.Ping(ctx); err != nil {
			dbErr = fmt.Errorf("pool.Ping: %w", err)
			return
		}
		if err = applyMigration(dsn, migration.SQLFiles); err != nil {
			dbErr = fmt.Errorf("applyMigration: %w", err)
			return
		}
		db = pool
	})
	return db, dbErr
}

func applyMigration(dsn string, fsys fs.FS) error {
	//TODO ask about conv between pool
	//db := stdlib.OpenDBFromPool(db, nil)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(fsys)
	goose.SetSequential(true)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}
	return nil
}
