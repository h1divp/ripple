package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/h1divp/echo-chat-v2/internal/chat"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/session"
)

type Api struct {
	Router         *chi.Mux
	Logger         *zerolog.Logger
	SessionHandler *session.Handler
	SessionManager *session.Manager
	ChatHandler    *chat.Handler
}

func New(logger *zerolog.Logger, sessionHdl *session.Handler, sessionMgr *session.Manager, chatHdl *chat.Handler) *Api {
	r := chi.NewRouter()
	r.Use(SessionMiddleware(sessionMgr))
	r.Use(middleware.Recoverer)

	api := &Api{
		Router:      r,
		Logger:      logger,
		ChatHandler: chatHdl,
	}

	api.CreateRoutes()

	return api
}

func (api *Api) CreateRoutes() {
	AllowedOrigins := config.Load().AllowedOrigins

	api.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api.Router.Get("/ws/chat", api.ChatHandler.JoinChat)
	api.Router.Get("/register", api.SessionHandler.Register)
}
