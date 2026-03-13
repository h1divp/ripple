package main

import (
	"net/http"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/rs/zerolog/log"

	"github.com/h1divp/echo-chat-v2/internal/api"
	"github.com/h1divp/echo-chat-v2/internal/logger"
)

// init logger
// init config
// init api service

func main() {
	logger := logger.New()
	config := config.Load()

	// redis
	// middlewear

	// repos, services, handlers

	api := api.CreateApi(&logger)

	logger.Info().Msgf("Api is listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, api.Router); err != nil {
		log.Fatal().Msgf("Failed to start http server!")
	}
}
