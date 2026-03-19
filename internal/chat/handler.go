package chat

import (
	"net/http"

	gorilla "github.com/gorilla/websocket"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/websocket"
	"github.com/rs/zerolog"
)

const sessionIdKey = "session_id"

type Handler struct {
	logger   zerolog.Logger
	service  *Service
	hub      *websocket.Hub
	upgrader gorilla.Upgrader
}

func NewHandler(logger zerolog.Logger, service *Service, hub *websocket.Hub) *Handler {
	allowedOrigins := config.Load().AllowedOrigins
	return &Handler{
		logger:   logger.With().Str("handler", "chat").Logger(),
		service:  service,
		hub:      hub,
		upgrader: websocket.NewUpgrader(allowedOrigins),
	}
}

func (h *Handler) JoinChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(sessionIdKey).(string)
	if !ok || userID == "" {
		http.Error(w, "You must have a session before joining the chat.", http.StatusUnauthorized)
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Err(err).Msg("Failed to upgrade connection to websocket")
		// Gorilla handles closing the connection here. No need for conn.Close() (this will result in a panic)
		return
	}

	clientLogger := h.logger.With().Str("userID", userID).Logger()

	client := websocket.NewClient(h.hub, conn, userID, clientLogger)
	h.hub.Register <- client

	go client.ReadPump(h.service)
	go client.WritePump()

	clientLogger.Info().Msg("Client connected successfully.")
}
