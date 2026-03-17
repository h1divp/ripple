package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan Message
	ID   string

	logger zerolog.Logger
}

func NewClient(hub *Hub, conn *Conn, userID string, logger zerolog.Logger) *Client {
	return &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan Message),
		ID:     userID,
		logger: logger,
	}
}

func (c *Client) ReadPump(service *Service) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// The default read deadline isn't automatically extended by Gorilla
	// so we set a custom pong handler and deadline after every pong.
	// This is important so we can remove dead connections
	// accurately (the ones that miss the deadline).
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.logger.Info().Msg("Client disconnected")
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error().Err(err).Msg("unexpected close")
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			c.logger.Error().Err(err).Msg("Failed to unmarshal incoming message")
			continue
		}

		c.logger.Debug().
			Str("msgId", msg.ID).
			Str("text", msg.Text).
			Msg("Recieved incoming message")

		service.HandleIncomingMessage(context.Background(), msg)
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
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
