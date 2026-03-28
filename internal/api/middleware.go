package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/h1divp/echo-chat-v2/internal/session"
	"github.com/rs/zerolog"
)

const (
	sessionIdKey  = "session_id"
	sessionCtxKey = "session_id"
	userCtxKey    = "user_id"
)

// TODO: fix by decoding session cookie into string, then parsing into uuid
func SessionMiddleware(sm *session.Manager, logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionIdKey)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			var userID uuid.UUID
			var sessionID uuid.UUID
			sessionID, userID, err = sm.GetSessionIDAndUserID(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), sessionCtxKey, sessionID)
			ctx = context.WithValue(ctx, userCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
