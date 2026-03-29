package websocket

import (
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Hub struct {
	clients      map[uuid.UUID]*Client
	Register     chan *Client
	Unregister   chan *Client
	OnDisconnect func(sessionID uuid.UUID, userID uuid.UUID)

	mu     sync.RWMutex
	logger zerolog.Logger
}

func NewHub(logger zerolog.Logger) *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
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
			// h.logger.Debug().Str("name", client.profile.DisplayName).Msg("Registered client")
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

func (h *Hub) DeliverToLocalClients(userIDs []uuid.UUID, msg any) {
	// h.logger.Debug().Msg("DeliverToLocalClients(): recieved message")

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, id := range userIDs {
		// h.logger.Debug().Str("userID", id.String()).Msg("Sending message to client")
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

func (h *Hub) HasClient(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clients[userID]
	return ok
}
