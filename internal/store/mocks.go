package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStorage() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, expiry time.Duration) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
