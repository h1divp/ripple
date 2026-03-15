package user

import (
	"context"
	// "database/sql"
	"fmt"
	"time"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	// "github.com/stephenafamo/bob"
)

type Repository struct {
	logger zerolog.Logger
	// db     *bob.DB
	redis *redis.Client
}

// func NewRepository(logger zerolog.Logger, db *bob.Executor, rdb *redis.Client) *repository {
func NewRepository(logger zerolog.Logger, rdb *redis.Client) *Repository {
	return &Repository{
		logger: logger.With().Str("repository", "user").Logger(),
		// db:     db,
		redis: rdb,
	}
}

// TODO: make sure we limit the amount of sessions one user can make. Should be done with a combination of ipaddr and some other identifier.
// (this should be done in the service file)
func (r *Repository) SaveSession(ctx context.Context, s *Session) error {
	key := fmt.Sprintf("session:%s", s.ID)

	err := r.redis.HSet(ctx, key, s).Err()
	if err != nil {
		r.logger.Err(err).Msg("Failed to save session to redis")
		return err
	}

	expireTime := config.Load().SessionExpireTime
	err = r.redis.Expire(ctx, key, time.Duration(expireTime)*time.Minute).Err() // TODO: test
	if err != nil {
		r.logger.Err(err).Msg("Failed to set expiriation for session")
		return err
	}

	return nil
}
