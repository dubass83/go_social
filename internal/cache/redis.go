// Package cache provides caching functionality using Redis.
// It includes utilities for creating and managing Redis client connections.
package cache

import (
	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates and returns a new Redis client instance with the provided configuration.
//
// Parameters:
//   - addr: Redis server address in the format "host:port" (e.g., "localhost:6379")
//   - pw: Password for Redis authentication (use empty string if no password is required)
//   - db: Redis database number to use (0-15, typically 0 for default)
//
// Returns:
//   - A configured *redis.Client ready for use
//
// Example:
//
//	client := NewRedisClient("localhost:6379", "", 0)
//	defer client.Close()
func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}
