package user

import "context"

type User struct {
	ID          string
	Email       string
	DisplayName string
	TrustScore  int
}

type Session struct {
	ID          string
	DisplayName string
	IPAddress   string
	TrustScore  int
	IsAnon      bool
}

type Repository interface {
	GetProfile(ctx context.Context, userID string) (*User, error)

	SaveSession(ctx context.Context, s *Session) error
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	UpdateLocation(ctx context.Context, sessionID string, lat, lon float64) error
}
