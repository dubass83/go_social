package cache

import (
	"context"

	"github.com/dubass83/go_social/internal/store"
)

func NewMockStoreCache() *StoreCache {
	return &StoreCache{
		User: &MockUserCache{},
	}
}

type MockUserCache struct{}

func (muc *MockUserCache) Get(ctx context.Context, id int64) (*store.User, error) {
	return &store.User{}, nil
}
func (muc *MockUserCache) Set(ctx context.Context, user *store.User) error {
	return nil
}
