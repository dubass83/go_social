package cache

import (
	"context"

	"github.com/dubass83/go_social/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStoreCache() *StoreCache {
	return &StoreCache{
		User: &MockUserCache{},
	}
}

type MockUserCache struct {
	mock.Mock
}

func (muc *MockUserCache) Get(ctx context.Context, id int64) (*store.User, error) {
	args := muc.Called(id)
	return nil, args.Error(1)
}
func (muc *MockUserCache) Set(ctx context.Context, user *store.User) error {
	args := muc.Called(user)
	return args.Error(0)
}
