// Package store provides a data access layer for the application.
package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (ps *PostsStore) GetUserFeed(ctx context.Context, userID int64, pg PaginatedFeedQuery) ([]*PostWithMetadata, error) {
	var tagsCondition string

	if len(pg.Tags) == 0 {
		tagsCondition = "($5 = $5 OR TRUE)"
	} else {
		tagsCondition = "(p.tags @> $5)"
	}

	query := `
	SELECT
      p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
      u.username,
      COUNT(c.id) AS comments_count
    FROM posts p
    LEFT JOIN users u ON u.id = p.user_id
    LEFT JOIN comments c ON p.id = c.post_id
    WHERE
      (p.user_id = $1 OR p.user_id IN (
           SELECT follower_id
           FROM followers
           WHERE user_id = $1
      ))
      AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
      AND ` + tagsCondition + `
    GROUP BY p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username
    ORDER BY p.created_at ` + pg.Sort + `, p.id ` + pg.Sort + `
    LIMIT $2 OFFSET $3;
    `

	log.Debug().Msgf("userID: %d, limit: %d, offset: %d, tags: %+v, search: '%s'",
		userID, pg.Limit, pg.Offset, pg.Tags, pg.Search)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := ps.db.QueryContext(ctx, query, userID, pg.Limit, pg.Offset, pg.Search, pq.Array(pg.Tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feed := []*PostWithMetadata{}
	for rows.Next() {
		p := &PostWithMetadata{}
		err = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}
	// Check for any iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return feed, nil
}

func (ps *PostsStore) GetByID(ctx context.Context, id string) (*Post, error) {
	query := `
	SELECT id, title, content, created_at, updated_at, user_id, version, tags
	FROM posts
	WHERE id = $1
	LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
