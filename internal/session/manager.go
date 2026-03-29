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

func (m *Manager) CreateSession(ctx context.Context) (string, uuid.UUID, error) {
	_, ok := ctx.Value(userCtxKey).(uuid.UUID)
	if ok {
		// m.logger.Debug().Str("userID", val.String()).Msg("Session already exists")
		return "", uuid.Nil, ErrSessionAlreadyExists
	}

	sessionID := uuid.New()
	userID := uuid.New()

	err := m.rdb.Set(ctx, fmt.Sprintf("session:%s", sessionID.String()), userID.String(), m.expiry).Err()
	if err != nil {
		return "", uuid.Nil, err
	}

	// TODO: test to see if the signing works correctly (try copying into a different web browser)
	signed, err := m.sc.Encode("session_id", sessionID.String())
	m.logger.Debug().Str("sessionID", sessionID.String()).Msg("signing sessionID")
	m.logger.Debug().Msg("created session")
	return signed, userID, err
}

func (m *Manager) GetSessionIDAndUserID(ctx context.Context, signedCookie string) (uuid.UUID, uuid.UUID, error) {
	var sessionIDString string
	if err := m.sc.Decode("session_id", signedCookie, &sessionIDString); err != nil {
		m.logger.Err(err).Msg("Could not decode cookie value")
		return uuid.Nil, uuid.Nil, err
	}
	sessionID, err := uuid.Parse(sessionIDString)
	if err != nil {
		m.logger.Err(err).Msg("decrypted sessionID is malformed")
		return uuid.Nil, uuid.Nil, err
	}

	rawUserID, err := m.rdb.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
	if err != nil {
		m.logger.Err(err).Msg("Could not retrieve userID")
		return uuid.Nil, uuid.Nil, err
	}
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		m.logger.Warn().Msg("userID in redis was not stored as a parsable uuid")
		return uuid.Nil, uuid.Nil, err
	}

	return sessionID, userID, nil
}

func (m *Manager) HasSession(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	exists, error := m.rdb.Exists(ctx, fmt.Sprintf("session:%s", sessionID.String())).Result()
	return exists > 0, error
}

func (m *Manager) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	err := m.rdb.Del(ctx, fmt.Sprintf("session:%s", sessionID.String())).Err()
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
