package store

import (
	"context"
	"database/sql"
)

type CommentsStore struct {
	db *sql.DB
}

func NewCommentsStore(db *sql.DB) *CommentsStore {
	return &CommentsStore{
		db: db,
	}
}

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

func (cs *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
	   INSERT INTO comments (post_id, user_id, content)
       VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := cs.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommentsStore) GetByPostID(ctx context.Context, id int64) ([]Comment, error) {
	query := `
	    SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id
        FROM comments as c
        JOIN users ON users.id = c.user_id
        WHERE c.post_id = $1
        ORDER BY c.created_at DESC;
        `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := cs.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil

}
