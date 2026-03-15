package user

import (
	"context"
	"net/http"
)

type User struct {
	ID          string
	Email       string
	DisplayName string
	TrustScore  int
}

type Session struct {
	ID          string `redis:"session_id"`
	DisplayName string `redis:"display_name"`
	IPAddress   string `redis:"ip_address"`
	TrustScore  int    `redis:"trust_score"`
}

type UserRepository interface {
	// CreateProfile(ctx context.Context, userID string) (*User, error)
	// GetProfile(ctx context.Context, userID string) (*User, error)

	SaveSession(ctx context.Context, s *Session) error
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	// UpdateLocation(ctx context.Context, sessionID string, lat, lon float64) error
}

type UserService interface {
	JoinChat(ctx context.Context, displayName string, ip string) (*Session, error)
}

type UserHandler interface {
	Connect(w http.ResponseWriter, r *http.Request)
}
