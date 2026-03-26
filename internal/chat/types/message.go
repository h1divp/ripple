package types

import "github.com/google/uuid"

// Interface used for Service.ProcessIncomingMessage
type Message interface{}

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
