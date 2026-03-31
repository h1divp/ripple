package config

import (
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type ChatConfig struct {
	RateLimitWindowSeconds time.Duration `env:"RATE_LIMIT_WINDOW_SECONDS"`
	RateLimitMaxMessages   int           `env:"RATE_LIMIT_MAX_MESSAGES"`
	MaxMessageLength       int           `env:"MAX_MESSAGE_LENGTH"`
	MessageSearchRadius    float64       `env:"MESSAGE_SEARCH_RADIUS_METERS"`
}

type Config struct {
	Port                 string   `env:"PORT" envDefault:"8080"`
	DatabaseURL          string   `env:"DATABASE_URL"`
	RedisURL             string   `env:"REDIS_URL"`
	AllowedOriginsString string   `env:"ALLOWED_ORIGINS"`
	AllowedOrigins       []string ``
	SessionExpireTime    uint32   `env:"SESSION_EXPIRE_TIME_MIN"`
	CookieHashKey        string   `env:"COOKIE_HASH_KEY"`
	CookieBlockKey       string   `env:"COOKIE_BLOCK_KEY"`
	IsProd               bool     `env:"IS_PROD" envDefault:"false"`

	Chat ChatConfig `envPrefix:"CHAT_"`
}

func Load() *Config {
	LoadEnv()

	config, err := env.ParseAs[Config]()
	config.AllowedOrigins = strings.Split(config.AllowedOriginsString, ",")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse .env: %v", err)
	}

	return &config
}

func LoadEnv() {
	files := []string{".env.dev", ".env"}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err != nil {
				log.Fatal().Err(err).Msgf("Failed to load .env file: %s", file)
				break
			}
		}
	}
}
