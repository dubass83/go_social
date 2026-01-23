package main

import (
	"net/http"
	"testing"
)

func TestGetUserByIDHandler(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

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

}
