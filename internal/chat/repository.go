package chat

// TODO: convert userIDs to uuids as passed values to ensure type integrity. Ensure that uuids stored as strings are de-serialized back to uuids correctly.
import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const userLocationsKey = "user_locations"

var ErrNoLocationFound = errors.New("No location found for user.")

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

func (r *Repository) GetLocationFromUserID(ctx context.Context, userID uuid.UUID) (float64, float64, error) {
	pos, err := r.rdb.GeoPos(ctx, userLocationsKey, userID.String()).Result()
	if err != nil {
		r.logger.Err(err).Msg("Could not get location from userID")
		return 0, 0, err
	}

	if len(pos) == 0 || pos[0] == nil {
		return 0, 0, ErrNoLocationFound
	}

	return pos[0].Latitude, pos[0].Longitude, nil
}

func (r *Repository) FindNearbyUserIDs(ctx context.Context, lat, lon, radius float64) ([]uuid.UUID, error) {
	idStrings, err := r.rdb.GeoSearch(ctx, userLocationsKey, &redis.GeoSearchQuery{
		Latitude:   lat,
		Longitude:  lon,
		Radius:     radius,
		RadiusUnit: "m",
	}).Result()

	if err != nil {
		r.logger.Err(err).Msg("Could not retrieve nearby user IDs.")
		return nil, err
	}

	UUIDs := make([]uuid.UUID, 0, len(idStrings))
	for _, idStr := range idStrings {
		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			r.logger.Warn().Msg("Skipping malformed UUID from Redis")
			continue
		}
		UUIDs = append(UUIDs, parsedID)
	}

	return UUIDs, nil

}

func (r *Repository) UpdateUserLocation(ctx context.Context, userID uuid.UUID, lat, lon float64) error {
	err := r.rdb.GeoAdd(ctx, userLocationsKey, &redis.GeoLocation{
		// TODO: fix
		Name:      userID.String(),
		Latitude:  lat,
		Longitude: lon,
	}).Err()

	if err != nil {
		r.logger.Err(err).
			Str("userID", userID.String()).
			Msg("Failed to update user location in Redis")
		return err
	}
	r.logger.Debug().Float64("lat", lat).Float64("lon", lon).Msg("Updated location")
	return nil
}

func (r *Repository) RemoveUserLocation(ctx context.Context, userID uuid.UUID) error {
	err := r.rdb.ZRem(ctx, userLocationsKey, userID.String()).Err()
	if err != nil {
		r.logger.Err(err).Str("UserID", userID.String()).Msg("Failed to remove user location")
		return err
	}
	r.logger.Debug().Str("UserID", userID.String()).Msg("Deleted user location from redis")
	return nil
}

func (r *Repository) RemoveAllLocations(ctx context.Context) error {
	err := r.rdb.Del(ctx, userLocationsKey).Err()
	if err != nil {
		r.logger.Err(err).Msg("Failed to remove all locations during cleanup")
		return err
	}
	r.logger.Info().Msg("Successfully cleaned user_locations in Redis")
	return nil
}
