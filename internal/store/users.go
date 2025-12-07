package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
	SELECT id, username, email, password, created_at 
	FROM users 
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	var user User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string) error {
	// transaction wrapper
	// create the user
	// create the user invite
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		// TODO: Implement the transaction logic once done with documentation
		return nil

	})
}
