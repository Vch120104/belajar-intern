// Package cache provides utility functions for interacting with Redis.
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetRedisValue sets the value for a specific key in the Redis cache.
func SetRedisValue(ctx context.Context, rdb *redis.Client, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set the key's value in Redis with JSON data and a 10-minute expiration.
	err = rdb.Set(ctx, key, jsonData, 10*time.Minute).Err()
	return err
}

// GetRedisValue retrieves the value from the Redis cache based on the key.
func GetRedisValue(ctx context.Context, rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
