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
	messageSearchRadius float64
}

type HubInterface interface {
	DeliverToLocalClients(userIDs []string, msg types.Message)
}

func NewService(logger zerolog.Logger, repo *Repository, hub HubInterface) *Service {
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
		messageSearchRadius: radius,
	}
}

func (s *Service) HandleIncomingMessage(ctx context.Context, msg types.Message) error {
	logEvent := s.logger.Debug().
		Float64("lat", msg.Latitude).
		Float64("lon", msg.Longitude)
	if msg.Text != nil {
		logEvent.Str("msg", *msg.Text)
	}
	logEvent.Msg("Processing incoming message")

	// UpdateUserLocation: runs for both chat and location_update message types.
	err := s.repo.UpdateUserLocation(ctx, msg.SenderID, msg.Latitude, msg.Longitude)
	if err != nil {
		s.logger.Err(err).Msg("Failed to update user location")
	}

	if msg.Type == types.ChatMsgType {
		nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, msg.Latitude, msg.Longitude, s.messageSearchRadius)
		s.logger.Debug().Msgf("nearby %v", len(nearbyIDs))
		if err != nil {
			s.logger.Err(err).Msg("Could not find nearby users")
		}

		recipients := append(nearbyIDs, msg.SenderID)

		s.logger.Debug().Int("count", len(recipients)).Msg("Delivering message")
		s.hub.DeliverToLocalClients(recipients, msg)
		// TODO: Put onto redis pub/sub
	}
	return nil
}

func (s *Service) HandleDisconnect(ctx context.Context, userID string) {
	err := s.repo.RemoveUserLocation(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Error while disconnecting user")
	}
}
