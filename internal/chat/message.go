package chat

import (
	"context"

	"github.com/google/uuid"
)

type Message interface {
	Handle(s *Service, ctx context.Context) error
}

type ChatMessage struct {
	ID        uuid.UUID `json:"id",omitempty`
	Text      string    `json:"text,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

type SystemMessage struct {
	ID uuid.UUID `json:"id",omitempty`
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
	Delta int `json:"delta",omitempty`
}

func (m *ChatMessage) Handle(s *Service, ctx context.Context) error {
	// get lat & lon from user_locations. the db is our source of truth :)
	// TODO: convert from messageSearchRadius to range stored in redis for user.
	nearbyIDs, err := s.repo.FindNearbyUserIDs(ctx, /*Latitude*/, /*Longitude*/, s.messageSearchRadius)
	if err != nil {
		s.logger.Err(err).Msg("Could not find nearby users")
	}

	recipients := append(nearbyIDs)

	logEvent.Int("count", len(recipients)).Msg("Delivering message")
	s.hub.DeliverToLocalClients(recipients, msg)
	// TODO: Put onto redis pub/sub
}

func (m *LocationUpdate) Handle() {
	// Update location in redis
	// get userID 
	err := s.repo.UpdateUserLocation(ctx, /*userID*/, m.Latitude, m.Longitude)
	if err != nil {
		s.logger.Err(err).Msg("Failed to update user location")
	}
	return nil
}

func (m *UsernameChange) Handle() {
	// Send to nearby users
}

func (m *IconChange) Handle() {
	// Send to nearby users
}
