package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/rs/zerolog"
)

type Api struct {
	Router *chi.Mux
	Logger *zerolog.Logger
}

func CreateApi(logger *zerolog.Logger) *Api {
	api := &Api{
		Router: chi.NewRouter(),
		Logger: logger,
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

	api.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello~"))
	})
}
