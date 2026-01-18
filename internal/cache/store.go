package cache

import (
	"context"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-redis/redis/v8"
)

type StoreCache struct {
	User interface {
		Get(ctx context.Context, id int64) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
	}
}

func NewStoreCache(rdb *redis.Client) *StoreCache {
	return &StoreCache{
		User: NewUserCache(rdb),
	}
}
