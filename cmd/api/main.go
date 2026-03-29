package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/h1divp/echo-chat-v2/internal/api"
	"github.com/h1divp/echo-chat-v2/internal/chat"
	"github.com/h1divp/echo-chat-v2/internal/config"
	"github.com/h1divp/echo-chat-v2/internal/db"
	"github.com/h1divp/echo-chat-v2/internal/logger"
	"github.com/h1divp/echo-chat-v2/internal/profile"
	"github.com/h1divp/echo-chat-v2/internal/session"
	"github.com/h1divp/echo-chat-v2/internal/websocket"
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

	/*
		Dependency injections
	*/

	profileSvc := profile.NewService(logger)
	profileHdl := profile.NewHandler(logger, profileSvc)

	hashKey, err := base64.StdEncoding.DecodeString(cfg.CookieHashKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to decide COOKIE_HASH_KEY")
	}
	blockKey, err := base64.StdEncoding.DecodeString(cfg.CookieBlockKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to decide COOKIE_BLOCK_KEY")
	}
	sessionMgr := session.NewManager(logger, rdb, hashKey, blockKey)
	sessionHdl := session.NewHandler(logger, sessionMgr, *cfg)

	wsHub := websocket.NewHub(logger)
	chatRepo := chat.NewRepository(logger, rdb)
	chatSvc := chat.NewService(logger, chatRepo, wsHub, sessionMgr)
	chatHdl := chat.NewHandler(logger, chatSvc, wsHub, profileSvc, sessionMgr)

	// Messy but needed...
	wsHub.OnDisconnect = func(sessionID uuid.UUID, userID uuid.UUID) {
		chatSvc.HandleDisconnect(context.Background(), sessionID, userID)
	}

	// Remove stale user_locations
	// We do not defer in case of needing to clean up potentially after a panic
	if err := chatRepo.RemoveAllLocations(context.Background()); err != nil {
		logger.Warn().Msg("Could not clean up stale user_locations from Redis.")
	}
	if err := sessionMgr.RemoveAllSessions(context.Background()); err != nil {
		logger.Warn().Msg("Could not clean up stale sessions from Redis.")
	}

	go wsHub.Run()
	api := api.New(&logger, sessionHdl, sessionMgr, chatHdl, profileHdl)

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
