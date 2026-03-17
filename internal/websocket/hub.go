package websocket

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Hub struct {
	clients    map[string]*Client
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex

	logger zerolog.Logger
}

func NewHub(logger zerolog.Logger, rdb *redis.Client) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		logger: logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) DeliverToLocalClients(userIDs []string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, id := range userIDs {
		if client, ok := h.clients[id]; ok {
			select {
			case client.Send <- msg:
			default:
				// If a cliet's message buffer is full, then they are likely disconnected or are lagging,
				// so we skip them to stay performant.
			}
		}
	}
}
