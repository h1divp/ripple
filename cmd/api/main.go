package main

import (
	"net/http"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/rs/zerolog/log"

	"github.com/h1divp/echo-chat-v2/internal/api"
	"github.com/h1divp/echo-chat-v2/internal/db"
	"github.com/h1divp/echo-chat-v2/internal/logger"
)

func main() {
	logger := logger.New()
	config := config.Load()

	// redis
	// middlewear

	store, err := store.New(config.DatabaseURL)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	defer store.Close()

	api := api.CreateApi(&logger)

	logger.Info().Msgf("Api is listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, api.Router); err != nil {
		log.Fatal().Msgf("Failed to start http server!")
	}
}
