package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dubass83/go_social/internal/cache"
	"github.com/dubass83/go_social/internal/store"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	mockStorage := store.NewMockStorage()
	mockCache := cache.NewMockStoreCache()

	return &application{
		store: mockStorage,
		cache: mockCache,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}
