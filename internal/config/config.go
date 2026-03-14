package config

import (
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port                 string   `env:"PORT" envDefault:"8080"`
	DatabaseURL          string   `env:"PSQL_DSN"`
	AllowedOriginsString string   `env:"ALLOWED_ORIGINS"`
	AllowedOrigins       []string ``
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
