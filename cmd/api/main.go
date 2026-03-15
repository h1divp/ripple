package main

import (
	"net/http"

	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/user"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/h1divp/echo-chat-v2/internal/api"
	"github.com/h1divp/echo-chat-v2/internal/db"
	"github.com/h1divp/echo-chat-v2/internal/logger"
)

func main() {
	logger := logger.New()
	config := config.Load()

	// TODO: add redis & middlewear

	store, err := store.New(config.DatabaseURL)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	defer store.Close()

	redisOpt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	rdb := redis.NewClient(redisOpt)

	// Dependency injections
	userRepo := user.NewRepository(logger, rdb)
	userSvc := user.NewService(userRepo)
	userHdl := user.NewHandler(logger, userSvc)

	api := api.CreateApi(&logger, userHdl)

	logger.Info().Msgf("Api is listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, api.Router); err != nil {
		log.Fatal().Msgf("Failed to start http server!")
	}
}
