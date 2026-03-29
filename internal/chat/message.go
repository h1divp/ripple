package chat

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/h1divp/echo-chat-v2/internal/chat/types"
	"github.com/h1divp/echo-chat-v2/internal/profile"
)

var ErrCouldNotGetUserID = errors.New("Could not get userID from context.")

/*
Handlers for message types in types/message.go
*/

func (s *Service) HandleChatMessage(ctx context.Context, m *types.ChatMessageInbound, userID uuid.UUID, profile *profile.Profile) error {
	lat, lon, err := s.repo.GetLocationFromUserID(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not get location from userID")
	}

	s.logger.Debug().Str("content", m.Text).Msg("Recieved chat message")

	// TODO: convert from messageSearchRadius to range stored in redis for user.
	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, lat, lon, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not find nearby users")
	}

	s.logger.Debug().Int("count", len(nearbyIDs)).Msg("Delivering message")

	outboundMsg := types.ChatMessageOutbound{
		ID:          m.ID,
		Text:        m.Text,
		Timestamp:   m.Timestamp,
		Type:        "chat",
		Status:      "sent",
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
	}

	s.hub.DeliverToLocalClients(nearbyIDs, &outboundMsg)
	// TODO: Put onto redis pub/sub
	return nil
}

func (s *Service) HandleLocationUpdate(ctx context.Context, m *types.LocationUpdate, userID uuid.UUID) error {
	// TODO: how should we get the userID?
	err := s.repo.UpdateUserLocation(ctx, userID, m.Latitude, m.Longitude)
	if err != nil {
		s.logger.Err(err).Msg("Failed to update user location")
	}
	return nil
}

func (s *Service) HandleUsernameUpdate(ctx context.Context, m *types.UsernameUpdate) error {
	// TODO: Send to nearby users
	return nil
}

func (s *Service) HandleIconUpdate(ctx context.Context, m *types.IconUpdate) error {
	// TODO: Send to nearby users
	return nil
}
