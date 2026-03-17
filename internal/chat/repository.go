package chat

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Repository struct {
	logger zerolog.Logger
	rdb    *redis.Client
}

func NewRepository(logger zerolog.Logger, rdb *redis.Client) *Repository {
	return &Repository{
		logger: logger.With().Str("repository", "chat").Logger(),
		rdb:    rdb,
	}
}

func (r *Repository) FindNearbyUserIDs(ctx context.Context, lat, lon, radius float64) ([]string, error) {
	return r.rdb.GeoSearch(ctx, "user_locations", &redis.GeoSearchQuery{
		Latitude:   lat,
		Longitude:  lon,
		Radius:     radius,
		RadiusUnit: "m",
	}).Result()
}
