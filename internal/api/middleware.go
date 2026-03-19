package api

import (
	"context"
	"net/http"

	"github.com/h1divp/echo-chat-v2/internal/session"
)

const (
	sessionIdKey = "session_id"
	userCtxKey   = "user_id"
)

func SessionMiddleware(sm *session.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionIdKey)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			var userID string
			userID, err = sm.GetUserID(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
