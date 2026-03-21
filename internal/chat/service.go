package chat

import (
	"context"

	"github.com/h1divp/echo-chat-v2/internal/chat/types"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/rs/zerolog"
)

type Service struct {
	logger              zerolog.Logger
	repo                *Repository
	hub                 HubInterface
	sessionManager      SessionInterface
	messageSearchRadius float64
}

type HubInterface interface {
	DeliverToLocalClients(userIDs []string, msg types.Message)
}

type SessionInterface interface {
	DeleteSession(ctx context.Context, sessionID string) error
}

func NewService(logger zerolog.Logger, repo *Repository, hub HubInterface, sessionMgr SessionInterface) *Service {
	log := logger.With().Str("service", "chat").Logger()
	radius := config.Load().MessageSearchRadius
	if radius <= 0 {
		log.Warn().Msg("Fallback is being used for message search radius (50 meters).")
		radius = 50.0
	}
	return &Service{
		logger:              log,
		repo:                repo,
		hub:                 hub,
		sessionManager:      sessionMgr,
		messageSearchRadius: radius,
	}
}

func (s *Service) HandleIncomingMessage(ctx context.Context, msg types.Message) error {
	logEvent := s.logger.Info().
		Float64("lat", msg.Latitude).
		Float64("lon", msg.Longitude)
	if msg.Text != nil {
		logEvent.Str("msg", *msg.Text)
	}

	// UpdateUserLocation: runs for both chat and location_update message types.
	err := s.repo.UpdateUserLocation(ctx, msg.SenderID, msg.Latitude, msg.Longitude)
	if err != nil {
		s.logger.Err(err).Msg("Failed to update user location")
	}

	// TODO: refactor
	if msg.Type == types.ChatMsgType {
		nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, msg.Latitude, msg.Longitude, s.messageSearchRadius)
		if err != nil {
			s.logger.Err(err).Msg("Could not find nearby users")
		}

		recipients := append(nearbyIDs, msg.SenderID)

		logEvent.Int("count", len(recipients)).Msg("Delivering message")
		s.hub.DeliverToLocalClients(recipients, msg)
		// TODO: Put onto redis pub/sub
	}
	return nil
}

func (s *Service) HandleDisconnect(ctx context.Context, sessionID string, userID string) {
	err := s.repo.RemoveUserLocation(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not delete location from redis while disconnecting user")
	}
	err = s.sessionManager.DeleteSession(ctx, sessionID)
	if err != nil {
		s.logger.Err(err).Msg("Could not delete session from redis while disconnecting user")
	}
	s.logger.Info().Msg("Client disconnected from chat.")
}
