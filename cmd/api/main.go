package main

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/h1divp/echo-chat-v2/internal/api"
	"github.com/h1divp/echo-chat-v2/internal/chat"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/db"
	"github.com/h1divp/echo-chat-v2/internal/logger"
)

func main() {
	logger := logger.New().With().Timestamp().Logger()
	cfg := config.Load()

	dbStore, err := db.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer dbStore.Close()

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to parse Redis URL")
	}
	rdb := redis.NewClient(redisOpt)

	// Dependency injections
	// userRepo := user.NewRepository(logger, rdb)
	// userSvc := user.NewService(userRepo)
	// userHdl := user.NewHandler(logger, userSvc)

	chatHub := chat.NewHub(logger, rdb)
	go chatHub.Run()

	chatRepo := chat.NewRepository(logger, rdb)
	chatSvc := chat.NewService(logger, chatRepo, chatHub)
	chatHdl := chat.NewHandler(logger, chatSvc, chatHub)

	api := api.New(&logger, nil, chatHdl)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      api.Router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	logger.Info().Msgf("Api is listening on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Failed to start http server!")
	}
}
