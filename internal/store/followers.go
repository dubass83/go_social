package store

import (
	"context"
	"database/sql"
)

type Follow struct {
	UserID    int64  `json:"user_id"`
	FollowID  int64  `json:"follow_id"`
	CreatedAt string `json:"created_at"`
}

type FollowsStore struct {
	db *sql.DB
}

func NewFollowsStore(db *sql.DB) *FollowsStore {
	return &FollowsStore{
		db: db,
	}
}

func (fs *FollowsStore) CreateFollow(ctx context.Context, userID, followID int64) error {
	query := `INSERT INTO followers (user_id, follow_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return fs.db.QueryRowContext(ctx, query, userID, followID).Err()
}

func (fs *FollowsStore) DeleteFollow(ctx context.Context, userID, followID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follow_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return fs.db.QueryRowContext(ctx, query, userID, followID).Err()
}
