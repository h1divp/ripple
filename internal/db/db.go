package store

import (
	"context"
	"database/sql"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stephenafamo/bob"

	"github.com/h1divp/echo-chat-v2/internal/models"

	_ "github.com/jackc/pgx/v5"
)

var embedMigrations embed.FS

type Store struct {
	db   bob.Executor
	pool *pgxpool.Pool
	repo *models.Store
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
		repo: models.New(executor),
	}, nil

}

func (s *Store) Close() error {
	s.pool.Close()
	return nil
}

func (s *Store) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Store, *sql.Tx, error) {
	tx, err := s.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	txExecutor := bob.NewTx(tx)

	return &Store{
		db: txExecutor,
		pool: s.pool,
		repo: models.New(txExecutor)
	}
}

func (s *Store) executor(tx *sql.Tx) bob.Executor {
	if tx != nil {
		return bob.NewTx(tx)
	}

	return s.db
}
