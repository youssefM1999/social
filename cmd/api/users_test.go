package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/youssefM1999/social/internal/store"
	"github.com/youssefM1999/social/internal/store/cache"
)

func TestGetUser(t *testing.T) {
	config := config{
		redisCfg: redisConfig{
			enabled: true,
		},
	}
	app := newTestApplication(t, config)
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

		rr := executeRequest(mux, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		// First Get call (from middleware) - cache miss
		mockCacheStore.On("Get", int64(1)).Return(nil, nil).Once()
		// Set call (from middleware) - cache the user
		mockCacheStore.On("Set", mock.Anything).Return(nil).Once()
		// Second Get call (from handler) - cache hit (user was cached by middleware)
		cachedUser := &store.User{ID: 1}
		mockCacheStore.On("Get", int64(1)).Return(cachedUser, nil).Once()

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(mux, req)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.AssertNumberOfCalls(t, "Set", 1)
		mockCacheStore.Calls = nil // reset the calls expectations
	})

	t.Run("should hit cache first and if not exists set the user in the cache", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		mockCacheStore.On("Get", int64(1)).Return(nil, nil)
		mockCacheStore.On("Get", int64(2)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/2", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(mux, req)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)

		mockCacheStore.Calls = nil // reset the calls expectations
	})
}
