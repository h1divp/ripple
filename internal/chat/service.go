package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/h1divp/echo-chat-v2/internal/chat/types"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/profile"
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
	// Defined in websocket/hub.go
	DeliverToLocalClients(userIDs []uuid.UUID, msg any)
}

type SessionInterface interface {
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
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

func (s *Service) ProcessIncomingMessage(ctx context.Context, msg types.Message, userID uuid.UUID, profile *profile.Profile) error {
	// Handle message depending on type (e.g. ChatMessageInbound, LocationUpdate, etc)
	// TODO: error?
	switch msg := msg.(type) {
	case *types.ChatMessageInbound:
		s.HandleChatMessage(ctx, msg, userID, profile)
	case *types.LocationUpdate:
		s.HandleLocationUpdate(ctx, msg, userID)
	}
	return nil
}

func (s *Service) HandleDisconnect(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) {
	s.BroadcastNearbyUpdate(context.Background(), userID, -1)

	err := s.repo.RemoveUserLocation(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not delete location from redis while disconnecting user")
	}

	// err = s.sessionManager.DeleteSession(ctx, sessionID)
	// if err != nil {
	// s.logger.Err(err).Msg("Could not delete session from redis while disconnecting user")
	// }

	s.logger.Info().Msg("Client disconnected from chat.")
}
