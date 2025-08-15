package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
}

type PostsStore struct {
	db *sql.DB
}

func NewPostsStore(db *sql.DB) *PostsStore {
	return &PostsStore{
		db: db,
	}
}

func (ps *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	   INSERT INTO posts (title, content, user_id, tags)
       VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := ps.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
