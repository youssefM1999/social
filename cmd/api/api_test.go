package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/youssefM1999/social/internal/ratelimiter"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Second * 5,
			Enabled:              true,
		},
	}

	app := newTestApplication(t, cfg)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := ts.Client()
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < cfg.rateLimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("X-Forwarded-For", mockIP)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()
		if i < cfg.rateLimiter.RequestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status %d, got %d, index: %d", http.StatusOK, resp.StatusCode, i)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Fatalf("expected status %d, got %d, index: %d", http.StatusTooManyRequests, resp.StatusCode, i)
			}
		}
	}
}
