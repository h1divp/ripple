package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/user"
)

type Api struct {
	Router      *chi.Mux
	Logger      *zerolog.Logger
	UserHandler *user.Handler
}

func CreateApi(logger *zerolog.Logger, userHdl *user.Handler) *Api {
	api := &Api{
		Router:      chi.NewRouter(),
		Logger:      logger,
		UserHandler: userHdl,
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

	api.Router.Get("/ws/chat", api.UserHandler.Connect)
}
