package profile

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const userIdKey = "user_id"

type Handler struct {
	logger  zerolog.Logger
	service *Service
}

func NewHandler(logger zerolog.Logger, service *Service) *Handler {
	return &Handler{
		logger:  logger.With().Str("handler", "profile").Logger(),
		service: service,
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdKey).(uuid.UUID)
	if !ok {
		http.Error(w, "You must have a session before getting a profile.", http.StatusUnauthorized)
		h.logger.Warn().Msg("User attempted to retreive profile data before obtaining a session.")
		return
	}

	profile, err := h.service.GetOrCreateProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "Could not retrieve profile.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		h.logger.Err(err).Msg("Failed to encode profile response")
		return
	}
}
