package chat

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/h1divp/echo-chat-v2/internal/chat/types"
)

var ErrCouldNotGetUserID = errors.New("Could not get userID from context.")

/*
Handlers for message types in types/message.go
*/

func (s *Service) HandleChatMessage(m *types.ChatMessage, ctx context.Context) error {
	userID, ok := ctx.Value(userIdKey).(uuid.UUID)
	if !ok {
		s.logger.Warn().Msg("Attempted to handle message by a user with no stored userID")
		return ErrCouldNotGetUserID
	}

	lat, lon, err := s.repo.GetLocationFromUserID(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not get location from userID")
	}

	// TODO: convert from messageSearchRadius to range stored in redis for user.
	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, lat, lon, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not find nearby users")
	}

	s.logger.Debug().Int("count", len(nearbyIDs)).Msg("Delivering message")
	s.hub.DeliverToLocalClients(nearbyIDs, m)
	// TODO: Put onto redis pub/sub
	return nil
}

func (s *Service) HandleLocationUpdate(m *types.LocationUpdate, ctx context.Context) error {
	// Update location in redis
	userID, ok := ctx.Value(userIdKey).(uuid.UUID)
	if !ok {
		s.logger.Warn().Msg("Attempted to handle message by a user with no stored userID")
		return ErrCouldNotGetUserID
	}

	err := s.repo.UpdateUserLocation(ctx, userID, m.Latitude, m.Longitude)
	if err != nil {
		s.logger.Err(err).Msg("Failed to update user location")
	}
	return nil
}

func (s *Service) HandleUsernameUpdate(m *types.UsernameUpdate, ctx context.Context) error {
	// TODO: Send to nearby users
	return nil
}

func (s *Service) HandleIconUpdate(m *types.IconUpdate, ctx context.Context) error {
	// TODO: Send to nearby users
	return nil
}
