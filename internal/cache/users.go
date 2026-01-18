package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-redis/redis/v8"
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
	user := &store.User{}
	idStr := uch.getUserKey(id)
	err := uch.rdb.Get(ctx, idStr).Scan(user)
	if err == redis.Nil {
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}
	return user, nil
}

func (uch *userCache) Set(ctx context.Context, user *store.User) error {
	idStr := uch.getUserKey(user.ID)
	return uch.rdb.Set(ctx, idStr, user, uch.ttl).Err()
}
