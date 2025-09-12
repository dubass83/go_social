package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func NewFollowersStore(db *sql.DB) *FollowersStore {
	return &FollowersStore{
		db: db,
	}
}

func (fs *FollowersStore) CreateFollower(ctx context.Context, userID, followerID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return fs.db.QueryRowContext(ctx, query, userID, followerID).Err()
}

func (fs *FollowersStore) DeleteFollower(ctx context.Context, userID, followerID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return fs.db.QueryRowContext(ctx, query, userID, followerID).Err()
}
