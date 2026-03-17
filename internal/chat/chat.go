package chat

import "time"

type Message struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	DisplayName string `json:"displayName"`
	AvatarSeed  string `json:"avatarSeed"`
	SenderID    string `json:"senderId"`
	Timestamp   int64  `json:"timestamp"`
	Status      string `json:"status"`
	Type        string `json:"type"`

	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func NewSystemMessage(text string) Message {
	return Message{
		ID:        ,
		Text:      text,
		Timestamp: time.Now().UnixMilli(),
		Type:      "system",
		Status:    "sent",
	}
}
