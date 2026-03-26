package chat

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrCouldNotGetUserID = errors.New("Could not get userID from context.")

type Message interface {
	Handle(s *Service, ctx context.Context) error
}

type ChatMessage struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Text      string    `json:"text,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

type SystemMessage struct {
	ID uuid.UUID `json:"id,omitempty"`
}

// Range included for future ability to send chat messages to users
// 50m/500m/5000m away.
type LocationUpdate struct {
	Latitude  float64  `json:"lat,omitempty"`
	Longitude float64  `json:"lon,omitempty"`
	Range     *float64 `json:"distance,omitempty"`
}

// A request for these should be sent from the client. The server should generate
// these and then send the updated versions to all nearby users (including the
// requester).
// We don't want the client to generate their stored username or icon
// url since that leaves a security and moderation risk.
type UsernameUpdate struct {
	Username string `json:"username,omitempty"`
}

type IconUpdate struct {
	IconURL string `json:"icon_url,omitempty"`
}

// Represents a change in the number of users within an area.
// This should be sent out after a new connection sends its first location update,
// but in situations where there is a high volume of join/leaves we can use a delta
// to prevent needing to send many messages.
type NearbyUpdate struct {
	Delta int `json:"delta,omitempty"`
}

func (m *ChatMessage) Handle(s *Service, ctx context.Context) error {
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
	// TODO: update DeliverToLocalClients to accept new message type
	s.hub.DeliverToLocalClients(nearbyIDs, m)
	// TODO: Put onto redis pub/sub
	return nil
}

func (m *LocationUpdate) Handle(s *Service, ctx context.Context) error {
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

func (m *UsernameUpdate) Handle(s *Service, ctx context.Context) error {
	// Send to nearby users
	return nil
}

func (m *IconUpdate) Handle(s *Service, ctx context.Context) error {
	// Send to nearby users
	return nil
}
