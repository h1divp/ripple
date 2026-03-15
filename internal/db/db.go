package store

import (
	"context"
	"database/sql"
	"embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stephenafamo/bob"

	_ "github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Store struct {
	db   bob.DB
	pool *pgxpool.Pool
}

func New(connString string) (*Store, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	if err := sqlDB.Ping(); err != nil {
		pool.Close()
		return nil, err
	}

	executor := bob.NewDB(sqlDB)

	// Run migrations
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return nil, err
	}

	return &Store{
		db:   executor,
		pool: pool,
	}, nil
}

func (s *Store) Close() error {
	s.pool.Close()
	return nil
}

func (s *Store) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return s.pool.BeginTx(ctx, opts)
}

func (s *Store) executor(tx *sql.Tx) bob.Executor {
	if tx != nil {
		return bob.NewTx(tx)
	}

	return s.db
}
