package cache

import (
	"context"
	"strconv"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-redis/redis/v8"
)

type userCache struct {
	rdb *redis.Client
}

func (uch *userCache) Get(ctx context.Context, id int64) (*store.User, error) {
	user := &store.User{}
	idStr := strconv.FormatInt(id, 10)
	err := uch.rdb.Get(ctx, idStr).Scan(user)
	if err == redis.Nil {
		return nil, err
	}
	return user, err
}

func (uch *userCache) Set(ctx context.Context, user *store.User) error {
	idStr := strconv.FormatInt(user.ID, 10)
	return uch.rdb.Set(ctx, idStr, user, 0).Err()
}
