package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/h1divp/echo-chat-v2/internal/chat"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/profile"
	"github.com/h1divp/echo-chat-v2/internal/session"
)

type Api struct {
	Router         *chi.Mux
	SessionHandler *session.Handler
	SessionManager *session.Manager
	ChatHandler    *chat.Handler
	ProfileHandler *profile.Handler

	logger *zerolog.Logger
}

func New(logger *zerolog.Logger, sessionHdl *session.Handler, sessionMgr *session.Manager, chatHdl *chat.Handler, profileHdl *profile.Handler) *Api {
	r := chi.NewRouter()

	AllowedOrigins := config.Load().AllowedOrigins

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(SessionMiddleware(sessionMgr, logger))
	r.Use(middleware.Recoverer)

	api := &Api{
		Router:         r,
		logger:         logger,
		SessionHandler: sessionHdl,
		SessionManager: sessionMgr,
		ChatHandler:    chatHdl,
		ProfileHandler: profileHdl,
	}

	api.CreateRoutes()

	return api
}

func (api *Api) CreateRoutes() {
	api.Router.Get("/chat/ws", api.ChatHandler.JoinChat)
	api.Router.Post("/register", api.SessionHandler.Register)
	api.Router.Get("/profile", api.ProfileHandler.GetProfile)
}
