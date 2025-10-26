package store

import (
	"context"
	"database/sql"
	"time"
)

type Invitation struct {
	ID     int64     `json:"id"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
	UserID int64     `json:"user_id"`
}

type InvitationStore struct {
	db *sql.DB
}

func NewInvitationStore(db *sql.DB) *InvitationStore {
	return &InvitationStore{db: db}
}

func (inv *InvitationStore) CleanByID(ctx context.Context, userID int64) error {
	query := `DELETE FROM invitations WHERE user_id = $1`
	_, err := inv.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
