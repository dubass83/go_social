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

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func createUserTx(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
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

func createInvitationTx(ctx context.Context, tx *sql.Tx, user *User) error {
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

	_, err := tx.ExecContext(ctx, query, user.ID, token, expiry)

	log.Debug().Msgf("Invitation token: %s  for userID: %d", plainToken, user.ID)

	if err != nil {
		return err
	}

	return nil
}
