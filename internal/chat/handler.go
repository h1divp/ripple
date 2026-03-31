package chat

import (
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/profile"
	"github.com/h1divp/echo-chat-v2/internal/session"
	"github.com/h1divp/echo-chat-v2/internal/websocket"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	userIdKey    = "user_id"
	sessionIdKey = "session_id"
)

type Handler struct {
	logger         zerolog.Logger
	service        *Service
	hub            *websocket.Hub
	profileService *profile.Service
	sessionManager *session.Manager
	upgrader       gorilla.Upgrader
	config         *config.Config
}

func NewHandler(logger zerolog.Logger, service *Service, hub *websocket.Hub, profileSvc *profile.Service, sessionMgr *session.Manager, cfg *config.Config) *Handler {
	allowedOrigins := config.Load().AllowedOrigins
	return &Handler{
		logger:         logger.With().Str("handler", "chat").Logger(),
		service:        service,
		hub:            hub,
		profileService: profileSvc,
		sessionManager: sessionMgr,
		upgrader:       websocket.NewUpgrader(allowedOrigins),
		config:         cfg,
	}
}

func (h *Handler) JoinChat(w http.ResponseWriter, r *http.Request) {
	// Check if session exists in redis
	sessionID, ok := r.Context().Value(sessionIdKey).(uuid.UUID)
	if !ok {
		http.Error(w, "You must have a session before joining the chat.", http.StatusUnauthorized)
		return
	}

	// Check if client exists in Hub
	userID, ok := r.Context().Value(userIdKey).(uuid.UUID)
	if !ok {
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

	clientLogger := h.logger.With().Str("userID", userID.String()).Logger()

	userProfile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		h.logger.Warn().Msg("Could not retrive profile during websocket handshake. Has the user generated a profile or logged in?")
		conn.WriteMessage(gorilla.CloseMessage, gorilla.FormatCloseMessage(
			gorilla.CloseProtocolError, "Profile not found"))
		conn.Close()
		return
	}

	client := websocket.NewClient(clientLogger, h.hub, conn, sessionID, userID, userProfile, h.config)
	h.hub.Register <- client
	go client.ReadPump(h.service)
	go client.WritePump()
	clientLogger.Info().Msg("Client connected to chat.")
}
