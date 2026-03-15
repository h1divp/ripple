package user

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Needed to interact with Svelte(Kit) during development. Should be restricted to domain for prod.
		return true
	},
}

type Handler struct {
	logger zerolog.Logger
	svc    *Service
}

func NewHandler(logger zerolog.Logger, svc *Service) *Handler {
	return &Handler{
		logger: logger.With().Str("handler", "user").Logger(),
		svc:    svc,
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	displayName := r.URL.Query().Get("displayName")
	// TODO: randomize name
	if displayName == "" {
		displayName = "Anonymous"
	}
	ip := r.RemoteAddr

	sess, err := h.svc.JoinChat(ctx, displayName, ip)
	if err != nil {
		h.logger.Err(err).Msg("Failed to join chat")
		http.Error(w, "Internal server error", http.StatusInternalServerError) // TODO: Maybe should refactor into api/response for reusability
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Err(err).Msg("Failed to upgrade connection from http to websocket")
		// By this point status code status code 101 (Switching Protocols) has been sent,
		// so we do not need to fail with an http.Error()
		return
	}

	h.logger.Info().Str("session_id", sess.ID).Msg("User connected to chat")

	go h.handleConnection(conn, sess)
}

func (h *Handler) handleConnection(conn *websocket.Conn, sess *Session) {
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			h.logger.Info().Str("session_id", sess.ID).Msg("User disconnected from chat")
			return
		}

		h.logger.Debug().Str("payload", string(p)).Msg("Message recieved")
		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}
