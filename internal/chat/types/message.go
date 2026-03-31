package types

import (
	"time"

	"github.com/google/uuid"
)

const (
	MessageTypeChatInbound    = "chat"
	MessageTypeChatOutbound   = "chat"
	MessageTypeLocationUpdate = "location_update"
	MessageTypeUsernameUpdate = "username_update"
	MessageTypeIconUpdate     = "icon_update"
	MessageTypeNearbyUpdate   = "nearby_update"
	MessageTypeSystem         = "system"
)

const (
	SystemMessageTooManyMessages = "too_many_messages"
	SystemMessageMessageTooLong  = "message_too_long"
)

// General message type used for Service.ProcessIncomingMessage
type Message any

type ChatMessageInbound struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Text      string    `json:"text,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

type ChatMessageOutbound struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Text        string    `json:"text,omitempty"`
	Timestamp   int64     `json:"timestamp,omitempty"`
	Type        string    `json:"type,omitempty"`
	Status      string    `json:"status,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
}

type SystemMessage struct {
	ID               uuid.UUID `json:"id,omitempty"`
	Type             string    `json:"type,omitempty"`
	Code             string    `json:"code,omitempty"`
	Text             string    `json:"text,omitempty"`
	IsConsoleMessage bool      `json:"is_console_message"`
	RateLimitEndTime time.Time `json:"rate_limit_end_time"`
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
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}
