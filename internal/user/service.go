package user

import (
	"context"
	"github.com/google/uuid"
)

// TODO: JoinChat, ChangeName

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) JoinChat(ctx context.Context, displayName string, ip string) (*Session, error) {
	sessionID := uuid.New().String()

	sess := &Session{
		ID:          sessionID,
		DisplayName: displayName,
		IPAddress:   ip,
		TrustScore:  0,
	}

	if err := s.repo.SaveSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}
