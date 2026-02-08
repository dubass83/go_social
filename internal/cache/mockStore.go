package cache

import (
	"context"

	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog/log"
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
	log.Debug().Msgf("MockUserCache.Get(%d)", id)
	args := muc.Called(ctx, id)
	if args.Get(0) == nil {
		log.Debug().Msgf("MockUserCache.Get(%d)", id)
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}
func (muc *MockUserCache) Set(ctx context.Context, user *store.User) error {
	args := muc.Called(ctx, user)
	return args.Error(0)
}
