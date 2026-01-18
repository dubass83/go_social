package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

const (
	// UserKeyPrefix is the prefix for user cache keys
	UserKeyPrefix = "user:"
	// UserCacheTTL is the default TTL for cached users
	UserCacheTTL = 24 * time.Hour
)

type userCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewUserCache creates a new user cache instance with the provided Redis client
// and default TTL for cached users
func NewUserCache(rdb *redis.Client) *userCache {
	return &userCache{
		rdb: rdb,
		ttl: UserCacheTTL,
	}
}

// getUserKey generates a Redis key for a user ID
func (uch *userCache) getUserKey(id int64) string {
	return UserKeyPrefix + strconv.FormatInt(id, 10)
}

func (uch *userCache) Get(ctx context.Context, id int64) (*store.User, error) {
	// If Redis client is not configured, return cache miss
	if uch.rdb == nil {
		return nil, nil
	}

	key := uch.getUserKey(id)

	// Get the JSON string from Redis
	data, err := uch.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Cache miss - not an error, just return nil
		log.Debug().Str("key", key).Msg("Cache miss")
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	// Unmarshal JSON into user struct
	var user store.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user from cache: %w", err)
	}

	log.Debug().Int64("user_id", id).Str("username", user.Username).Msg("Cache hit")
	return &user, nil
}

func (uch *userCache) Set(ctx context.Context, user *store.User) error {
	// If Redis client is not configured, silently skip
	if uch.rdb == nil {
		return nil
	}

	key := uch.getUserKey(user.ID)

	// Marshal user to JSON
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user for cache: %w", err)
	}

	// Store in Redis with TTL
	if err := uch.rdb.Set(ctx, key, data, uch.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set user in cache: %w", err)
	}

	log.Debug().Int64("user_id", user.ID).Str("username", user.Username).Msg("User cached")
	return nil
}
