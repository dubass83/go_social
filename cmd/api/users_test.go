package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dubass83/go_social/internal/cache"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByIDHandler(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		app.config.cache.enable = true
		mockCacheStore := app.cache.User.(*cache.MockUserCache)

		mockCacheStore.On("Get", int64(42)).Return(nil, fmt.Errorf("no such user in the cache"))
		// mockCacheStore.On("Get", int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		mockCacheStore.AssertNumberOfCalls(t, "Get", 1)
		mockCacheStore.AssertNumberOfCalls(t, "Set", 1)
		mockCacheStore.Calls = nil // Reset the calls to avoid interference with other tests
	})

}
