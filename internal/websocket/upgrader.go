package websocket

import (
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

func NewUpgrader(allowedOrigins []string) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// A CheckOrigin bool is required for connections with a Svelte frontend
		CheckOrigin: func(r *http.Request) bool {
			if len(allowedOrigins) == 0 {
				return false
			}
			origin := r.Header.Get("Origin")
			return slices.Contains(allowedOrigins, origin)
		},
	}
}
