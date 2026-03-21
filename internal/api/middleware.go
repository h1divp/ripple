package api

import (
	"context"
	"net/http"

	"github.com/h1divp/echo-chat-v2/internal/session"
	"github.com/rs/zerolog"
)

const (
	sessionIdKey  = "session_id"
	sessionCtxKey = "session_id"
	userCtxKey    = "user_id"
)

func SessionMiddleware(sm *session.Manager, logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionIdKey)
			if err != nil {
				logger.Debug().Msg("no session cookie found")
				next.ServeHTTP(w, r)
				return
			}

			var userID string
			var sessionID string
			sessionID, userID, err = sm.GetSessionIDAndUserID(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			logger.Debug().Str("userID", userID).Str("sessionID", sessionID).Msg("passing through middleware")

			ctx := context.WithValue(r.Context(), sessionCtxKey, sessionID)
			ctx = context.WithValue(ctx, userCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
