package store

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/stephenafamo/bob"

	_ "github.com/jackc/pgx/v5"
)

var embedMigrations embed.FS

type Store struct {
	db bob.DB
}

func New(dsn string) (*Store, error) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return nil, err
	}

	return &Store{db: bob.NewDB(sqlDB)}, nil
}

func (s *Store) Close() error {
	return s.db.DB.Close()
}

func (s *Store) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.DB.BeginTx(ctx, opts)
}

func (s *Store) executor(tx *sql.Tx) bob.Executor {
	if tx != nil {
		return bob.NewTx(tx)
	}

	return s.db
}
