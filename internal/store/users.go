package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"github.com/dubass83/go_social/internal/util"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"-"`
	CreatedAt       string `json:"created_at"`
	Active          bool   `json:"active"`
	ActivationToken string `json:"activation_token"`
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
		&user.Active,
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

	plainToken := util.GenerateToken(user.ID)
	hash := sha256.Sum256([]byte(plainToken))
	token := fmt.Sprintf("%x", hash)

	expiry := time.Now().Add(30 * time.Minute)

	invite := &Invitation{
		UserID: user.ID,
		Token:  token,
		Expiry: expiry,
	}
	err := us.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		token,
		expiry,
	).Scan(
		&invite.ID,
	)

	if err != nil {
		return err
	}

	log.Debug().Msgf("user: %v", user)
	log.Debug().Msgf("invite: %v", invite)

	return nil
}

func (us *UsersStore) CreateAndInviteTx(ctx context.Context, user *User) error {
	return withTx(us.db, ctx, func(tx *sql.Tx) error {
		if err := createUserTx(ctx, tx, user); err != nil {
			return err
		}
		if err := createInvitationTx(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
}

func (us *UsersStore) Activate(ctx context.Context, plainToken string) error {
	query := `
	UPDATE users
	SET active = TRUE
	WHERE id = (
		SELECT user_id
		FROM invitations
		WHERE token = $1
		AND expiry > NOW()
	)
	RETURNING id
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(plainToken))
	token := fmt.Sprintf("%x", hash)

	var userID int64
	err := us.db.QueryRowContext(
		ctx,
		query,
		token,
	).Scan(
		&userID,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	log.Debug().Msgf("user: %v", userID)

	return nil
}

func (us *UsersStore) DeleteByID(ctx context.Context, userID int64) error {
	query := `
	DELETE FROM users
	WHERE id = $1 and active = FALSE
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := us.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (us *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
    SELECT *
    FROM users
    WHERE email = $1
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := us.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.Active,
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
