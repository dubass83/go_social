package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	UserID    int64     `json:"user_id"`
	Version   int64     `json:"version"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
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

func (ps *PostsStore) GetByID(ctx context.Context, id string) (*Post, error) {
	query := `
	SELECT id, title, content, created_at, updated_at, user_id, version, tags
	FROM posts
	WHERE id = $1
	LIMIT 1
	`

	post := &Post{}

	err := ps.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.Version,
		pq.Array(&post.Tags),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return post, nil
}

func (ps *PostsStore) DeleteByID(ctx context.Context, id string) error {
	query := `
	Delete FROM posts
	WHERE id = $1
	`
	err := ps.db.QueryRowContext(ctx, query, id).Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (ps *PostsStore) Update(ctx context.Context, postID int64, version int64, post *Post) error {
	query := `
	   UPDATE posts SET title = $1, content = $2, tags = $3, updated_at = NOW(), version = version + 1
       WHERE id = $4 AND version = $5
       RETURNING id, user_id, created_at, updated_at, version
	`
	err := ps.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		postID,
		version,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		return err
	}
	return nil
}
