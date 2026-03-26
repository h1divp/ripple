package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/h1divp/echo-chat-v2/internal/chat/types"
	"github.com/rs/zerolog"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	// TODO: convert id types from string to uuid
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan types.Message
	sessionID string
	userID    string

	logger zerolog.Logger
}

func NewClient(hub *Hub, conn *websocket.Conn, sessionID string, userID string, logger zerolog.Logger) *Client {
	return &Client{
		Hub:       hub,
		Conn:      conn,
		Send:      make(chan types.Message),
		sessionID: sessionID,
		userID:    userID,
		logger:    logger,
	}
}

type MessageHandler interface {
	ProcessIncomingMessage(ctx context.Context, msg types.Message) error
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
			msg = &types.ChatMessage{}
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

		handler.ProcessIncomingMessage(context.Background(), msg)
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
			c.logger.Debug().Msg("Writepump recieved message")
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
			c.logger.Debug().Msg("Writepump recieved ticker")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
