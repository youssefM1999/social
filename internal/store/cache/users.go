package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/youssefM1999/social/internal/store"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Hour * 24

func (u *UserStore) Get(ctx context.Context, userId int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userId)
	data, err := u.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return u.rdb.SetEX(ctx, cacheKey, jsonData, UserExpTime).Err()
}
