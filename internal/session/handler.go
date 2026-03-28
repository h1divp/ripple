package session

import (
	"net/http"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/rs/zerolog"
)

const sessionIdKey = "session_id"

type Handler struct {
	logger  zerolog.Logger
	manager *Manager
	config  config.Config
}

func NewHandler(logger zerolog.Logger, mgr *Manager, cfg config.Config) *Handler {
	return &Handler{
		logger:  logger,
		manager: mgr,
		config:  cfg,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	signed, _, err := h.manager.CreateSession(r.Context())
	if err == ErrSessionAlreadyExists {
		http.Error(w, "User already has a session", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Err(err).Msg("error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionIdKey,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.IsProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600, // TODO: session to be extended once a message is sent
	})
	w.WriteHeader(http.StatusOK)
}
