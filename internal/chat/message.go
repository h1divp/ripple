package chat

import (
	"context"
	"errors"
	"slices"
	"time"

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
		return err
	}
	s.logger.Debug().Str("content", m.Text).Msg("Recieved chat message")

	isRateLimited, rateLimitEnd, err := s.hub.IsClientRateLimited(userID)
	if err != nil {
		s.logger.Err(err).Msg("Error while checking if user is rate limited")
		isRateLimited = false
	}
	if isRateLimited {
		s.logger.Warn().Str("userID", userID.String()).Msg("User rate limited")

		if (rateLimitEnd.Equal(time.Time{})) {
			s.logger.Warn().Msg("end time returned for user rate limit is a 0 value")
		}

		// Send rate limit error back to the user
		errorMsg := &types.SystemMessage{
			ID:               uuid.New(),
			Type:             types.MessageTypeSystem,
			Code:             types.SystemMessageTooManyMessages,
			Text:             "You're sending messages too quickly. Please slow down.",
			RateLimitEndTime: rateLimitEnd,
		}
		s.hub.DeliverToLocalClients([]uuid.UUID{userID}, errorMsg)
		return nil
	}

	if len(m.Text) > s.config.Chat.MaxMessageLength {
		s.logger.Warn().Int("length", len(m.Text)).Str("userID", userID.String()).Msg("Recieved message exceeding max message length.")

		errorMsg := &types.SystemMessage{
			ID:               uuid.New(),
			Type:             types.MessageTypeSystem,
			Code:             types.SystemMessageMessageTooLong,
			Text:             "Message is too long.",
			IsConsoleMessage: true,
		}
		s.hub.DeliverToLocalClients([]uuid.UUID{userID}, errorMsg)
		return nil
	}

	// TODO: convert from messageSearchRadius to range stored in redis for user.
	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, lat, lon, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not find nearby users")
		return err
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
	hasPreviousLocation, err := s.repo.HasUserLocation(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not determine if a user has a location stored in Redis")
	}

	err = s.repo.UpdateUserLocation(ctx, userID, m.Latitude, m.Longitude)
	if err != nil {
		s.logger.Err(err).Msg("Failed to update user location")
	}

	// s.logger.Debug().Float64("lat", m.Latitude).Float64("lon", m.Longitude).Msg("")

	if !hasPreviousLocation {
		// First location ping
		s.BroadcastJoinNearbyUpdates(ctx, userID)
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

func (s *Service) BroadcastNearbyUpdate(ctx context.Context, userID uuid.UUID, delta int) error {
	lat, lon, err := s.repo.GetLocationFromUserID(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not get location for nearby update")
		return err
	}

	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, lat, lon, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not find nearby users for update")
		return err
	}
	// Remove joining user from broadcast. We send a seperate message to them specifically.
	nearbyIDs = slices.DeleteFunc(nearbyIDs, func(id uuid.UUID) bool {
		return id == userID
	})

	nearbyBroadcastMsg := &types.NearbyUpdate{
		Type:  types.MessageTypeNearbyUpdate,
		Delta: delta,
	}
	s.hub.DeliverToLocalClients(nearbyIDs, nearbyBroadcastMsg)

	return nil
}

func (s *Service) BroadcastJoinNearbyUpdates(ctx context.Context, userID uuid.UUID) error {
	lat, lon, err := s.repo.GetLocationFromUserID(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("Could not get location for nearby update")
		return err
	}

	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, lat, lon, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not find nearby users for update")
		return err
	}
	// Remove joining user from broadcast. We send a seperate message to them specifically.
	nearbyIDs = slices.DeleteFunc(nearbyIDs, func(id uuid.UUID) bool {
		return id == userID
	})

	nearbyBroadcastMsg := &types.NearbyUpdate{
		Type:  types.MessageTypeNearbyUpdate,
		Delta: 1,
	}
	s.hub.DeliverToLocalClients(nearbyIDs, nearbyBroadcastMsg)

	nearbyToJoiningUserMsg := &types.NearbyUpdate{
		Type:  types.MessageTypeNearbyUpdate,
		Delta: len(nearbyIDs),
	}
	s.hub.DeliverToLocalClients([]uuid.UUID{userID}, nearbyToJoiningUserMsg)
	return nil
}
