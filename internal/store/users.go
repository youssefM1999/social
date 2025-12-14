package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
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
		&user.Password.hash,
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

func (u *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		// Find and validate the invitation
		user, err := u.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// Update the user
		user.IsActive = true
		err = u.update(ctx, tx, user)
		if err != nil {
			return err
		}

		// Delete the invitation
		err = u.deleteUserInvitation(ctx, tx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExpiry time.Duration) error {
	// transaction wrapper
	// create the user
	// create the user invite
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := u.createUserInvitation(ctx, tx, user.ID, token, invitationExpiry); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, userID int64, token string, exp time.Duration) error {
	query := `
	INSERT INTO user_invitations (user_id, token, expiry)
	VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		userID,
		token,
		time.Now().Add(exp),
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
	DELETE FROM user_invitations
	WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
	SELECT u.id, u.username, u.email, u.created_at, u.is_active
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	user := User{}
	err := tx.QueryRowContext(ctx,
		query,
		hashToken,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	UPDATE users
	SET username = $1, email = $2, is_active = $3
	WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}
