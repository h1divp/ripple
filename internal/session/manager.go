package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	ErrSessionAlreadyExists = errors.New("this user already has a session cookie")
)

const (
	userCtxKey = "user_id"
	sessionKey = "session"
)

type Manager struct {
	rdb    *redis.Client
	sc     *securecookie.SecureCookie
	expiry time.Duration
	logger zerolog.Logger
}

func NewManager(logger zerolog.Logger, rdb *redis.Client, hashKey, blockKey []byte) *Manager {
	return &Manager{
		rdb:    rdb,
		sc:     securecookie.New(hashKey, blockKey),
		expiry: 24 * time.Hour,
		logger: logger,
	}
}

func (m *Manager) CreateSession(ctx context.Context) (string, string, error) {
	val, ok := ctx.Value(userCtxKey).(string)
	if ok {
		m.logger.Debug().Str("userID", val).Msg("Session already exists")
		return "", "", ErrSessionAlreadyExists
	}

	sessionID := uuid.New().String()
	userID := uuid.New().String()

	err := m.rdb.Set(ctx, fmt.Sprintf("session:%s", sessionID), userID, m.expiry).Err()
	if err != nil {
		return "", "", err
	}

	signed, err := m.sc.Encode("session_id", sessionID)
	m.logger.Debug().Msg("created session")
	return signed, userID, err
}

func (m *Manager) DeleteSession(ctx context.Context, sessionID string) error {
	err := m.rdb.Del(ctx, fmt.Sprintf("session:%s", sessionID)).Err()
	if err != nil {
		m.logger.Err(err).Msg("Failed to delete session from redis")
		return err
	}
	m.logger.Debug().Msg("Deleted session from redis")
	return nil
}

func (r *Manager) RemoveAllSessions(ctx context.Context) error {
	err := r.rdb.Del(ctx, sessionKey).Err()
	if err != nil {
		r.logger.Err(err).Msg("Failed to remove all sessions during cleanup")
		return err
	}
	r.logger.Info().Msg("Successfully cleaned sessions in Redis")
	return nil
}

// TODO: reorder
func (m *Manager) GetSessionIDAndUserID(ctx context.Context, signedCookie string) (string, string, error) {
	var sessionID string
	if err := m.sc.Decode("session_id", signedCookie, &sessionID); err != nil {
		m.logger.Err(err).Msg("Could not decode cookie value")
		return "", "", err
	}

	userID, err := m.rdb.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
	if err != nil {
		m.logger.Err(err).Msg("Could not retrieve userID")
		return "", "", err
	}

	return sessionID, userID, nil
}

func (m *Manager) HasSession(ctx context.Context, sessionID string) (bool, error) {
	exists, error := m.rdb.Exists(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
	return exists > 0, error
}
