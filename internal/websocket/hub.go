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
	OnDisconnect func(userID string)

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
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)

				go h.OnDisconnect(client.ID)
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
