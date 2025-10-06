package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/dubass83/go_social/internal/util"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

func (us *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := us.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (us *UsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
    SELECT *
    FROM users
    WHERE ID = $1
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := us.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (us *UsersStore) CreateAndInvite(ctx context.Context, user *User) error {
	if err := us.Create(ctx, user); err != nil {
		return err
	}

	query := `
	INSERT INTO invitations (user_id, token, expiry)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := us.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		util.GenerateToken(user.ID),
		time.Now().Add(30*time.Minute),
	).Scan(
		&user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
