package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ratelimiter "github.com/dubass83/go_social/internal/rateLimiter"
)

func TestRateLimiterMiddleware(t *testing.T) {
	app := newTestApplication(t)

	testRateLimitConfig := ratelimiter.Config{
		RequestPerTimeFrame: 20,
		TimeFrame:           5 * time.Second,
		Enabled:             true,
	}
	app.config.rateLimiter = testRateLimitConfig
	app.rateLimiter = ratelimiter.NewFixedWindowLimeter(testRateLimitConfig)
	app.config.addr = ":8080"

	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := range app.config.rateLimiter.RequestPerTimeFrame + marginOfError {

		req, err := http.NewRequest("GET", ts.URL+"/", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not sent rrequest: %v", err)
		}
		defer resp.Body.Close()

		if i < app.config.rateLimiter.RequestPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status code %d, got %d", http.StatusTooManyRequests, resp.StatusCode)
			}
		}

	}

}
