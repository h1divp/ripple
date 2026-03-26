package chat

import (
	"net/http"

	gorilla "github.com/gorilla/websocket"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/session"
	"github.com/h1divp/echo-chat-v2/internal/websocket"
	"github.com/rs/zerolog"
)

const (
	userIdKey    = "user_id"
	sessionIdKey = "session_id"
)

type Handler struct {
	logger         zerolog.Logger
	service        *Service
	hub            *websocket.Hub
	sessionManager *session.Manager
	upgrader       gorilla.Upgrader
}

func NewHandler(logger zerolog.Logger, service *Service, hub *websocket.Hub, sessionMgr *session.Manager) *Handler {
	allowedOrigins := config.Load().AllowedOrigins
	return &Handler{
		logger:         logger.With().Str("handler", "chat").Logger(),
		service:        service,
		hub:            hub,
		sessionManager: sessionMgr,
		upgrader:       websocket.NewUpgrader(allowedOrigins),
	}
}

func (h *Handler) JoinChat(w http.ResponseWriter, r *http.Request) {
	// Check if session exists in redis
	sessionID, ok := r.Context().Value(sessionIdKey).(string)
	if !ok || sessionID == "" {
		http.Error(w, "You must have a session before joining the chat.", http.StatusUnauthorized)
		return
	}

	// Check if client exists in Hub
	userID, ok := r.Context().Value(userIdKey).(string)
	if !ok || userID == "" {
		h.logger.Info().Msg("User attempted to join chat before session was created")
		http.Error(w, "You must have a session before joining the chat.", http.StatusUnauthorized)
		return
	}
	if h.hub.HasClient(userID) {
		h.logger.Info().Msg("User attempted to join chat while having an active connection")
		http.Error(w, "You already have an active connection.", http.StatusTooManyRequests)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Err(err).Msg("Failed to upgrade connection to websocket")
		// Gorilla handles closing the connection here. No need for conn.Close() (this will result in a panic)
		return
	}

	clientLogger := h.logger.With().Str("userID", userID).Logger()

	client := websocket.NewClient(h.hub, conn, sessionID, userID, clientLogger)
	h.hub.Register <- client
	go client.ReadPump(h.service)
	go client.WritePump()
	clientLogger.Info().Msg("Client connected to chat.")
}
