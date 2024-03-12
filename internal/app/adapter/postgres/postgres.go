package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"github.com/denis-oreshkevich/gophkeeper/migration"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
	"sync"
)

var _ service.Repository = (*Repository)(nil)

var (
	db     *pgxpool.Pool
	pgOnce sync.Once

	dbErr error
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(ctx context.Context, dsn string) (*Repository, error) {
	pool, err := initPool(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("initPool: %w", err)
	}
	return &Repository{
		db: pool,
	}, nil
}

func (r *Repository) InTransaction(ctx context.Context,
	transact func(context.Context) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tx begin. %w", err)
	}
	defer tx.Rollback(ctx)
	err = transact(ctx)
	if err != nil {
		return fmt.Errorf("transact: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *Repository) Close() {
	r.db.Close()
}

func initPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
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
