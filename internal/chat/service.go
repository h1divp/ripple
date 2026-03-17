package chat

import (
	"context"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/rs/zerolog"
)

type Service struct {
	logger zerolog.Logger
	repo   *Repository
	hub    *Hub

	messageSearchRadius float64
}

func NewService(logger zerolog.Logger, repo *Repository, hub *Hub) *Service {
	return &Service{
		logger:              logger.With().Str("service", "chat").Logger(),
		repo:                repo,
		hub:                 hub,
		messageSearchRadius: config.Load().MessageSearchRadius,
	}
}

func (s *Service) HandleIncomingMessage(ctx context.Context, msg Message) error {
	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, msg.Latitude, msg.Longitude, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not handle incoming message")
		return err
	}

	s.hub.DeliverToLocalClients(nearbyIDs, msg)
	return nil
}
