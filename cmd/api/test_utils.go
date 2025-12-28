package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/youssefM1999/social/internal/auth"
	"github.com/youssefM1999/social/internal/store"
	"github.com/youssefM1999/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStorage()
	mockCache := cache.NewMockStore()
	testAuth := auth.NewTestAuthenticator()
	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCache,
		authenticator: testAuth,
		config:        cfg,
	}
}

func executeRequest(mux http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected status code %d, got %d", expected, actual)
	}
}
