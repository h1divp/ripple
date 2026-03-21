package websocket

import (
	"sync"

	"github.com/h1divp/echo-chat-v2/internal/chat/types"
	"github.com/rs/zerolog"
)

type Hub struct {
	clients      map[string]*Client
	Register     chan *Client
	Unregister   chan *Client
	OnDisconnect func(sessionID string, userID string)

	mu     sync.RWMutex
	logger zerolog.Logger
}

func NewHub(logger zerolog.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, exists := h.clients[client.userID]; exists {
				// Reject duplicate connection
				h.mu.Unlock()
				client.Conn.Close()
				continue
			}
			h.clients[client.userID] = client
			h.logger.Debug().Str("userID", client.userID).Msg("Registered client")
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()

			if existing, ok := h.clients[client.userID]; ok && existing == client {
				delete(h.clients, client.userID)
				close(client.Send)
				go h.OnDisconnect(client.sessionID, client.userID)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) DeliverToLocalClients(userIDs []string, msg types.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, id := range userIDs {
		if client, ok := h.clients[id]; ok {
			select {
			case client.Send <- msg:
			default:
				// If a client's message buffer is full, then they are likely disconnected or are lagging,
				// so we skip them to stay performant.
			}
		}
	}
}

func (h *Hub) HasClient(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clients[userID]
	return ok
}
