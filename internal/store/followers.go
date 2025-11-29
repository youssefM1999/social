package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id)
	VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followerID, userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
	DELETE FROM followers
	WHERE follower_id = $1 AND user_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followerID, userID)
	return err
}
