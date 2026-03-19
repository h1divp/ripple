package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/redis/go-redis/v9"
)

var (
	ErrSessionAlreadyExists = errors.New("this user already has a session cookie")
)

const userCtxKey = "user_id"

type Manager struct {
	rdb    *redis.Client
	sc     *securecookie.SecureCookie
	expiry time.Duration
}

func NewManager(rdb *redis.Client, hashKey, blockKey []byte) *Manager {
	return &Manager{
		rdb:    rdb,
		sc:     securecookie.New(hashKey, blockKey),
		expiry: 24 * time.Hour,
	}
}

func (m *Manager) CreateSession(ctx context.Context) (string, string, error) {
	_, ok := ctx.Value(userCtxKey).(string)
	if ok {
		return "", "", ErrSessionAlreadyExists
	}

	sessionID := uuid.New().String()
	userID := uuid.New().String()

	err := m.rdb.Set(ctx, fmt.Sprintf("session:%s", sessionID), userID, m.expiry).Err()
	if err != nil {
		return "", "", err
	}

	signed, err := m.sc.Encode("session_id", sessionID)
	return signed, userID, err
}

func (m *Manager) GetUserID(ctx context.Context, signedCookie string) (string, error) {
	var sessionID string
	if err := m.sc.Decode("session_ID", signedCookie, &sessionID); err != nil {
		return "", err
	}

	userID, err := m.rdb.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}
