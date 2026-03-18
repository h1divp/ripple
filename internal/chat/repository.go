package chat

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const userLocationsKey = "user_locations"

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
	return r.rdb.GeoSearch(ctx, userLocationsKey, &redis.GeoSearchQuery{
		Latitude:   lat,
		Longitude:  lon,
		Radius:     radius,
		RadiusUnit: "m",
	}).Result()
}

func (r *Repository) UpdateUserLocation(ctx context.Context, userID string, lat, lon float64) error {
	err := r.rdb.GeoAdd(ctx, userLocationsKey, &redis.GeoLocation{
		Name:      userID,
		Latitude:  lat,
		Longitude: lon,
	}).Err()

	if err != nil {
		r.logger.Err(err).
			Str("userID", userID).
			Msg("Failed to update user location in Redis")
		return err
	}
	r.logger.Debug().Float64("lat", lat).Float64("lon", lon).Msg("Updated location")
	return nil
}

func (r *Repository) RemoveUserLocation(ctx context.Context, userID string) error {
	err := r.rdb.ZRem(ctx, userLocationsKey, userID).Err()
	if err != nil {
		r.logger.Err(err).Str("UserID", userID).Msg("Failed to remove user location")
		return err
	}
	r.logger.Debug().Str("UserID", userID).Msg("Removed user location")
	return nil
}
