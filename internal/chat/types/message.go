package types

type Message struct {
	ID        string  `json:"id"`
	SenderID  string  `json:"senderId"`
	Type      string  `json:"type"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`

	Text        *string `json:"text,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	AvatarSeed  *string `json:"avatarSeed,omitempty"`
	Timestamp   *int64  `json:"timestamp,omitempty"`
	Status      *string `json:"status"`
	MessageType *string `json:"msgType"`
}

const (
	ChatMsgType   = "chat"
	SystemMsgType = "system"
)
