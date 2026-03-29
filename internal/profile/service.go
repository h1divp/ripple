package profile

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var ErrCouldNotGetProfile = errors.New("Could not get profile from userID.")

type Service struct {
	logger   zerolog.Logger
	profiles map[uuid.UUID]Profile
}

func NewService(logger zerolog.Logger) *Service {
	return &Service{
		logger:   logger.With().Str("service", "profile").Logger(),
		profiles: make(map[uuid.UUID]Profile),
	}
}

// Used explicitly for pre websocket profile information generation and retrieval.
func (s *Service) GetOrCreateProfile(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	// 1. Is the userID in the map?
	//    If so, return from map
	profile, ok := s.profiles[userID]
	if ok {
		return &profile, nil
	}
	// 2. Is the user anonymous or logged in?
	//    If anon, generate temporary profile and insert into map
	//    If logged in, TODO: retrieve data from postgres
	//    -> return profile struct
	profile, err := s.CreateAnonymousProfile(userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not create anonymous profile.")
		return nil, err
	}
	return &profile, nil
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	if s.profiles == nil {
		s.logger.Warn().Msg("Profiles map is nil!")
		return nil, ErrCouldNotGetProfile
	}
	profile, ok := s.profiles[userID]
	if !ok {
		return nil, ErrCouldNotGetProfile
	}
	return &profile, nil
}

func (s *Service) CreateAnonymousProfile(userID uuid.UUID) (Profile, error) {
	// TODO: handle errors
	profile := Profile{
		DisplayName: GenerateRandomName(),
		AvatarURL:   GenerateRandomAvatarURL(),
		CreatedAt:   time.Now(),
	}
	s.profiles[userID] = profile
	return profile, nil
}
