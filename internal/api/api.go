package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/h1divp/echo-chat-v2/internal/chat"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/user"
)

type Api struct {
	Router      *chi.Mux
	Logger      *zerolog.Logger
	UserHandler *user.Handler
	ChatHandler *chat.Handler
}

func New(logger *zerolog.Logger, userHdl *user.Handler, chatHdl *chat.Handler) *Api {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	api := &Api{
		Router:      r,
		Logger:      logger,
		UserHandler: userHdl,
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

	api.Router.Get("/ws/chat", api.ChatHandler.HandleWS)
}
