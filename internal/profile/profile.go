package profile

import (
	"time"
)

type Profile struct {
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
}
