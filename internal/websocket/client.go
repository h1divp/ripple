package websocket

import (
	"context"
	"encoding/json"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/h1divp/echo-chat-v2/internal/chat/types"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/profile"
	"github.com/rs/zerolog"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub                    *Hub
	Conn                   *websocket.Conn
	Send                   chan any // Can be used for multiple message types
	OldestMessageTimestamp time.Time

	sessionID         uuid.UUID
	userID            uuid.UUID
	profile           *profile.Profile
	config            *config.Config
	messageTimestamps []time.Time
	rateLimitMutex    sync.RWMutex

	logger zerolog.Logger
}

func NewClient(logger zerolog.Logger, hub *Hub, conn *websocket.Conn, sessionID uuid.UUID, userID uuid.UUID, profile *profile.Profile, cfg *config.Config) *Client {
	return &Client{
		logger:    logger,
		Hub:       hub,
		Conn:      conn,
		Send:      make(chan any),
		sessionID: sessionID,
		userID:    userID,
		profile:   profile,
		config:    cfg,
	}
}

type MessageHandler interface {
	// Defined in chat/service
	ProcessIncomingMessage(ctx context.Context, msg types.Message, userID uuid.UUID, profile *profile.Profile) error
}

func (c *Client) ReadPump(handler MessageHandler) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// For safe and accurate cleanup we use timeouts for connections. We "ping" our
	// connection and expect a "pong" in return. If the pong doesn't make a deadline,
	// we assume that the connection is dead and clean it up.
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.logger.Debug().Msg("Client recieved pong")
		return nil
	})

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		// c.logger.Debug().Msg("Readpump recieved websocket message")
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error().Err(err).Msg("unexpected close")
			}
			break
		}

		var envelope struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawMsg, &envelope); err != nil {
			c.logger.Err(err).Msg("Could not parse incoming message.")
			continue
		}

		var msg types.Message
		switch envelope.Type {
		case "chat":
			msg = &types.ChatMessageInbound{}
		case "location_update":
			msg = &types.LocationUpdate{}
		case "username_update":
			msg = &types.UsernameUpdate{}
		case "icon_update":
			msg = &types.IconUpdate{}
		default:
			c.logger.Warn().Str("type", envelope.Type).Msg("Recieved unknown message type")
			continue
		}

		if err := json.Unmarshal(rawMsg, msg); err != nil {
			c.logger.Err(err).Msg("Failed to unmarshall concrete message.")
			continue
		}

		// TODO: refactor
		handler.ProcessIncomingMessage(context.Background(), msg, c.userID, c.profile)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// c.logger.Debug().Msg("Writepump recieved message")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			// c.logger.Debug().Msg("Writepump recieved ticker")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) isRateLimited() (bool, time.Time, error) {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-c.config.Chat.RateLimitWindowSeconds)

	// Remove old timestamps
	c.messageTimestamps = slices.DeleteFunc(c.messageTimestamps,
		func(t time.Time) bool { return t.Before(cutoff) })

	if len(c.messageTimestamps) >= c.config.Chat.RateLimitMaxMessages {
		// The oldest message that's still counting against the limit
		// When this message expires, the user will no longer be rate limited
		oldestRelevantMessage := c.messageTimestamps[0]
		endTime := oldestRelevantMessage.Add(c.config.Chat.RateLimitWindowSeconds)
		return true, endTime, nil
	}

	c.messageTimestamps = append(c.messageTimestamps, now)
	return false, time.Time{}, nil
}
